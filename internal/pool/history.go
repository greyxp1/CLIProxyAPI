package pool

import (
	"net/http"
	"time"

	coreauth "github.com/router-for-me/CLIProxyAPI/v8/sdk/cliproxy/auth"
)

// Weekly readings are kept two weeks, enough for a full past window beside the forecast.
const historyKeep = 15 * 24 * time.Hour

// recordUsage keeps each weekly reading, so the dashboard can chart the pool over time and
// measure how fast it drains.
func (m *Maintenance) recordUsage(id string, snapshot coreauth.QuotaSnapshot) error {
	for _, window := range snapshot.Long {
		if window.Name != "" && window.Name != "seven_day" {
			continue
		}
		var reset *int64
		if window.Reset != nil {
			at := window.Reset.UnixMilli()
			reset = &at
		}
		if _, err := m.db.Exec(`INSERT OR REPLACE INTO quota_history (account, at, used, reset) VALUES (?, ?, ?, ?)`,
			id, snapshot.CheckedAt.UnixMilli(), window.Used, reset); err != nil {
			return err
		}
	}
	return nil
}

func (m *Maintenance) pruneHistory(now time.Time) error {
	_, err := m.db.Exec(`DELETE FROM quota_history WHERE at < ?`, now.Add(-historyKeep).UnixMilli())
	return err
}

// history lists each pooled account's weekly readings, at most one per half hour, as
// [time, used percent, reset time or null] in milliseconds, keyed like the status accounts.
func (p *Pool) history(w http.ResponseWriter, r *http.Request) {
	if !p.byIP.allow("history:"+clientIP(r), 30, time.Minute) {
		jsonReply(w, 429, map[string]string{"error": "Too many requests"})
		return
	}
	ids := map[string]string{}
	for _, auth := range p.manager.List() {
		if (auth.Provider == "claude" || auth.Provider == "codex") && !coreauth.IsConfigAPIKeyAuth(auth) {
			ids[auth.ID] = auth.EnsureIndex()
		}
	}
	// SQLite fills the bare columns from the row holding MAX(at), so each half hour keeps its latest reading.
	rows, err := p.db.Query(`SELECT account, MAX(at), used, reset FROM quota_history WHERE at >= ? GROUP BY account, at / 1800000 ORDER BY account, at`,
		time.Now().Add(-historyKeep).UnixMilli())
	if err != nil {
		jsonReply(w, 503, map[string]string{"error": "History unavailable"})
		return
	}
	defer rows.Close()
	accounts := map[string][][3]any{}
	for rows.Next() {
		var account string
		var at int64
		var used float64
		var reset *int64
		if err = rows.Scan(&account, &at, &used, &reset); err != nil {
			jsonReply(w, 503, map[string]string{"error": "History unavailable"})
			return
		}
		if id, ok := ids[account]; ok {
			accounts[id] = append(accounts[id], [3]any{at, used, reset})
		}
	}
	if rows.Err() != nil {
		jsonReply(w, 503, map[string]string{"error": "History unavailable"})
		return
	}
	jsonReply(w, 200, map[string]any{"accounts": accounts})
}
