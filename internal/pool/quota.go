package pool

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"math/rand/v2"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/MoeForge/moego/json"
	"github.com/router-for-me/CLIProxyAPI/v8/internal/auth/codex"
	sdkauth "github.com/router-for-me/CLIProxyAPI/v8/sdk/auth"
	coreauth "github.com/router-for-me/CLIProxyAPI/v8/sdk/cliproxy/auth"
	log "github.com/sirupsen/logrus"
)

const (
	// cedar_ember=1 adds the banked reset grants, which Anthropic only reports to recent Claude Code versions.
	claudeUsageURL             = "https://api.anthropic.com/api/oauth/usage?cedar_ember=1"
	codexUsageURL              = "https://chatgpt.com/backend-api/wham/usage"
	creditURL                  = "https://chatgpt.com/backend-api/wham/rate-limit-reset-credits"
	claudeProfileURL           = "https://api.anthropic.com/api/oauth/profile"
	antigravityQuotaURL        = sdkauth.AntigravityModelBaseURLDaily + sdkauth.AntigravityModelsPath
	antigravityQuotaSummaryURL = sdkauth.AntigravityModelBaseURLDaily + "/v1internal:retrieveUserQuotaSummary"
)

// RequestAccount uses native executor transport and token refresh without looping
// through the management API. Callers supply only canonical provider URLs.
func RequestAccount(ctx context.Context, manager *coreauth.Manager, account *coreauth.Auth, method, target string, body []byte) (*http.Response, error) {
	if account == nil || manager == nil {
		return nil, fmt.Errorf("pool: account unavailable")
	}
	current, ok := manager.GetByID(account.ID)
	// Disabled accounts still report quota. Only the quota endpoints may use their credentials.
	quotaRead := method == http.MethodGet && (target == claudeUsageURL || target == codexUsageURL) ||
		method == http.MethodPost && (target == antigravityQuotaSummaryURL || target == antigravityQuotaURL)
	if !ok || current.Disabled && !quotaRead {
		return nil, fmt.Errorf("pool: account unavailable")
	}
	account = current
	refresh := func() error {
		refreshed, ok := manager.RefreshAfterUnauthorized(ctx, account, &coreauth.Error{Code: "unauthorized", HTTPStatus: http.StatusUnauthorized}, false)
		if !ok {
			if current, found := manager.GetByID(account.ID); found && needsSignIn(current, time.Now()) {
				return errSignIn
			}
			return fmt.Errorf("pool: credential refresh failed")
		}
		account = refreshed
		return nil
	}
	if needsSignIn(account, time.Now()) {
		return nil, errSignIn
	}
	// Usage endpoints can rate limit even expired credentials, hiding the 401 that triggers renewal.
	if !account.HasValidAccessToken(time.Now()) {
		if err := refresh(); err != nil {
			return nil, err
		}
	}
	for attempt := 0; attempt < 2; attempt++ {
		request, err := http.NewRequestWithContext(ctx, method, target, bytes.NewReader(body))
		if err != nil {
			return nil, err
		}
		request.Header.Set("Content-Type", "application/json")
		if account.Provider == "claude" {
			request.Header.Set("anthropic-beta", "oauth-2025-04-20")
			request.Header.Set("anthropic-version", "2023-06-01")
			request.Header.Set("User-Agent", "claude-cli/2.1.280 (external, cli)")
		} else if account.Provider == "codex" {
			request.Header.Set("OpenAI-Beta", "codex-1")
			request.Header.Set("Originator", "Codex Desktop")
			if id, ok := account.Metadata["account_id"].(string); ok && id != "" {
				request.Header.Set("Chatgpt-Account-Id", id)
			}
		} else if account.Provider != "antigravity" {
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
		if err := refresh(); err != nil {
			return nil, err
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
		return nil, newProviderError(response, time.Now())
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

// Usage reads and stores an account's limits. It skips accounts that need a new sign-in or whose
// usage endpoint is rate limiting them, so maintenance and manual refreshes share one backoff.
func (m *Maintenance) Usage(ctx context.Context, account *coreauth.Auth) ([]byte, error) {
	m.limitMu.Lock()
	limit := m.limits[account.ID]
	if limit == nil {
		limit = &usageLimit{}
		m.limits[account.ID] = limit
	}
	m.limitMu.Unlock()
	// Serialize checks per account so concurrent manual and maintenance reads share the backoff.
	limit.mu.Lock()
	defer limit.mu.Unlock()
	if current, ok := m.manager.GetByID(account.ID); ok {
		account = current
	}
	now := time.Now()
	if needsSignIn(account, now) {
		return nil, errSignIn
	}
	// New credentials end the backoff early.
	if limit.version != account.CredentialVersion {
		limit.until, limit.failures = time.Time{}, 0
		limit.version = account.CredentialVersion
	}
	if limit.until.After(now) {
		return nil, &rateLimitedError{until: limit.until}
	}
	data, err := m.readUsage(ctx, account)
	var provider *providerError
	if err == nil {
		limit.until, limit.failures = time.Time{}, 0
	} else if errors.As(err, &provider) && provider.status == http.StatusTooManyRequests {
		limit.failures++
		limit.until = provider.retryAfter
		if limit.until.IsZero() {
			limit.until = time.Now().Add(usageBackoff(limit.failures))
		}
		accountLog(account).WithFields(log.Fields{"status": provider.status, "code": provider.code, "until": limit.until.Format(time.RFC3339)}).Warn("pool usage endpoint rate limited, pausing checks")
		err = &rateLimitedError{until: limit.until}
	}
	return data, err
}

func (m *Maintenance) readUsage(ctx context.Context, account *coreauth.Auth) ([]byte, error) {
	if account.Provider == "antigravity" {
		return m.antigravityUsage(ctx, account)
	}
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
	// Banked resets come from the maintenance listing, so the usage reading keeps them.
	m.quotas.SetUsage(account.ID, snapshot)
	if err = m.recordUsage(account.ID, snapshot); err != nil {
		log.Warnf("pool history recording failed: %v", err)
	}
	if account.Provider == "claude" && !account.Disabled {
		m.readSubscription(ctx, account)
	}
	return data, nil
}

func (m *Maintenance) antigravityUsage(ctx context.Context, account *coreauth.Auth) ([]byte, error) {
	project, _ := account.Metadata["project_id"].(string)
	body := struct {
		Project string `json:"project"`
	}{Project: project}
	if data, err := m.Call(ctx, account, antigravityQuotaSummaryURL, body); err == nil {
		if windows := antigravitySummaryWindows(data); len(windows) > 0 {
			snapshot := coreauth.QuotaSnapshot{CheckedAt: time.Now(), DisplayWindows: windows}
			m.quotas.SetUsage(account.ID, snapshot)
			if errRecord := m.recordUsage(account.ID, snapshot); errRecord != nil {
				log.Warnf("pool history recording failed: %v", errRecord)
			}
			return data, nil
		}
	}
	data, err := m.Call(ctx, account, antigravityQuotaURL, body)
	if err != nil {
		return nil, err
	}
	var catalog struct {
		Models map[string]struct {
			DisplayName string `json:"displayName"`
			Quota       struct {
				Remaining *float64   `json:"remainingFraction"`
				Reset     *time.Time `json:"resetTime"`
			} `json:"quotaInfo"`
		} `json:"models"`
	}
	if err := json.Unmarshal(data, &catalog); err != nil || catalog.Models == nil {
		return nil, fmt.Errorf("pool: invalid Antigravity quota response")
	}
	windows := make([]coreauth.QuotaWindow, 0, len(catalog.Models))
	labels := make(map[string]int, len(catalog.Models))
	for id, model := range catalog.Models {
		quota := model.Quota
		if quota.Remaining == nil || *quota.Remaining < 0 || *quota.Remaining > 1 {
			continue
		}
		label := strings.TrimSpace(model.DisplayName)
		labels[label]++
		windows = append(windows, coreauth.QuotaWindow{Name: id, Label: label, Used: (1 - *quota.Remaining) * 100, Reset: quota.Reset})
	}
	for i := range windows {
		if windows[i].Label == "" || labels[windows[i].Label] > 1 {
			windows[i].Label = windows[i].Name
		}
	}
	slices.SortFunc(windows, func(a, b coreauth.QuotaWindow) int { return strings.Compare(a.Name, b.Name) })
	m.quotas.SetUsage(account.ID, coreauth.QuotaSnapshot{CheckedAt: time.Now(), DisplayWindows: windows})
	return data, nil
}

func antigravitySummaryWindows(data []byte) []coreauth.QuotaWindow {
	type bucket struct {
		ID        string     `json:"bucketId"`
		Name      string     `json:"displayName"`
		Window    string     `json:"window"`
		Remaining *float64   `json:"remainingFraction"`
		Reset     *time.Time `json:"resetTime"`
		Disabled  bool       `json:"disabled"`
	}
	type group struct {
		Name    string   `json:"displayName"`
		Buckets []bucket `json:"buckets"`
	}
	var summary struct {
		Groups  []group  `json:"groups"`
		Buckets []bucket `json:"buckets"`
	}
	if json.Unmarshal(data, &summary) != nil {
		return nil
	}
	windows := []coreauth.QuotaWindow{}
	seen := map[string]bool{}
	for _, group := range append(summary.Groups, group{Buckets: summary.Buckets}) {
		for _, bucket := range group.Buckets {
			if bucket.Disabled || bucket.ID == "" || seen[bucket.ID] || bucket.Remaining == nil || *bucket.Remaining < 0 || *bucket.Remaining > 1 {
				continue
			}
			name := group.Name
			switch bucket.ID {
			case "gemini-weekly", "gemini-5h":
				name = "Gemini"
			case "3p-weekly", "3p-5h":
				name = "Claude & GPT"
			}
			label := bucket.Name
			var duration time.Duration
			switch bucket.Window {
			case "weekly":
				label = "Weekly"
				duration = 7 * 24 * time.Hour
			case "5h":
				label = "5-hour"
				duration = 5 * time.Hour
			}
			if label = strings.TrimSpace(name + " " + label); label == "" {
				label = bucket.ID
			}
			windows = append(windows, coreauth.QuotaWindow{Name: bucket.ID, Label: label, Used: (1 - *bucket.Remaining) * 100, Reset: bucket.Reset, Duration: duration})
			seen[bucket.ID] = true
		}
	}
	return windows
}

// errSignIn marks an account whose credentials expired after its refresh token was rejected.
var errSignIn = errors.New("pool: account needs a new sign-in")

// needsSignIn reports an account that cannot authenticate until it is signed in again.
// The auth manager clears this as soon as the stored refresh token changes.
func needsSignIn(account *coreauth.Auth, now time.Time) bool {
	return account.RenewalStateAt(now) == coreauth.RenewalStateExpiredSignInRequired
}

type rateLimitedError struct{ until time.Time }

func (e *rateLimitedError) Error() string {
	return "pool: usage checks paused until " + e.until.Format(time.RFC3339)
}

// usageLimit is an account's usage endpoint backoff for one credential version.
type usageLimit struct {
	mu       sync.Mutex
	until    time.Time
	failures int
	version  uint64
}

// usageBackoff doubles from ten minutes up to two hours, with jitter so accounts drift apart.
func usageBackoff(failures int) time.Duration {
	wait := min(10*time.Minute<<min(failures-1, 4), 2*time.Hour)
	return min(wait+rand.N(wait/4), 2*time.Hour)
}

// providerError keeps only the status, Retry-After and a sanitized error code, never the body.
type providerError struct {
	status     int
	code       string
	retryAfter time.Time
}

func (e *providerError) Error() string {
	if e.code == "" {
		return fmt.Sprintf("provider HTTP %d", e.status)
	}
	return fmt.Sprintf("provider HTTP %d (%s)", e.status, e.code)
}

func newProviderError(response *http.Response, now time.Time) *providerError {
	result := &providerError{status: response.StatusCode}
	// Retry-After is either delta seconds or an HTTP date. A day bounds a malformed value.
	if value := strings.TrimSpace(response.Header.Get("Retry-After")); value != "" {
		if seconds, err := strconv.Atoi(value); err == nil && seconds > 0 {
			result.retryAfter = now.Add(time.Duration(min(seconds, 86400)) * time.Second)
		} else if at, errDate := http.ParseTime(value); errDate == nil && at.After(now) {
			result.retryAfter = now.Add(min(at.Sub(now), 24*time.Hour))
		}
	}
	data, _ := io.ReadAll(io.LimitReader(response.Body, 64*1024))
	var body struct {
		Error json.RawMessage `json:"error"`
	}
	if json.Unmarshal(data, &body) != nil || len(body.Error) == 0 {
		return result
	}
	// OAuth errors are a string. Anthropic and OpenAI use an object with a code or type.
	var detail struct {
		Type string `json:"type"`
		Code string `json:"code"`
	}
	if json.Unmarshal(body.Error, &result.code) != nil && json.Unmarshal(body.Error, &detail) == nil {
		result.code = detail.Code
		if result.code == "" {
			result.code = detail.Type
		}
	}
	if len(result.code) > 64 || strings.ContainsFunc(result.code, func(r rune) bool {
		return !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '_' || r == '-' || r == '.')
	}) {
		result.code = ""
	}
	return result
}

// subscription holds a Claude plan's billing anchor, which only the profile endpoint reports.
type subscription struct {
	anchor  time.Time
	ends    bool
	checked time.Time
}

// readSubscription rereads the Claude profile at most daily, keeping the last answer on failure.
func (m *Maintenance) readSubscription(ctx context.Context, account *coreauth.Auth) {
	m.subMu.Lock()
	sub := m.subs[account.ID]
	m.subMu.Unlock()
	if time.Since(sub.checked) < 24*time.Hour {
		return
	}
	sub.checked = time.Now()
	var profile struct {
		Organization struct {
			Status  string     `json:"subscription_status"`
			Created *time.Time `json:"subscription_created_at"`
		} `json:"organization"`
	}
	if data, err := m.Call(ctx, account, claudeProfileURL, nil); err == nil && json.Unmarshal(data, &profile) == nil && profile.Organization.Created != nil {
		sub.anchor, sub.ends = *profile.Organization.Created, profile.Organization.Status != "active"
	}
	m.subMu.Lock()
	m.subs[account.ID] = sub
	m.subMu.Unlock()
}

// Renewal estimates the next monthly billing date and whether the plan ends there instead of renewing.
// Codex reports the current period's end in its ID token, which can trail renewals, so both
// providers roll their last known billing date forward by whole months.
func (m *Maintenance) Renewal(account *coreauth.Auth, now time.Time) (*time.Time, bool) {
	var anchor time.Time
	ends := false
	if account.Provider == "codex" {
		token, _ := account.Metadata["id_token"].(string)
		claims, err := codex.ParseJWTToken(token)
		if err != nil {
			return nil, false
		}
		until, _ := claims.CodexAuthInfo.ChatgptSubscriptionActiveUntil.(string)
		if anchor, err = time.Parse(time.RFC3339, until); err != nil {
			return nil, false
		}
	} else {
		m.subMu.Lock()
		anchor, ends = m.subs[account.ID].anchor, m.subs[account.ID].ends
		m.subMu.Unlock()
		if anchor.IsZero() {
			return nil, false
		}
	}
	for months := 0; ; months++ {
		next := anchor.AddDate(0, months, 0)
		// Billing on the 31st falls on the last day of shorter months.
		if next.Day() != anchor.Day() {
			next = next.AddDate(0, 0, -next.Day())
		}
		if next.After(now) {
			return &next, ends
		}
	}
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
