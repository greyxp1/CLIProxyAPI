package pool

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
	"time"

	coreauth "github.com/router-for-me/CLIProxyAPI/v8/sdk/cliproxy/auth"
)

const (
	claudeUsageURL = "https://api.anthropic.com/api/oauth/usage"
	codexUsageURL  = "https://chatgpt.com/backend-api/wham/usage"
	creditURL      = "https://chatgpt.com/backend-api/wham/rate-limit-reset-credits"
)

// RequestAccount uses native executor transport and token refresh without looping
// through the management API. Callers supply only canonical provider URLs.
func RequestAccount(ctx context.Context, manager *coreauth.Manager, account *coreauth.Auth, method, target string, body []byte) (*http.Response, error) {
	if account == nil || manager == nil {
		return nil, fmt.Errorf("pool: account unavailable")
	}
	current, ok := manager.GetByID(account.ID)
	if !ok || current.Disabled {
		return nil, fmt.Errorf("pool: account unavailable")
	}
	account = current
	for attempt := 0; attempt < 2; attempt++ {
		request, err := http.NewRequestWithContext(ctx, method, target, bytes.NewReader(body))
		if err != nil {
			return nil, err
		}
		request.Header.Set("Content-Type", "application/json")
		if account.Provider == "claude" {
			request.Header.Set("anthropic-beta", "oauth-2025-04-20")
			request.Header.Set("anthropic-version", "2023-06-01")
			request.Header.Set("User-Agent", "claude-code/2.1.273")
		} else if account.Provider == "codex" {
			request.Header.Set("OpenAI-Beta", "codex-1")
			request.Header.Set("Originator", "Codex Desktop")
			if id, ok := account.Metadata["account_id"].(string); ok && id != "" {
				request.Header.Set("Chatgpt-Account-Id", id)
			}
		} else {
			return nil, fmt.Errorf("pool: unsupported provider")
		}
		response, err := manager.HttpRequest(ctx, account, request)
		if err != nil {
			return nil, err
		}
		if response.StatusCode != http.StatusUnauthorized || attempt != 0 {
			return response, nil
		}
		_ = response.Body.Close()
		account, err = manager.ForceRefreshAuth(ctx, account.ID)
		if err != nil || account == nil {
			return nil, fmt.Errorf("pool: credential refresh failed")
		}
	}
	return nil, fmt.Errorf("pool: credential unavailable")
}

func (m *Maintenance) Call(ctx context.Context, account *coreauth.Auth, target string, body any) ([]byte, error) {
	var encoded []byte
	method := http.MethodGet
	if body != nil {
		var err error
		encoded, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
		method = http.MethodPost
	}
	response, err := RequestAccount(ctx, m.manager, account, method, target, encoded)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("provider HTTP %d", response.StatusCode)
	}
	data, err := io.ReadAll(io.LimitReader(response.Body, 4*1024*1024+1))
	if err != nil {
		return nil, err
	}
	if len(data) > 4*1024*1024 {
		return nil, fmt.Errorf("pool: provider response too large")
	}
	return data, nil
}

func (m *Maintenance) Usage(ctx context.Context, account *coreauth.Auth) ([]byte, error) {
	target := codexUsageURL
	if account.Provider == "claude" {
		target = claudeUsageURL
	}
	data, err := m.Call(ctx, account, target, nil)
	if err != nil {
		return nil, err
	}
	snapshot, err := parseUsage(account.Provider, data, time.Now())
	if err != nil {
		return nil, err
	}
	// Banked resets come from the maintenance listing. Keep them across usage reads.
	if previous, ok := m.quotas.Get(account.ID); ok {
		snapshot.Credits, snapshot.CreditExpiry = previous.Credits, previous.CreditExpiry
		snapshot.ResetCredits = previous.ResetCredits
	}
	m.quotas.Set(account.ID, snapshot)
	return data, nil
}

func (m *Maintenance) Credits(ctx context.Context, account *coreauth.Auth) ([]byte, error) {
	if account.Provider != "codex" {
		return nil, fmt.Errorf("pool: reset credits require Codex")
	}
	return m.Call(ctx, account, creditURL, nil)
}

func (m *Maintenance) Probe(ctx context.Context, account *coreauth.Auth) (coreauth.QuotaSnapshot, error) {
	_, err := m.Usage(ctx, account)
	if err != nil {
		return coreauth.QuotaSnapshot{}, err
	}
	snapshot, _ := m.quotas.Get(account.ID)
	return snapshot, nil
}

func parseUsage(provider string, data []byte, now time.Time) (coreauth.QuotaSnapshot, error) {
	snapshot := coreauth.QuotaSnapshot{CheckedAt: now}
	if provider == "claude" {
		var fields map[string]json.RawMessage
		if err := json.Unmarshal(data, &fields); err != nil {
			return snapshot, fmt.Errorf("pool: invalid Claude usage")
		}
		for name, raw := range fields {
			if name != "five_hour" && name != "seven_day" && !strings.HasPrefix(name, "seven_day_") {
				continue
			}
			if string(raw) == "null" {
				if name == "five_hour" {
					snapshot.Short = &coreauth.QuotaWindow{Duration: 5 * time.Hour}
				}
				continue
			}
			var window struct {
				Used  *float64 `json:"utilization"`
				Reset *string  `json:"resets_at"`
			}
			if json.Unmarshal(raw, &window) != nil || window.Used == nil {
				continue
			}
			parsed, err := parseClaudeWindow(*window.Used, window.Reset, name)
			if err != nil {
				return snapshot, err
			}
			if name == "five_hour" {
				snapshot.Short = &parsed
			} else {
				snapshot.Long = append(snapshot.Long, parsed)
			}
		}
		var scoped struct {
			Limits []struct {
				Kind     string   `json:"kind"`
				Percent  *float64 `json:"percent"`
				ResetsAt *string  `json:"resets_at"`
				Scope    *struct {
					Model *struct {
						DisplayName string `json:"display_name"`
					} `json:"model"`
				} `json:"scope"`
			} `json:"limits"`
		}
		if err := json.Unmarshal(data, &scoped); err != nil {
			return snapshot, fmt.Errorf("pool: invalid scoped usage")
		}
		for _, limit := range scoped.Limits {
			if limit.Kind != "weekly_scoped" || limit.Percent == nil || limit.Scope == nil || limit.Scope.Model == nil {
				continue
			}
			name := strings.ToLower(limit.Scope.Model.DisplayName)
			if name == "" {
				continue
			}
			for _, family := range []string{"haiku", "sonnet", "opus", "fable"} {
				if strings.Contains(name, family) {
					name = family
					break
				}
			}
			parsed, err := parseClaudeWindow(*limit.Percent, limit.ResetsAt, "seven_day_"+name)
			if err != nil {
				return snapshot, err
			}
			snapshot.Long = append(snapshot.Long, parsed)
		}
	} else if provider == "codex" {
		var response struct {
			Plan      string                     `json:"plan_type"`
			RateLimit map[string]json.RawMessage `json:"rate_limit"`
		}
		if json.Unmarshal(data, &response) != nil || response.RateLimit == nil {
			return snapshot, fmt.Errorf("pool: invalid Codex usage")
		}
		snapshot.Plan = response.Plan
		for _, raw := range response.RateLimit {
			var window struct {
				Used     *float64 `json:"used_percent"`
				Reset    *float64 `json:"reset_at"`
				Duration float64  `json:"limit_window_seconds"`
			}
			if json.Unmarshal(raw, &window) != nil || window.Used == nil || window.Duration <= 0 {
				continue
			}
			if math.IsNaN(*window.Used) || math.IsInf(*window.Used, 0) || *window.Used < 0 {
				return snapshot, fmt.Errorf("pool: invalid usage percentage")
			}
			parsed := coreauth.QuotaWindow{Used: *window.Used, Duration: time.Duration(window.Duration * float64(time.Second))}
			if window.Reset != nil {
				value := time.Unix(int64(*window.Reset), 0)
				parsed.Reset = &value
			}
			if parsed.Duration == 5*time.Hour {
				snapshot.Short = &parsed
			} else if parsed.Duration > 5*time.Hour {
				snapshot.Long = append(snapshot.Long, parsed)
			}
		}
	} else {
		return snapshot, fmt.Errorf("pool: unsupported usage provider")
	}
	if snapshot.Short == nil && len(snapshot.Long) == 0 {
		return snapshot, fmt.Errorf("pool: usage windows unavailable")
	}
	return snapshot, nil
}

func parseClaudeWindow(used float64, reset *string, name string) (coreauth.QuotaWindow, error) {
	window := coreauth.QuotaWindow{Used: used, Duration: 7 * 24 * time.Hour, Name: name}
	if math.IsNaN(used) || math.IsInf(used, 0) || used < 0 {
		return window, fmt.Errorf("pool: invalid usage percentage")
	}
	if name == "five_hour" {
		window.Duration = 5 * time.Hour
	}
	if reset != nil && *reset != "" {
		value, err := time.Parse(time.RFC3339Nano, *reset)
		if err != nil {
			return window, fmt.Errorf("pool: invalid reset timestamp")
		}
		window.Reset = &value
	}
	return window, nil
}
