package pool

import (
	"bytes"
	"context"
	"crypto/rand"
	"database/sql"
	"embed"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/router-for-me/CLIProxyAPI/v8/internal/api/handlers/management"
	"github.com/router-for-me/CLIProxyAPI/v8/internal/config"
	"github.com/router-for-me/CLIProxyAPI/v8/sdk/access"
	coreauth "github.com/router-for-me/CLIProxyAPI/v8/sdk/cliproxy/auth"
	coreusage "github.com/router-for-me/CLIProxyAPI/v8/sdk/cliproxy/usage"
	log "github.com/sirupsen/logrus"
)

//go:embed assets/*
var assets embed.FS

type enrollment struct {
	Provider, State string
	Expires         time.Time
}
type rateBucket struct {
	Hits    []time.Time
	Expires time.Time
}

type Pool struct {
	cfg         config.PoolConfig
	manager     *coreauth.Manager
	clients     *ClientRegistry
	db          *sql.DB
	quotas      *coreauth.QuotaStore
	maintenance *Maintenance
	oauth       *management.Handler
	server      *http.Server
	cancel      context.CancelFunc
	workers     sync.WaitGroup
	sessionMu   sync.Mutex
	sessions    map[string]enrollment
	rateMu      sync.Mutex
	rates       map[string]rateBucket
}

func New(cfg *config.Config, manager *coreauth.Manager) (*Pool, error) {
	if cfg == nil || manager == nil {
		return nil, fmt.Errorf("pool configuration and manager are required")
	}
	settings := cfg.Pool
	if settings.Address == "" {
		settings.Address = "127.0.0.1:8318"
	}
	host, _, err := net.SplitHostPort(settings.Address)
	if err != nil || net.ParseIP(host) == nil || !net.ParseIP(host).IsLoopback() {
		return nil, fmt.Errorf("pool address must be a loopback IP and port")
	}
	origin, err := url.Parse(settings.Origin)
	if err != nil || origin.Scheme != "https" || origin.Host == "" || origin.User != nil || origin.RawQuery != "" || origin.Fragment != "" || origin.Path != "" {
		return nil, fmt.Errorf("pool origin must be an HTTPS origin without a path")
	}
	if settings.ClientsFile == "" || settings.Database == "" {
		return nil, fmt.Errorf("pool clients-file and database are required")
	}
	registry := &ClientRegistry{path: settings.ClientsFile}
	if _, err = registry.Snapshot(); err != nil {
		return nil, fmt.Errorf("read pool clients: %w", err)
	}
	db, err := openDatabase(settings.Database)
	if err != nil {
		return nil, err
	}
	quotas := coreauth.NewQuotaStore()
	maintenance, err := NewMaintenance(manager, db, quotas)
	if err != nil {
		db.Close()
		return nil, err
	}
	p := &Pool{cfg: settings, manager: manager, clients: registry, db: db, quotas: quotas, maintenance: maintenance, sessions: make(map[string]enrollment), rates: make(map[string]rateBucket)}
	if err = p.syncClientNames(); err != nil {
		db.Close()
		return nil, fmt.Errorf("initialize client names: %w", err)
	}
	if err = p.applyOptions(context.Background()); err != nil {
		db.Close()
		return nil, fmt.Errorf("load options: %w", err)
	}
	return p, nil
}

func (p *Pool) Registry() access.Provider                   { return p.clients }
func (p *Pool) Quotas() *coreauth.QuotaStore                { return p.quotas }
func (p *Pool) SetOAuthHandler(handler *management.Handler) { p.oauth = handler }

func (p *Pool) Start(ctx context.Context) error {
	if p.oauth == nil {
		return fmt.Errorf("pool OAuth handler is not configured")
	}
	listener, err := net.Listen("tcp", p.cfg.Address)
	if err != nil {
		return fmt.Errorf("listen for pool: %w", err)
	}
	ctx, p.cancel = context.WithCancel(ctx)
	coreusage.RegisterNamedPlugin("pool-usage", p)
	p.server = &http.Server{Handler: http.HandlerFunc(p.serveHTTP), BaseContext: func(net.Listener) context.Context { return ctx }, ReadHeaderTimeout: 10 * time.Second, MaxHeaderBytes: 16384}
	p.workers.Add(1)
	go func() {
		defer p.workers.Done()
		if errServe := p.server.Serve(listener); errServe != nil && !errors.Is(errServe, http.ErrServerClosed) {
			log.WithError(errServe).Error("pool listener stopped")
		}
	}()
	p.workers.Add(1)
	go func() { defer p.workers.Done(); p.watchModels(ctx) }()
	if p.cfg.Maintenance {
		p.workers.Add(1)
		go func() { defer p.workers.Done(); p.maintenance.Run(ctx) }()
	}
	return nil
}

func (p *Pool) Shutdown(ctx context.Context) error {
	if p.cancel != nil {
		p.cancel()
	}
	if p.server != nil {
		if err := p.server.Shutdown(ctx); err != nil {
			return err
		}
	}
	done := make(chan struct{})
	go func() { p.workers.Wait(); close(done) }()
	select {
	case <-done:
		return p.db.Close()
	case <-ctx.Done():
		return ctx.Err()
	}
}

func jsonReply(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func decodeBody(w http.ResponseWriter, r *http.Request, body any) bool {
	if strings.Split(r.Header.Get("Content-Type"), ";")[0] != "application/json" {
		jsonReply(w, 400, map[string]string{"error": "Expected JSON"})
		return false
	}
	r.Body = http.MaxBytesReader(w, r.Body, 16384)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(body); err != nil {
		jsonReply(w, 400, map[string]string{"error": "Invalid request"})
		return false
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		jsonReply(w, 400, map[string]string{"error": "Invalid request"})
		return false
	}
	return true
}

func (p *Pool) limit(key string, count int, duration time.Duration) bool {
	p.rateMu.Lock()
	defer p.rateMu.Unlock()
	now := time.Now()
	for key, bucket := range p.rates {
		if bucket.Expires.Before(now) {
			delete(p.rates, key)
		}
	}
	bucket := p.rates[key]
	if bucket.Hits == nil && len(p.rates) >= 4096 {
		return false
	}
	first := 0
	for first < len(bucket.Hits) && !bucket.Hits[first].After(now.Add(-duration)) {
		first++
	}
	bucket.Hits = bucket.Hits[first:]
	if len(bucket.Hits) >= count {
		return false
	}
	bucket.Hits = append(bucket.Hits, now)
	bucket.Expires = now.Add(duration)
	p.rates[key] = bucket
	return true
}

func (p *Pool) authorize(w http.ResponseWriter, r *http.Request) bool {
	name, err := p.clients.match(bearer(r, true))
	if err != nil {
		jsonReply(w, 503, map[string]string{"error": "Client registry unavailable"})
		return false
	}
	if name == "" {
		jsonReply(w, 401, map[string]string{"error": "Invalid API key"})
		return false
	}
	if !p.limit("inference:"+name, 120, time.Minute) {
		w.Header().Set("Retry-After", "60")
		jsonReply(w, 429, map[string]string{"error": "Too many requests"})
		return false
	}
	return true
}

func (p *Pool) serveHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Referrer-Policy", "no-referrer")
	if r.Method == http.MethodGet {
		files := map[string]string{"/proxy": "proxy.html", "/proxy/": "proxy.html", "/proxy/app.js": "proxy.js", "/proxy/style.css": "style.css", "/proxy/stats.svg": "stats.svg"}
		if file, ok := files[r.URL.Path]; ok {
			data, err := assets.ReadFile("assets/" + file)
			if err != nil {
				http.NotFound(w, r)
				return
			}
			typeName := "text/html; charset=utf-8"
			if strings.HasSuffix(file, ".js") {
				typeName = "text/javascript"
			} else if strings.HasSuffix(file, ".css") {
				typeName = "text/css"
			} else if strings.HasSuffix(file, ".svg") {
				typeName = "image/svg+xml"
			}
			w.Header().Set("Content-Type", typeName)
			w.Header().Set("Content-Security-Policy", "default-src 'none'; script-src 'self'; style-src 'self'; img-src 'self'; connect-src 'self'; base-uri 'none'; frame-ancestors 'none'; form-action 'none'")
			_, _ = w.Write(data)
			return
		}
	}
	// Browser writes must come from the page itself.
	if r.Method == http.MethodPost && strings.HasPrefix(r.URL.Path, "/proxy/api/") && r.Header.Get("Origin") != p.cfg.Origin {
		jsonReply(w, 403, map[string]string{"error": "Use the proxy page"})
		return
	}
	switch r.Method + " " + r.URL.Path {
	case "GET /authorize":
		if !p.authorize(w, r) {
			return
		}
		w.Header().Set("Authorization", "Bearer "+bearer(r, true))
		jsonReply(w, 200, map[string]any{})
	case "GET /proxy/api/status":
		if !p.limit("status:"+clientIP(r), 120, time.Minute) {
			jsonReply(w, 429, map[string]string{"error": "Too many requests"})
			return
		}
		// Anyone may look. Sending the management key adds the owner's controls.
		role := "public"
		if bearer(r, false) != "" {
			if !p.authorizeOwner(w, r) {
				return
			}
			role = "owner"
		}
		data, err := p.statusData(role)
		if err != nil {
			jsonReply(w, 503, map[string]string{"error": "Status unavailable"})
			return
		}
		jsonReply(w, 200, data)
	case "POST /proxy/api/refresh":
		p.refreshAccount(w, r)
	case "POST /proxy/api/accounts":
		if p.authorizeOwner(w, r) {
			p.updateAccount(w, r)
		}
	case "POST /proxy/api/clients":
		if p.authorizeOwner(w, r) {
			p.adminClients(w, r)
		}
	case "POST /proxy/api/models":
		if p.authorizeOwner(w, r) {
			p.updateModel(w, r)
		}
	case "POST /proxy/api/options":
		if p.authorizeOwner(w, r) {
			p.updateOption(w, r)
		}
	case "POST /proxy/api/resets":
		if p.authorizeOwner(w, r) {
			p.redeemReset(w, r)
		}
	case "POST /proxy/api/join/start", "POST /proxy/api/join/complete", "POST /proxy/api/join/status":
		p.enroll(w, r)
	default:
		jsonReply(w, 404, map[string]string{"error": "Not found"})
	}
}

func clientIP(r *http.Request) string {
	if ip := strings.TrimSpace(strings.Split(r.Header.Get("X-Forwarded-For"), ",")[0]); ip != "" {
		return ip
	}
	return r.RemoteAddr
}

// updateAccount changes a display name locally or routing through the stock handler.
func (p *Pool) updateAccount(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID       string  `json:"id"`
		Disabled *bool   `json:"disabled"`
		Name     *string `json:"name"`
	}
	if !decodeBody(w, r, &body) {
		return
	}
	if (body.Disabled == nil) == (body.Name == nil) {
		jsonReply(w, 400, map[string]string{"error": "Choose a display name or routing state"})
		return
	}
	for _, account := range p.manager.List() {
		if account.EnsureIndex() != body.ID || coreauth.IsConfigAPIKeyAuth(account) || (account.Provider != "claude" && account.Provider != "codex") {
			continue
		}
		if body.Name != nil {
			name := strings.TrimSpace(*body.Name)
			if len([]rune(name)) > 80 {
				jsonReply(w, 400, map[string]string{"error": "Account names can be at most 80 characters"})
				return
			}
			var err error
			if name == "" {
				_, err = p.db.ExecContext(r.Context(), `DELETE FROM account_names WHERE account=?`, account.ID)
			} else {
				_, err = p.db.ExecContext(r.Context(), `INSERT INTO account_names(account,name) VALUES(?,?) ON CONFLICT(account) DO UPDATE SET name=excluded.name`, account.ID, name)
			}
			if err != nil {
				jsonReply(w, 503, map[string]string{"error": "Could not save the account name"})
				return
			}
		} else if _, err := p.invokeOAuth(r.Context(), p.oauth.PatchAuthFileStatus, nil, map[string]any{"name": account.ID, "auth_index": body.ID, "disabled": *body.Disabled}); err != nil {
			jsonReply(w, 502, map[string]string{"error": "Could not update the account"})
			return
		} else if !*body.Disabled {
			// Maintenance skips disabled accounts, so a re-enabled one has no fresh quota to rank it by.
			if _, err = p.maintenance.Usage(r.Context(), account); err != nil {
				log.WithError(err).Warn("pool: could not read usage for a re-enabled account")
			}
		}
		jsonReply(w, 200, map[string]string{})
		return
	}
	jsonReply(w, 404, map[string]string{"error": "Unknown account"})
}

// refreshAccount rereads one account's limits, at most once a minute per account.
func (p *Pool) refreshAccount(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID string `json:"id"`
	}
	if !decodeBody(w, r, &body) {
		return
	}
	for _, account := range p.manager.List() {
		if account.EnsureIndex() != body.ID || coreauth.IsConfigAPIKeyAuth(account) || (account.Provider != "claude" && account.Provider != "codex") {
			continue
		}
		if !p.limit("refresh:"+account.ID, 1, time.Minute) {
			jsonReply(w, 429, map[string]string{"error": "Refreshed less than a minute ago"})
			return
		}
		if _, err := p.maintenance.Usage(r.Context(), account); err != nil {
			jsonReply(w, 502, map[string]string{"error": "The provider could not return usage"})
			return
		}
		jsonReply(w, 200, map[string]string{})
		return
	}
	jsonReply(w, 404, map[string]string{"error": "Unknown account"})
}

// invokeOAuth calls the existing handler with a constructed request, never caller headers or URLs.
func (p *Pool) invokeOAuth(ctx context.Context, handler gin.HandlerFunc, query url.Values, body any) (map[string]any, error) {
	var data []byte
	if body != nil {
		var err error
		data, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request, _ = http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost/?"+query.Encode(), bytes.NewReader(data))
	if body != nil {
		c.Request.Method = http.MethodPost
		c.Request.Header.Set("Content-Type", "application/json")
	}
	handler(c)
	var result map[string]any
	if recorder.Code < 200 || recorder.Code >= 300 || json.Unmarshal(recorder.Body.Bytes(), &result) != nil {
		return nil, fmt.Errorf("OAuth operation failed")
	}
	return result, nil
}

func (p *Pool) enroll(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/proxy/api/join/start" {
		ip := clientIP(r)
		if !p.limit("enroll:"+ip, 6, time.Hour) || !p.limit("enroll:global", 60, time.Hour) {
			jsonReply(w, 429, map[string]string{"error": "Too many sign-ins. Try again later."})
			return
		}
		var body struct {
			Provider string `json:"provider"`
		}
		if !decodeBody(w, r, &body) {
			return
		}
		if body.Provider != "claude" && body.Provider != "codex" {
			jsonReply(w, 400, map[string]string{"error": "Choose Claude or Codex"})
			return
		}
		p.sessionMu.Lock()
		for token, session := range p.sessions {
			if session.Expires.Before(time.Now()) {
				delete(p.sessions, token)
			}
		}
		if len(p.sessions) >= 20 {
			p.sessionMu.Unlock()
			jsonReply(w, 429, map[string]string{"error": "Enrollment is busy. Try again in five minutes."})
			return
		}
		// Reserve capacity before starting the native asynchronous OAuth exchange.
		random := make([]byte, 32)
		if _, err := rand.Read(random); err != nil {
			p.sessionMu.Unlock()
			jsonReply(w, 500, map[string]string{"error": "Could not start sign-in"})
			return
		}
		token := base64.RawURLEncoding.EncodeToString(random)
		provider := body.Provider
		if provider == "claude" {
			provider = "anthropic"
		}
		p.sessions[token] = enrollment{Provider: provider, Expires: time.Now().Add(5 * time.Minute)}
		p.sessionMu.Unlock()
		handler := p.oauth.RequestCodexToken
		if provider == "anthropic" {
			handler = p.oauth.RequestAnthropicToken
		}
		result, err := p.invokeOAuth(r.Context(), handler, nil, nil)
		if err != nil {
			p.sessionMu.Lock()
			delete(p.sessions, token)
			p.sessionMu.Unlock()
			jsonReply(w, 502, map[string]string{"error": "Could not start sign-in"})
			return
		}
		p.sessionMu.Lock()
		session := p.sessions[token]
		session.State, _ = result["state"].(string)
		p.sessions[token] = session
		p.sessionMu.Unlock()
		jsonReply(w, 200, map[string]string{"session": token, "url": fmt.Sprint(result["url"])})
		return
	}
	token := bearer(r, false)
	p.sessionMu.Lock()
	session, ok := p.sessions[token]
	p.sessionMu.Unlock()
	if !ok || session.State == "" || session.Expires.Before(time.Now()) {
		jsonReply(w, 400, map[string]string{"error": "Sign-in expired. Start again."})
		return
	}
	if !p.limit("session:"+token, 60, time.Minute) {
		jsonReply(w, 429, map[string]string{"error": "Too many requests"})
		return
	}
	if r.URL.Path == "/proxy/api/join/status" {
		var body map[string]any
		if !decodeBody(w, r, &body) {
			return
		}
		result, err := p.invokeOAuth(r.Context(), p.oauth.GetAuthStatus, url.Values{"state": {session.State}}, nil)
		if err != nil {
			jsonReply(w, 502, map[string]string{"error": "Could not check sign-in"})
			return
		}
		status, _ := result["status"].(string)
		if status == "ok" {
			p.sessionMu.Lock()
			delete(p.sessions, token)
			p.sessionMu.Unlock()
		}
		reply := map[string]string{"status": status}
		if status == "error" {
			reply["error"] = "Sign-in failed. Start again."
		}
		jsonReply(w, 200, reply)
		return
	}
	var body struct {
		Callback string `json:"callback"`
	}
	if !decodeBody(w, r, &body) {
		return
	}
	raw := strings.TrimSpace(body.Callback)
	code, state, _ := strings.Cut(raw, "#")
	if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") {
		parsed, err := url.Parse(raw)
		if err != nil {
			jsonReply(w, 400, map[string]string{"error": "Invalid callback"})
			return
		}
		code, state = parsed.Query().Get("code"), parsed.Query().Get("state")
	}
	if code == "" || (state != "" && state != session.State) {
		jsonReply(w, 400, map[string]string{"error": "The callback does not match this sign-in"})
		return
	}
	_, err := p.invokeOAuth(r.Context(), p.oauth.PostOAuthCallback, nil, map[string]string{"provider": session.Provider, "state": session.State, "code": code})
	if err != nil {
		jsonReply(w, 400, map[string]string{"error": "Could not complete sign-in"})
		return
	}
	jsonReply(w, 200, map[string]string{"status": "wait"})
}
