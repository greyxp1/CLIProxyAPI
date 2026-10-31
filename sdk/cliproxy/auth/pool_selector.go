package auth

import (
	"context"
	"math"
	"strings"
	"sync"
	"time"

	cliproxyexecutor "github.com/router-for-me/CLIProxyAPI/v8/sdk/cliproxy/executor"
)

type QuotaWindow struct {
	Used     float64
	Reset    *time.Time
	Duration time.Duration
	Name     string
}

type ResetCredit struct {
	ID        string    `json:"id"`
	ExpiresAt time.Time `json:"expires_at"`
}

type QuotaSnapshot struct {
	Plan      string
	Short     *QuotaWindow
	Long      []QuotaWindow
	CheckedAt time.Time
	// Credits counts banked resets. CreditExpiry is when the soonest one expires.
	Credits      int
	CreditExpiry *time.Time
	ResetCredits []ResetCredit
}

// CreditRedeemWindow is how long before a banked reset expires maintenance redeems it.
const CreditRedeemWindow = 5 * time.Minute

// QuotaStore shares authoritative observations between account maintenance and selection.
// No selector performs network or filesystem work.
type QuotaStore struct {
	mu       sync.RWMutex
	accounts map[string]QuotaSnapshot
}

func NewQuotaStore() *QuotaStore { return &QuotaStore{accounts: make(map[string]QuotaSnapshot)} }

func (s *QuotaStore) Set(id string, snapshot QuotaSnapshot) {
	if s == nil {
		return
	}
	snapshot.Long = append([]QuotaWindow(nil), snapshot.Long...)
	snapshot.ResetCredits = append([]ResetCredit(nil), snapshot.ResetCredits...)
	s.mu.Lock()
	if s.accounts == nil {
		s.accounts = make(map[string]QuotaSnapshot)
	}
	s.accounts[id] = snapshot
	s.mu.Unlock()
}

func (s *QuotaStore) Get(id string) (QuotaSnapshot, bool) {
	if s == nil {
		return QuotaSnapshot{}, false
	}
	s.mu.RLock()
	snapshot, ok := s.accounts[id]
	s.mu.RUnlock()
	snapshot.Long = append([]QuotaWindow(nil), snapshot.Long...)
	snapshot.ResetCredits = append([]ResetCredit(nil), snapshot.ResetCredits...)
	return snapshot, ok
}

// Pressure normalizes each deadline by its window duration, so five-hour resets
// do not automatically outrank weekly resets. It is an urgency heuristic.
func (s *QuotaStore) Pressure(id, model string, now time.Time) (score float64, exhausted, fresh bool) {
	snapshot, ok := s.Get(id)
	if !ok || now.Sub(snapshot.CheckedAt) > 10*time.Minute || snapshot.CheckedAt.Sub(now) > time.Minute {
		return 0, false, false
	}
	windows := applicableQuotaWindows(snapshot, model)
	if len(windows) == 0 {
		return 0, false, false
	}
	for _, window := range windows {
		if math.IsNaN(window.Used) || math.IsInf(window.Used, 0) || window.Used < 0 || window.Duration <= 0 {
			return 0, false, false
		}
		if window.Used >= 100 && (window.Reset == nil || window.Reset.After(now)) {
			return 0, true, true
		}
	}
	// Maintenance redeems banked resets in their last five minutes.
	var credit *time.Time
	if snapshot.CreditExpiry != nil && snapshot.CreditExpiry.After(now) {
		redeemAt := snapshot.CreditExpiry.Add(-CreditRedeemWindow)
		credit = &redeemAt
	}
	for _, window := range windows {
		if window.Reset != nil && !window.Reset.After(now) {
			return 0, false, false
		}
		// Redeeming a banked reset refills the account, so quota left when it expires is lost.
		deadline := window.Reset
		if credit != nil && (deadline == nil || credit.Before(*deadline)) {
			deadline = credit
		}
		if deadline == nil {
			if window.Used > 0 {
				return 0, false, false
			}
			continue
		}
		pressure := (1 - window.Used/100) * window.Duration.Seconds() / math.Max(deadline.Sub(now).Seconds(), 300)
		score = math.Max(score, pressure)
	}
	return score, false, true
}

func applicableQuotaWindows(snapshot QuotaSnapshot, model string) []QuotaWindow {
	model = strings.ToLower(model)
	windows := make([]QuotaWindow, 0, len(snapshot.Long)+1)
	for _, window := range snapshot.Long {
		if window.Name != "" && window.Name != "seven_day" {
			family := strings.TrimPrefix(window.Name, "seven_day_")
			if family == "" || !strings.Contains(model, family) {
				continue
			}
		}
		windows = append(windows, window)
	}
	if snapshot.Short != nil {
		windows = append(windows, *snapshot.Short)
	}
	return windows
}

// QuotaPressureSelector applies only when native session affinity needs a new seat.
// Keep it outside the scheduler fast path, which understands built-in ordering only.
type QuotaPressureSelector struct {
	Quotas     *QuotaStore
	Fallback   Selector
	roundRobin RoundRobinSelector
}

func (s *SessionAffinitySelector) availableFallbackAuths(ctx context.Context, auths []*Auth, provider, model string, now time.Time) ([]*Auth, error) {
	if _, ok := s.fallback.(*QuotaPressureSelector); ok {
		return getSelectorAvailableAuthsAcrossPriorities(ctx, auths, provider, model, now)
	}
	return getSelectorAvailableAuths(ctx, auths, provider, model, now)
}

func (s *SessionAffinitySelector) fallbackPriorityAuths(auths []*Auth) []*Auth {
	if _, ok := s.fallback.(*QuotaPressureSelector); ok {
		return auths
	}
	return highestPriorityAuths(auths)
}

func (s *QuotaPressureSelector) Pick(ctx context.Context, provider, model string, opts cliproxyexecutor.Options, auths []*Auth) (*Auth, error) {
	now := time.Now()
	available, err := getSelectorAvailableAuthsAcrossPriorities(ctx, auths, provider, model, now)
	if err != nil {
		return nil, err
	}
	// Filter authoritative quota before choosing a priority tier. An exhausted
	// high-priority account must not hide an eligible lower-priority account.
	eligible := make([]*Auth, 0, len(available))
	for _, account := range available {
		_, full, _ := s.Quotas.Pressure(account.ID, model, now)
		if !full {
			eligible = append(eligible, account)
		}
	}
	if len(eligible) > 0 {
		available = preferCodexWebsocketAuths(ctx, provider, highestPriorityAuths(eligible))
	}
	var candidates []*Auth
	best := -1.0
	unknown := make([]*Auth, 0, len(available))
	var next time.Time
	for _, account := range available {
		score, full, fresh := s.Quotas.Pressure(account.ID, model, now)
		if full {
			var accountReady time.Time
			unknownReset := false
			if snapshot, ok := s.Quotas.Get(account.ID); ok {
				for _, window := range applicableQuotaWindows(snapshot, model) {
					if window.Used < 100 || (window.Reset != nil && !window.Reset.After(now)) {
						continue
					}
					if window.Reset == nil {
						unknownReset = true
						continue
					}
					if window.Reset.After(accountReady) {
						accountReady = *window.Reset
					}
				}
			}
			if unknownReset {
				accountReady = now.Add(5 * time.Minute)
			}
			if !accountReady.IsZero() && (next.IsZero() || accountReady.Before(next)) {
				next = accountReady
			}
			continue
		}
		if !fresh {
			unknown = append(unknown, account)
			continue
		}
		if score > best {
			best = score
			candidates = []*Auth{account}
		} else if score == best {
			candidates = append(candidates, account)
		}
	}
	if len(candidates) == 0 {
		candidates = unknown
	}
	if len(candidates) == 0 {
		if next.IsZero() {
			next = now.Add(5 * time.Minute)
		}
		return nil, NewModelCooldownError(model, provider, next.Sub(now))
	}
	fallback := s.Fallback
	if fallback == nil {
		fallback = &s.roundRobin
	}
	return fallback.Pick(ctx, provider, model, opts, candidates)
}
