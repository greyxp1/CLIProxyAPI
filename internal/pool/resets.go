package pool

import (
	"database/sql"
	"net/http"

	coreauth "github.com/router-for-me/CLIProxyAPI/v8/sdk/cliproxy/auth"
)

// redeemReset spends the requested banked reset, never a different credit on retry.
// The caller authenticates the owner and enforces the browser origin.
func (p *Pool) redeemReset(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID       string `json:"id"`
		CreditID string `json:"credit_id"`
	}
	if !decodeBody(w, r, &body) {
		return
	}
	if body.ID == "" || body.CreditID == "" {
		jsonReply(w, 400, map[string]string{"error": "Choose an account and a banked reset"})
		return
	}
	if !p.maintenance.resetMu.TryLock() {
		jsonReply(w, 409, map[string]string{"error": "Reset maintenance is in progress. Try again shortly."})
		return
	}
	defer p.maintenance.resetMu.Unlock()
	var account *coreauth.Auth
	for _, candidate := range p.manager.List() {
		if candidate.EnsureIndex() == body.ID && candidate.Provider == "codex" && !coreauth.IsConfigAPIKeyAuth(candidate) {
			account = candidate
			break
		}
	}
	if account == nil {
		jsonReply(w, 404, map[string]string{"error": "Unknown Codex account"})
		return
	}
	if account.Disabled {
		jsonReply(w, 400, map[string]string{"error": "Enable the account before using a reset"})
		return
	}
	// An ambiguous earlier redemption must finish before another credit is spent.
	var pending string
	err := p.db.QueryRowContext(r.Context(), `SELECT credit FROM pool_resets WHERE account=? AND credit<>? AND (outcome='' OR (outcome IN ('reset','already_redeemed') AND cooldown_cleared=0)) LIMIT 1`, account.ID, body.CreditID).Scan(&pending)
	if err != nil && err != sql.ErrNoRows {
		jsonReply(w, 503, map[string]string{"error": "Reset ledger unavailable"})
		return
	}
	if err == nil {
		jsonReply(w, 409, map[string]string{"error": "Another reset is still pending for this account. Try again after maintenance completes."})
		return
	}
	var recorded string
	err = p.db.QueryRowContext(r.Context(), `SELECT request_id FROM pool_resets WHERE account=? AND credit=?`, account.ID, body.CreditID).Scan(&recorded)
	if err != nil && err != sql.ErrNoRows {
		jsonReply(w, 503, map[string]string{"error": "Reset ledger unavailable"})
		return
	}
	if err == sql.ErrNoRows {
		credits, errList := p.maintenance.creditList(r.Context(), account)
		if errList != nil {
			jsonReply(w, 502, map[string]string{"error": "The provider could not return available resets"})
			return
		}
		available := false
		for _, credit := range credits {
			available = available || credit.ID == body.CreditID
		}
		if !available {
			jsonReply(w, 409, map[string]string{"error": "That reset is no longer available. Refresh the account."})
			return
		}
	}
	outcome, errRedeem := p.maintenance.redeemCredit(r.Context(), account, body.CreditID)
	applied := outcome == "reset" || outcome == "already_redeemed"
	if errRedeem != nil && !applied {
		jsonReply(w, 502, map[string]string{"error": "Could not complete the reset. Retry this same credit."})
		return
	}
	refreshed := false
	if errRedeem == nil {
		_, errRefresh := p.maintenance.refreshResetState(r.Context(), account)
		refreshed = errRefresh == nil
	}
	jsonReply(w, 200, map[string]any{"code": outcome, "refreshed": refreshed})
}
