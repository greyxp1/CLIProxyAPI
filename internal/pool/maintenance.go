package pool

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	coreauth "github.com/router-for-me/CLIProxyAPI/v8/sdk/cliproxy/auth"
	log "github.com/sirupsen/logrus"
)

type Maintenance struct {
	manager *coreauth.Manager
	db      *sql.DB
	quotas  *coreauth.QuotaStore
	resetMu sync.Mutex
}

func NewMaintenance(manager *coreauth.Manager, db *sql.DB, quotas *coreauth.QuotaStore) (*Maintenance, error) {
	if manager == nil || db == nil || quotas == nil {
		return nil, fmt.Errorf("pool: maintenance dependencies unavailable")
	}
	for _, statement := range []string{
		`CREATE TABLE IF NOT EXISTS pool_accounts (account TEXT PRIMARY KEY, not_before INTEGER NOT NULL DEFAULT 0, short_reset INTEGER, checked_at INTEGER NOT NULL DEFAULT 0)`,
		`CREATE TABLE IF NOT EXISTS pool_resets (account TEXT NOT NULL, credit TEXT NOT NULL, request_id TEXT NOT NULL, outcome TEXT NOT NULL DEFAULT '', cooldown_cleared INTEGER NOT NULL DEFAULT 0, PRIMARY KEY(account,credit))`,
	} {
		if _, err := db.Exec(statement); err != nil {
			return nil, fmt.Errorf("pool: initialize maintenance ledger: %w", err)
		}
	}
	return &Maintenance{manager: manager, db: db, quotas: quotas}, nil
}

func (m *Maintenance) Run(ctx context.Context) {
	for {
		m.cycle(ctx)
		timer := time.NewTimer(m.nextWake(time.Now()))
		select {
		case <-ctx.Done():
			timer.Stop()
			return
		case <-timer.C:
		}
	}
}

// nextWake runs maintenance every five minutes, and every minute once a banked reset is due,
// so its short redemption window is never skipped and a failed redemption is retried.
func (m *Maintenance) nextWake(now time.Time) time.Duration {
	wait := 5 * time.Minute
	for _, account := range m.manager.List() {
		snapshot, ok := m.quotas.Get(account.ID)
		if ok && snapshot.CreditExpiry != nil && snapshot.CreditExpiry.After(now) {
			wait = min(wait, max(snapshot.CreditExpiry.Add(-coreauth.CreditRedeemWindow).Sub(now), time.Minute))
		}
	}
	return wait
}

func (m *Maintenance) cycle(ctx context.Context) {
	for _, account := range m.manager.List() {
		if ctx.Err() != nil {
			return
		}
		if account.Disabled || (account.Provider != "claude" && account.Provider != "codex") || coreauth.IsConfigAPIKeyAuth(account) {
			continue
		}
		digest := sha256.Sum256([]byte(account.ID))
		entry := log.WithFields(log.Fields{"account": hex.EncodeToString(digest[:6]), "provider": account.Provider})
		snapshot, err := m.maintainAccount(ctx, account)
		if err != nil {
			if ctx.Err() == nil {
				entry.Warn("pool account maintenance failed")
			}
			continue
		}
		if errWarm := m.warmIfIdle(ctx, account, snapshot); errWarm != nil && ctx.Err() == nil {
			entry.Warn("pool warmup failed")
		}
	}
}

func (m *Maintenance) maintainAccount(ctx context.Context, account *coreauth.Auth) (coreauth.QuotaSnapshot, error) {
	if account.Provider == "codex" {
		m.resetMu.Lock()
		defer m.resetMu.Unlock()
	}
	snapshot, err := m.Probe(ctx, account)
	if err != nil || account.Provider != "codex" {
		return snapshot, err
	}
	credits, err := m.creditList(ctx, account)
	if err != nil {
		return snapshot, err
	}
	snapshot = m.noteCredits(account.ID, snapshot, credits)
	// Resets are only spent just before they expire. The owner may be draining accounts
	// ahead of an announced reset, so an empty account is never a reason to spend one.
	changed, err := m.redeemExpiring(ctx, account, snapshot, credits)
	if err != nil || !changed {
		return snapshot, err
	}
	return m.refreshResetState(ctx, account)
}

func (m *Maintenance) refreshResetState(ctx context.Context, account *coreauth.Auth) (coreauth.QuotaSnapshot, error) {
	snapshot, err := m.Probe(ctx, account)
	if err != nil {
		return snapshot, err
	}
	credits, err := m.creditList(ctx, account)
	if err != nil {
		return snapshot, err
	}
	return m.noteCredits(account.ID, snapshot, credits), nil
}

func (m *Maintenance) warmIfIdle(ctx context.Context, account *coreauth.Auth, snapshot coreauth.QuotaSnapshot) error {
	if snapshot.Short == nil {
		return nil
	}
	now := time.Now()
	model := "gpt-6.1-sol"
	if account.Provider == "claude" {
		model = "claude-haiku-4-5-20251001"
	}
	_, exhausted, _ := m.quotas.Pressure(account.ID, model, now)
	if exhausted {
		return nil
	}
	var notBefore, checked int64
	var previousReset sql.NullInt64
	err := m.db.QueryRowContext(ctx, `SELECT not_before,short_reset,checked_at FROM pool_accounts WHERE account=?`, account.ID).Scan(&notBefore, &previousReset, &checked)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	short := snapshot.Short
	var reset any
	if short.Reset != nil {
		reset = short.Reset.Unix()
	}
	// A known natural reset releases a prior conservative five-hour warmup claim.
	if notBefore > now.Unix() && short.Reset != nil && short.Reset.After(now) && short.Reset.Unix() < notBefore {
		notBefore = short.Reset.Unix()
	}
	ready := short.Reset == nil || !short.Reset.After(now)
	if short.Used == 0 && short.Reset != nil && previousReset.Valid && short.Reset.Unix()-previousReset.Int64 >= 60 && now.Unix()-checked <= 900 {
		ready = true
	}
	ready = ready && notBefore <= now.Unix()
	if ready {
		notBefore = now.Add(5 * time.Hour).Unix()
	}
	if _, errSave := m.db.ExecContext(ctx, `INSERT INTO pool_accounts(account,not_before,short_reset,checked_at) VALUES(?,?,?,?) ON CONFLICT(account) DO UPDATE SET not_before=excluded.not_before,short_reset=excluded.short_reset,checked_at=excluded.checked_at`, account.ID, notBefore, reset, now.Unix()); errSave != nil {
		return errSave
	}
	if !ready {
		return nil
	}
	var target string
	var body any
	if account.Provider == "claude" {
		target = "https://api.anthropic.com/v1/messages"
		body = map[string]any{"model": model, "max_tokens": 1, "messages": []any{map[string]any{"role": "user", "content": "Hi"}}}
	} else {
		target = "https://chatgpt.com/backend-api/codex/responses"
		body = map[string]any{"model": model, "instructions": "Reply with OK.", "input": []any{map[string]any{"role": "user", "content": []any{map[string]any{"type": "input_text", "text": "Hi"}}}}, "reasoning": map[string]any{"effort": "low"}, "store": false, "stream": true}
	}
	data, errCall := m.Call(ctx, account, target, body)
	if errCall != nil {
		return errCall
	}
	if account.Provider == "claude" {
		var response struct {
			Type string `json:"type"`
		}
		if json.Unmarshal(data, &response) != nil || response.Type != "message" {
			return fmt.Errorf("pool: warmup incomplete")
		}
	} else {
		completed := false
		for _, line := range strings.Split(string(data), "\n") {
			if !strings.HasPrefix(line, "data: ") {
				continue
			}
			var event struct {
				Type string `json:"type"`
			}
			if json.Unmarshal([]byte(strings.TrimPrefix(line, "data: ")), &event) == nil && event.Type == "response.completed" {
				completed = true
			}
		}
		if !completed {
			return fmt.Errorf("pool: warmup incomplete")
		}
	}
	return nil
}

var resetNamespace = uuid.UUID{0x6f, 0x1c, 0x2a, 0x9e, 0x2d, 0x3d, 0x40, 0xc4, 0x9a, 0x9a, 0x7f, 0x3b, 0x5e, 0x0c, 0x1a, 0x42}

// redeemExpiring redeems a banked reset in its last minutes.
func (m *Maintenance) redeemExpiring(ctx context.Context, account *coreauth.Auth, snapshot coreauth.QuotaSnapshot, credits []resetCredit) (bool, error) {
	var credit string
	// Finish a previously redeemed credit's cooldown before considering another.
	err := m.db.QueryRowContext(ctx, `SELECT credit FROM pool_resets WHERE account=? AND (outcome='' OR (outcome IN ('reset','already_redeemed') AND cooldown_cleared=0)) ORDER BY credit LIMIT 1`, account.ID).Scan(&credit)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	if err == sql.ErrNoRows {
		used := false
		if snapshot.Short != nil && snapshot.Short.Used > 0 {
			used = true
		}
		for _, window := range snapshot.Long {
			if window.Used > 0 {
				used = true
			}
		}
		if !used {
			return false, nil
		}
		now := time.Now()
		var earliest time.Time
		for _, candidate := range credits {
			expires := candidate.expires
			if !expires.After(now) || expires.Sub(now) > coreauth.CreditRedeemWindow {
				continue
			}
			var previous string
			errPrevious := m.db.QueryRowContext(ctx, `SELECT outcome FROM pool_resets WHERE account=? AND credit=?`, account.ID, candidate.ID).Scan(&previous)
			if errPrevious != nil && errPrevious != sql.ErrNoRows {
				return false, errPrevious
			}
			if previous == "reset" || previous == "already_redeemed" || previous == "no_credit" {
				continue
			}
			if earliest.IsZero() || expires.Before(earliest) {
				credit = candidate.ID
				earliest = expires
			}
		}
		if credit == "" {
			return false, nil
		}
	}
	outcome, err := m.redeemCredit(ctx, account, credit)
	return outcome == "reset" || outcome == "already_redeemed", err
}

// redeemCredit is shared by manual and automatic redemption under resetMu.
// The durable claim and deterministic ID make lost responses safe to retry.
func (m *Maintenance) redeemCredit(ctx context.Context, account *coreauth.Auth, credit string) (string, error) {
	var requestID, outcome string
	var cleared int
	err := m.db.QueryRowContext(ctx, `SELECT request_id,outcome,cooldown_cleared FROM pool_resets WHERE account=? AND credit=?`, account.ID, credit).Scan(&requestID, &outcome, &cleared)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}
	if err == sql.ErrNoRows {
		identity, _ := account.Metadata["account_id"].(string)
		if identity == "" {
			identity = account.ID
		}
		requestID = uuid.NewSHA1(resetNamespace, []byte(identity+":"+credit)).String()
		if _, errClaim := m.db.ExecContext(ctx, `INSERT INTO pool_resets(account,credit,request_id) VALUES(?,?,?)`, account.ID, credit, requestID); errClaim != nil {
			return "", errClaim
		}
	}
	if outcome == "no_credit" {
		return outcome, nil
	}
	if (outcome == "reset" || outcome == "already_redeemed") && cleared != 0 {
		return "already_redeemed", nil
	}
	if outcome == "nothing_to_reset" {
		if _, errRetry := m.db.ExecContext(ctx, `UPDATE pool_resets SET outcome='' WHERE account=? AND credit=?`, account.ID, credit); errRetry != nil {
			return "", errRetry
		}
		outcome = ""
	}
	if outcome == "" {
		data, errConsume := m.Call(ctx, account, creditURL+"/consume", map[string]any{"credit_id": credit, "redeem_request_id": requestID})
		if errConsume != nil {
			return "", errConsume
		}
		var response struct {
			Code string `json:"code"`
		}
		if json.Unmarshal(data, &response) != nil {
			return "", fmt.Errorf("pool: invalid credit response")
		}
		switch response.Code {
		case "reset", "already_redeemed", "nothing_to_reset", "no_credit":
			outcome = response.Code
		default:
			return "", fmt.Errorf("pool: unexpected credit response")
		}
		if _, errRecord := m.db.ExecContext(ctx, `UPDATE pool_resets SET outcome=? WHERE account=? AND credit=?`, outcome, account.ID, credit); errRecord != nil {
			return outcome, errRecord
		}
	}
	if outcome != "reset" && outcome != "already_redeemed" {
		return outcome, nil
	}
	if _, _, errReset := m.manager.ResetQuota(ctx, account.ID); errReset != nil {
		return outcome, errReset
	}
	tx, errBegin := m.db.BeginTx(ctx, nil)
	if errBegin != nil {
		return outcome, errBegin
	}
	defer tx.Rollback()
	if _, errUpdate := tx.ExecContext(ctx, `UPDATE pool_resets SET cooldown_cleared=1 WHERE account=? AND credit=?`, account.ID, credit); errUpdate != nil {
		return outcome, errUpdate
	}
	if _, errUpdate := tx.ExecContext(ctx, `UPDATE pool_accounts SET not_before=0 WHERE account=?`, account.ID); errUpdate != nil {
		return outcome, errUpdate
	}
	if errCommit := tx.Commit(); errCommit != nil {
		return outcome, errCommit
	}
	return outcome, nil
}

type resetCredit struct {
	ID      string
	expires time.Time
}

// creditList returns the account's available banked resets.
func (m *Maintenance) creditList(ctx context.Context, account *coreauth.Auth) ([]resetCredit, error) {
	var list struct {
		Credits []struct {
			ID        string `json:"id"`
			Status    string `json:"status"`
			ResetType string `json:"reset_type"`
			ExpiresAt string `json:"expires_at"`
		} `json:"credits"`
	}
	data, err := m.Credits(ctx, account)
	if err != nil {
		return nil, err
	}
	if json.Unmarshal(data, &list) != nil {
		return nil, fmt.Errorf("pool: invalid credit list")
	}
	credits := make([]resetCredit, 0, len(list.Credits))
	for _, candidate := range list.Credits {
		expires, errParse := time.Parse(time.RFC3339Nano, candidate.ExpiresAt)
		if errParse == nil && candidate.Status == "available" && candidate.ResetType == "codex_rate_limits" && expires.After(time.Now()) {
			credits = append(credits, resetCredit{ID: candidate.ID, expires: expires})
		}
	}
	return credits, nil
}

// noteCredits records banked resets on the snapshot so routing treats their expiry as a deadline.
func (m *Maintenance) noteCredits(id string, snapshot coreauth.QuotaSnapshot, credits []resetCredit) coreauth.QuotaSnapshot {
	snapshot.Credits, snapshot.CreditExpiry = len(credits), nil
	snapshot.ResetCredits = make([]coreauth.ResetCredit, 0, len(credits))
	for _, credit := range credits {
		snapshot.ResetCredits = append(snapshot.ResetCredits, coreauth.ResetCredit{ID: credit.ID, ExpiresAt: credit.expires})
		if snapshot.CreditExpiry == nil || credit.expires.Before(*snapshot.CreditExpiry) {
			expires := credit.expires
			snapshot.CreditExpiry = &expires
		}
	}
	sort.Slice(snapshot.ResetCredits, func(i, j int) bool {
		a, b := snapshot.ResetCredits[i], snapshot.ResetCredits[j]
		if !a.ExpiresAt.Equal(b.ExpiresAt) {
			return a.ExpiresAt.Before(b.ExpiresAt)
		}
		return a.ID < b.ID
	})
	m.quotas.Set(id, snapshot)
	return snapshot
}
