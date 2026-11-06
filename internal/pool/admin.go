package pool

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"sort"
	"time"

	coreauth "github.com/router-for-me/CLIProxyAPI/v8/sdk/cliproxy/auth"
)

// authorizeOwner reuses the management key check, including its failed-attempt bans.
func (p *Pool) authorizeOwner(w http.ResponseWriter, r *http.Request) bool {
	allowed, status, message := p.oauth.AuthenticateManagementKey(clientIP(r), false, bearer(r, false))
	if !allowed {
		jsonReply(w, status, map[string]string{"error": message})
	}
	return allowed
}

type adminMember struct {
	Name               string           `json:"name"`
	Requests           int64            `json:"requests"`
	Total              int64            `json:"total"`
	TodayRequests      int64            `json:"today_requests"`
	TodayTotal         int64            `json:"today_total"`
	LastUsed           int64            `json:"last_used,omitempty"`
	LastUsedByProvider map[string]int64 `json:"last_used_by_provider,omitempty"`
	Managed            bool             `json:"managed"`
	// Days holds each provider's tokens per usage day, from the pool's first usage day to today.
	Days map[string][]int64 `json:"days"`
}

type adminWindow struct {
	Name  string     `json:"name"`
	Used  float64    `json:"used"`
	Reset *time.Time `json:"reset,omitempty"`
}

type adminAccount struct {
	ID           string                 `json:"id"`
	Provider     string                 `json:"provider"`
	Plan         string                 `json:"plan,omitempty"`
	Email        string                 `json:"email,omitempty"`
	Name         string                 `json:"name,omitempty"`
	State        string                 `json:"state"`
	Next         bool                   `json:"next"`
	Score        *float64               `json:"score,omitempty"`
	Windows      []adminWindow          `json:"windows"`
	CheckedAt    *time.Time             `json:"checked_at,omitempty"`
	RetryAt      *time.Time             `json:"retry_at,omitempty"`
	LastUsed     int64                  `json:"last_used,omitempty"`
	LastMembers  []string               `json:"last_members,omitempty"`
	Credits      int                    `json:"credits,omitempty"`
	CreditExpiry *time.Time             `json:"credit_expiry,omitempty"`
	ResetCredits []coreauth.ResetCredit `json:"reset_credits,omitempty"`
	// Members estimates how much of the current weekly limit each member used, in percent,
	// by splitting it by their tokens on the account since the window began.
	Members map[string]float64 `json:"members,omitempty"`
}

func (p *Pool) memberName(names map[string]string, client string) string {
	if name, ok := names[client]; ok {
		return name
	}
	return "Key " + client[:6]
}

// statusData lists accounts in routing order with their limits and recent use, and every
// member's usage.
func (p *Pool) statusData(role string) (any, error) {
	clients, err := p.clients.Snapshot()
	if err != nil {
		return nil, err
	}
	names := make(map[string]string, len(clients))
	rows, err := p.db.Query(`SELECT client,name FROM clients`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var client, name string
		if err = rows.Scan(&client, &name); err != nil {
			rows.Close()
			return nil, err
		}
		names[client] = name
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return nil, err
	}
	active := make(map[string]bool, len(clients))
	for name, client := range clients {
		id := Fingerprint(client.Inference)
		names[id], active[id] = name, true
	}
	today := usageDay(time.Now())
	members := make(map[string]*adminMember)
	member := func(client string) *adminMember {
		if members[client] == nil {
			members[client] = &adminMember{Name: p.memberName(names, client), Managed: active[client]}
		}
		return members[client]
	}
	for client := range active {
		member(client)
	}
	rows, err = p.db.Query(`SELECT client, SUM(requests), SUM(total), SUM(CASE WHEN day=? THEN requests ELSE 0 END), SUM(CASE WHEN day=? THEN total ELSE 0 END) FROM usage GROUP BY client`, today, today)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var client string
		var item adminMember
		if err = rows.Scan(&client, &item.Requests, &item.Total, &item.TodayRequests, &item.TodayTotal); err != nil {
			rows.Close()
			return nil, err
		}
		target := member(client)
		item.Name, item.Managed = target.Name, target.Managed
		*target = item
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return nil, err
	}

	start, _ := time.Parse("2006-01-02", today)
	var first sql.NullString
	if err = p.db.QueryRow(`SELECT MIN(day) FROM usage`).Scan(&first); err != nil {
		return nil, err
	}
	length := 1
	if parsed, errDay := time.Parse("2006-01-02", first.String); errDay == nil && parsed.Before(start) {
		length += int(start.Sub(parsed).Hours() / 24)
	}
	history := func(item *adminMember) map[string][]int64 {
		if item.Days == nil {
			item.Days = map[string][]int64{"claude": make([]int64, length), "codex": make([]int64, length)}
		}
		return item.Days
	}
	for _, item := range members {
		history(item)
	}
	// The pool only serves Claude and Codex, and every Claude model name starts with claude.
	rows, err = p.db.Query(`SELECT client, day, CASE WHEN model LIKE 'claude%' THEN 'claude' ELSE 'codex' END, SUM(total) FROM usage GROUP BY 1, 2, 3`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var client, day, provider string
		var total int64
		if err = rows.Scan(&client, &day, &provider, &total); err != nil {
			rows.Close()
			return nil, err
		}
		if parsed, errDay := time.Parse("2006-01-02", day); errDay == nil {
			if index := length - 1 - int(start.Sub(parsed).Hours()/24); index >= 0 && index < length {
				history(member(client))[provider][index] += total
			}
		}
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return nil, err
	}

	type dayUsage struct {
		day    string
		client string
		total  int64
	}
	usage := map[string][]dayUsage{}
	lastUsed := map[string]int64{}
	memberLast := map[string]map[string]int64{}
	rows, err = p.db.Query(`SELECT account, client, day, total, last_used FROM account_usage`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var account string
		var item dayUsage
		var last int64
		if err = rows.Scan(&account, &item.client, &item.day, &item.total, &last); err != nil {
			rows.Close()
			return nil, err
		}
		usage[account] = append(usage[account], item)
		lastUsed[account] = max(lastUsed[account], last)
		if memberLast[account] == nil {
			memberLast[account] = map[string]int64{}
		}
		name := p.memberName(names, item.client)
		memberLast[account][name] = max(memberLast[account][name], last)
		target := member(item.client)
		target.LastUsed = max(target.LastUsed, last)
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return nil, err
	}
	rows, err = p.db.Query(`SELECT client,provider,last_used FROM provider_usage`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var client, provider string
		var last int64
		if err = rows.Scan(&client, &provider, &last); err != nil {
			rows.Close()
			return nil, err
		}
		target := member(client)
		if target.LastUsedByProvider == nil {
			target.LastUsedByProvider = map[string]int64{}
		}
		target.LastUsedByProvider[provider] = last
		target.LastUsed = max(target.LastUsed, last)
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return nil, err
	}

	memberList := make([]*adminMember, 0, len(members))
	for _, item := range members {
		memberList = append(memberList, item)
	}
	sort.Slice(memberList, func(i, j int) bool { return memberList[i].Name < memberList[j].Name })
	accountNames := map[string]string{}
	rows, err = p.db.Query(`SELECT account,name FROM account_names`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var id, name string
		if err = rows.Scan(&id, &name); err != nil {
			rows.Close()
			return nil, err
		}
		accountNames[id] = name
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	accountList := make([]*adminAccount, 0)
	for _, state := range p.accountStates(now) {
		auth := state.auth
		item := &adminAccount{ID: auth.EnsureIndex(), Provider: auth.Provider, State: state.state, Next: state.next, Score: state.score, RetryAt: state.retryAt, Windows: []adminWindow{}, LastUsed: lastUsed[auth.ID], LastMembers: recentMembers(memberLast[auth.ID], now)}
		item.Email, _ = auth.Metadata["email"].(string)
		item.Name = accountNames[auth.ID]
		if snapshot, ok := p.quotas.Get(auth.ID); ok {
			item.Plan = snapshot.Plan
			checked := snapshot.CheckedAt
			item.CheckedAt = &checked
			item.Credits, item.CreditExpiry = snapshot.Credits, snapshot.CreditExpiry
			item.ResetCredits = snapshot.ResetCredits
			if snapshot.Short != nil {
				item.Windows = append(item.Windows, adminWindow{Name: "five_hour", Used: snapshot.Short.Used, Reset: snapshot.Short.Reset})
			}
			for _, window := range snapshot.Long {
				name := window.Name
				if name == "" {
					name = "seven_day"
				}
				if name == "seven_day" && window.Reset != nil && window.Used > 0 {
					length := window.Duration
					if length <= 0 {
						length = 7 * 24 * time.Hour
					}
					day := usageDay(window.Reset.Add(-length))
					tokens := map[string]int64{}
					var sum int64
					for _, entry := range usage[auth.ID] {
						if entry.day >= day {
							tokens[entry.client] += entry.total
							sum += entry.total
						}
					}
					if sum > 0 {
						item.Members = map[string]float64{}
						for client, total := range tokens {
							item.Members[p.memberName(names, client)] += window.Used * float64(total) / float64(sum)
						}
					}
				}
				item.Windows = append(item.Windows, adminWindow{Name: name, Used: window.Used, Reset: window.Reset})
			}
		}
		accountList = append(accountList, item)
	}
	data := map[string]any{"role": role, "members": memberList, "accounts": accountList, "today": today}
	if role == "owner" {
		models, err := p.models(context.Background())
		if err != nil {
			return nil, err
		}
		options, err := p.poolOptions(context.Background())
		if err != nil {
			return nil, err
		}
		data["models"] = models
		data["options"] = adminOptionList(options, models)
	}
	data["server_time"] = time.Now().UnixMilli()
	return data, nil
}

type accountState struct {
	auth    *coreauth.Auth
	state   string
	next    bool
	score   *float64
	retryAt *time.Time
	// back is when an account that can't take requests is expected to recover.
	back time.Time
}

// accountStates lists pooled accounts in routing order and marks each provider's next pick.
func (p *Pool) accountStates(now time.Time) []*accountState {
	states := make([]*accountState, 0)
	best := map[string]*accountState{}
	for _, auth := range p.manager.List() {
		if (auth.Provider != "claude" && auth.Provider != "codex") || coreauth.IsConfigAPIKeyAuth(auth) {
			continue
		}
		item := &accountState{auth: auth, state: "unknown quota"}
		// Model-specific weekly limits can change the choice for that model.
		score, exhausted, fresh := p.quotas.Pressure(auth.ID, "", now)
		snapshot, _ := p.quotas.Get(auth.ID)
		// Blocked accounts are exhausted only when their week is spent. Otherwise they wait on a shorter limit.
		blocked := "cooling down"
		if weekSpent(snapshot, now) {
			blocked = "exhausted"
		}
		switch {
		case auth.Disabled:
			item.state = "disabled"
		case auth.Unavailable && auth.NextRetryAfter.After(now):
			retry := auth.NextRetryAfter
			item.state, item.retryAt, item.back = blocked, &retry, retry
		case exhausted:
			item.state, item.back = blocked, exhaustedUntil(snapshot)
		case fresh:
			item.state, item.score = "available", &score
			if best[auth.Provider] == nil || score > *best[auth.Provider].score || (score == *best[auth.Provider].score && auth.ID < best[auth.Provider].auth.ID) {
				best[auth.Provider] = item
			}
		}
		states = append(states, item)
	}
	for _, item := range best {
		item.next = true
	}
	sort.Slice(states, func(i, j int) bool {
		a, b := states[i], states[j]
		if a.auth.Provider != b.auth.Provider {
			return a.auth.Provider < b.auth.Provider
		}
		if (a.score == nil) != (b.score == nil) {
			return a.score != nil
		}
		if a.score != nil && *a.score != *b.score {
			return *a.score > *b.score
		}
		// Accounts that can't take requests follow, the soonest to recover first.
		if rank := map[string]int{"unknown quota": 0, "cooling down": 1, "exhausted": 2, "disabled": 3}; rank[a.state] != rank[b.state] {
			return rank[a.state] < rank[b.state]
		}
		if !a.back.Equal(b.back) {
			return a.back.Before(b.back)
		}
		return a.auth.ID < b.auth.ID
	})
	return states
}

// weekSpent reports whether the weekly window that applies to every model is used up.
func weekSpent(snapshot coreauth.QuotaSnapshot, now time.Time) bool {
	for _, window := range snapshot.Long {
		if (window.Name == "" || window.Name == "seven_day") && window.Used >= 100 && (window.Reset == nil || window.Reset.After(now)) {
			return true
		}
	}
	return false
}

// exhaustedUntil is when the last full window that applies to every model resets.
func exhaustedUntil(snapshot coreauth.QuotaSnapshot) (until time.Time) {
	full := func(window coreauth.QuotaWindow) {
		if window.Used >= 100 && window.Reset != nil && window.Reset.After(until) {
			until = *window.Reset
		}
	}
	for _, window := range snapshot.Long {
		if window.Name == "" || window.Name == "seven_day" {
			full(window)
		}
	}
	if snapshot.Short != nil {
		full(*snapshot.Short)
	}
	return until
}

func (p *Pool) adminClients(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Action  string `json:"action"`
		Name    string `json:"name"`
		NewName string `json:"new_name"`
	}
	if !decodeBody(w, r, &body) {
		return
	}
	if body.Action == "show" {
		clients, err := p.clients.Snapshot()
		client, exists := clients[body.Name]
		if err != nil || !exists {
			jsonReply(w, 400, map[string]string{"error": "unknown client"})
			return
		}
		jsonReply(w, 200, map[string]string{"inference": client.Inference, "prompt": setupPrompt(p.cfg.Origin, client)})
		return
	}
	var created Client
	err := UpdateClients(p.cfg.ClientsFile, func(clients map[string]Client) error {
		switch body.Action {
		case "add":
			var err error
			created, err = AddClient(clients, body.Name)
			return err
		case "rename":
			return RenameClient(clients, body.Name, body.NewName)
		case "revoke":
			return RevokeClient(clients, body.Name)
		}
		return fmt.Errorf("unknown action")
	})
	if err != nil {
		jsonReply(w, 400, map[string]string{"error": err.Error()})
		return
	}
	if errSync := p.syncUsage(); errSync != nil {
		jsonReply(w, 500, map[string]string{"error": "Saved, but usage names could not be updated"})
		return
	}
	if body.Action != "add" {
		jsonReply(w, 200, map[string]string{})
		return
	}
	jsonReply(w, 200, map[string]string{"prompt": setupPrompt(p.cfg.Origin, created)})
}

// setupPrompt fits in one Discord message. The proxy enforces effort and fast-mode limits itself.
func setupPrompt(origin string, client Client) string {
	return fmt.Sprintf(`Set up Claude Code and Codex to use my friend's shared AI proxy. Preserve unrelated settings and chat history. Never print or commit the key.

URL: %[1]s
Key: %[2]s

1. Install missing CLIs. If Nix or another tool manages ~/.claude/settings.json or ~/.codex/config.toml, edit its source.
2. Claude Code: in settings.json "env", set ANTHROPIC_BASE_URL to the URL and ANTHROPIC_AUTH_TOKEN to the key. Remove conflicting ANTHROPIC_API_KEY, apiKeyHelper and proxy overrides, checking shell profiles too.
3. Codex: set top-level model_provider = "cliproxy" in config.toml and replace [model_providers.cliproxy] with:
[model_providers.cliproxy]
name = "CLIProxyAPI"
base_url = "%[1]s/v1"
model_catalog_url = "%[1]s/v1/models?provider=codex"
wire_api = "responses"
requires_openai_auth = false
supports_websockets = true
experimental_bearer_token = "%[2]s"
4. T3 Code: verify official pingdotgg/t3code works. Repair or replace wrong or broken installs, preserving data. Windows: use the desktop app (winget install T3Tools.T3Code if missing). Stop any standalone T3 server sharing its data. Open the app, enable Local environment and T3 Connect, help me sign in and leave it minimized. No Windows autostart. Linux/macOS: install the CLI if missing with curl -fsSL https://t3.codes/install.sh | sh, run t3 connect, help me sign in and verify service autostart. Open https://app.t3.codes and verify this machine connects. Use Claude Code for Claude and Codex ("Use existing CLI") for OpenAI. Hide models absent from the proxy catalog.
5. Check %[1]s/v1/models with Bearer auth and send one short request through each CLI using listed models. Distinguish exhausted usage and cooldowns from authentication or configuration errors.
6. Open %[1]s/proxy for usage and limits. Report changes and test results. End with that link.
`, origin, client.Inference)
}

// recentMembers lists who used an account in the last ten minutes, most recent first, or
// whoever used it last when nobody has since.
func recentMembers(last map[string]int64, now time.Time) []string {
	names := make([]string, 0, len(last))
	for name := range last {
		names = append(names, name)
	}
	sort.Slice(names, func(i, j int) bool { return last[names[i]] > last[names[j]] })
	cutoff := now.Add(-10 * time.Minute).Unix()
	for i, name := range names {
		if last[name] < cutoff {
			return names[:max(i, 1)]
		}
	}
	return names
}
