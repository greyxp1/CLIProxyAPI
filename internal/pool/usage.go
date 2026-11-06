package pool

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"time"

	_ "github.com/mattn/go-sqlite3"
	coreusage "github.com/router-for-me/CLIProxyAPI/v8/sdk/cliproxy/usage"
	log "github.com/sirupsen/logrus"
)

func openDatabase(path string) (*sql.DB, error) {
	dsn := (&url.URL{Scheme: "file", Path: path}).String() + "?_busy_timeout=5000&_journal_mode=WAL&_synchronous=FULL"
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	for _, statement := range []string{
		`CREATE TABLE IF NOT EXISTS usage (day TEXT, client TEXT, model TEXT, requests INTEGER, failures INTEGER, input INTEGER, output INTEGER, cached INTEGER, total INTEGER, PRIMARY KEY(day,client,model))`,
		`CREATE TABLE IF NOT EXISTS clients (client TEXT PRIMARY KEY, name TEXT)`,
		`CREATE TABLE IF NOT EXISTS account_names (account TEXT PRIMARY KEY, name TEXT NOT NULL)`,
		`CREATE TABLE IF NOT EXISTS blocked_models (model TEXT PRIMARY KEY)`,
		`CREATE TABLE IF NOT EXISTS allowed_models (model TEXT PRIMARY KEY)`,
		`CREATE TABLE IF NOT EXISTS options (provider TEXT, setting TEXT, enabled INTEGER NOT NULL, PRIMARY KEY(provider,setting))`,
		`CREATE TABLE IF NOT EXISTS account_usage (day TEXT, account TEXT, client TEXT, requests INTEGER, total INTEGER, last_used INTEGER, PRIMARY KEY(day,account,client))`,
		`CREATE TABLE IF NOT EXISTS provider_usage (client TEXT, provider TEXT, last_used INTEGER NOT NULL, PRIMARY KEY(client,provider))`,
	} {
		if _, err = db.Exec(statement); err != nil {
			db.Close()
			return nil, fmt.Errorf("initialize pool database: %w", err)
		}
	}
	return db, nil
}

func (p *Pool) syncUsage() error {
	clients, err := p.clients.Snapshot()
	if err != nil {
		return err
	}
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for name, client := range clients {
		if _, err = tx.Exec(`INSERT INTO clients VALUES (?,?) ON CONFLICT(client) DO UPDATE SET name=excluded.name`, Fingerprint(client.Inference), name); err != nil {
			return err
		}
	}
	// Older usage stores the timestamp against the account, so seed providers while
	// those accounts are still identifiable. New records survive account removal.
	for _, account := range p.manager.List() {
		if account.Provider != "claude" && account.Provider != "codex" {
			continue
		}
		if _, err = tx.Exec(`INSERT INTO provider_usage SELECT client, ?, MAX(last_used) FROM account_usage WHERE account=? GROUP BY client ON CONFLICT(client,provider) DO UPDATE SET last_used=MAX(last_used,excluded.last_used)`, account.Provider, account.ID); err != nil {
			return err
		}
	}
	return tx.Commit()
}

// HandleUsage commits native execution events directly, before the dispatcher forgets them.
func (p *Pool) HandleUsage(_ context.Context, record coreusage.Record) {
	if record.APIKey == "" {
		return
	}
	if err := p.recordUsage(record); err != nil {
		log.WithError(err).Error("could not persist pool usage")
	}
}

// Usage days begin at 08:00 in the server's timezone, including across DST changes.
func usageDay(stamp time.Time) string {
	local := stamp.In(time.Local)
	if local.Hour() < 8 {
		local = local.AddDate(0, 0, -1)
	}
	return local.Format("2006-01-02")
}

func (p *Pool) recordUsage(record coreusage.Record) error {
	stamp := record.RequestedAt
	if stamp.IsZero() {
		stamp = time.Now()
	}
	model := record.Model
	if model == "" {
		model = "unknown"
	}
	clientID := Fingerprint(record.APIKey)
	// Totals do not depend on the registry being readable when an in-flight request finishes.
	clients, _ := p.clients.Snapshot()
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for name, client := range clients {
		if Fingerprint(client.Inference) == clientID {
			if _, err = tx.Exec(`INSERT INTO clients VALUES (?,?) ON CONFLICT(client) DO UPDATE SET name=excluded.name`, clientID, name); err != nil {
				return err
			}
			break
		}
	}
	// Each attempt is recorded, so only successful ones count as requests and a retried request counts once.
	failed, succeeded := 0, 1
	if record.Failed {
		failed, succeeded = 1, 0
	}
	detail := coreusage.EnsureTokenBreakdownForProvider(record.Detail, record.Provider, record.ExecutorType)
	_, err = tx.Exec(`INSERT INTO usage VALUES (?,?,?,?,?,?,?,?,?) ON CONFLICT(day,client,model) DO UPDATE SET requests=requests+excluded.requests, failures=failures+excluded.failures, input=input+excluded.input, output=output+excluded.output, cached=cached+excluded.cached, total=total+excluded.total`,
		usageDay(stamp), clientID, model, succeeded, failed, detail.InputTokens, detail.OutputTokens, detail.CachedTokens, detail.TotalTokens)
	if err != nil {
		return err
	}
	if record.Provider == "claude" || record.Provider == "codex" {
		if _, err = tx.Exec(`INSERT INTO provider_usage VALUES (?,?,?) ON CONFLICT(client,provider) DO UPDATE SET last_used=MAX(last_used,excluded.last_used)`, clientID, record.Provider, stamp.Unix()); err != nil {
			return err
		}
	}
	if record.AuthID != "" {
		_, err = tx.Exec(`INSERT INTO account_usage VALUES (?,?,?,?,?,?) ON CONFLICT(day,account,client) DO UPDATE SET requests=requests+excluded.requests, total=total+excluded.total, last_used=MAX(last_used,excluded.last_used)`,
			usageDay(stamp), record.AuthID, clientID, succeeded, detail.TotalTokens, stamp.Unix())
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}
