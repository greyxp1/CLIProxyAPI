package pool

import (
	"context"
	"net/http"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/router-for-me/CLIProxyAPI/v8/internal/registry"
	"github.com/router-for-me/CLIProxyAPI/v8/internal/runtime/executor/helps"
	log "github.com/sirupsen/logrus"
)

type adminModel struct {
	ID       string `json:"id"`
	Name     string `json:"name,omitempty"`
	Provider string `json:"provider"`
	Enabled  bool   `json:"enabled"`
	// Older is set when a newer model of the same family exists.
	Older  bool `json:"older,omitempty"`
	levels []string
}

// generation splits a model ID into its family and version, so claude-opus-4-1-20250805 is
// claude-opus 4.1.20250805 and gpt-6.1-sol is gpt-sol 6.1. IDs without a version have none.
func generation(id string) (string, []int) {
	family := []string{}
	version := []int{}
	for _, part := range strings.Split(strings.ToLower(id), "-") {
		numbers := strings.Split(part, ".")
		parsed := make([]int, 0, len(numbers))
		for _, number := range numbers {
			value, err := strconv.Atoi(number)
			if err != nil {
				parsed = nil
				break
			}
			parsed = append(parsed, value)
		}
		if parsed == nil {
			family = append(family, part)
		} else {
			version = append(version, parsed...)
		}
	}
	return strings.Join(family, "-"), version
}

func newer(a, b []int) bool {
	for i := 0; i < len(a) && i < len(b); i++ {
		if a[i] != b[i] {
			return a[i] > b[i]
		}
	}
	return len(a) > len(b)
}

// olderModels returns the IDs that a newer model of the same family supersedes.
func olderModels(ids []string) map[string]bool {
	latest := map[string][]int{}
	for _, id := range ids {
		family, version := generation(id)
		if len(version) > 0 && (latest[family] == nil || newer(version, latest[family])) {
			latest[family] = version
		}
	}
	older := map[string]bool{}
	for _, id := range ids {
		family, version := generation(id)
		if len(version) > 0 && newer(latest[family], version) {
			older[id] = true
		}
	}
	return older
}

// modelSet reads one column of model IDs.
func (p *Pool) modelSet(ctx context.Context, query string) (map[string]bool, error) {
	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	set := map[string]bool{}
	for rows.Next() {
		var id string
		if err = rows.Scan(&id); err != nil {
			return nil, err
		}
		set[id] = true
	}
	return set, rows.Err()
}

// models lists every Claude and Codex model the pool's accounts offer. Older generations stay
// off unless switched on, so a new release replaces the previous one without the owner stepping in.
func (p *Pool) models(ctx context.Context) ([]adminModel, error) {
	blocked, err := p.modelSet(ctx, `SELECT model FROM blocked_models`)
	if err != nil {
		return nil, err
	}
	allowed, err := p.modelSet(ctx, `SELECT model FROM allowed_models`)
	if err != nil {
		return nil, err
	}
	options, err := p.poolOptions(ctx)
	if err != nil {
		return nil, err
	}
	list := []adminModel{}
	ids := []string{}
	for _, provider := range []string{"claude", "codex"} {
		for _, info := range registry.GetGlobalRegistry().GetAvailableModelsByProvider(provider) {
			model := adminModel{ID: info.ID, Name: info.DisplayName, Provider: provider}
			if info.Thinking != nil {
				model.levels = info.Thinking.Levels
			}
			list = append(list, model)
			ids = append(ids, info.ID)
		}
	}
	older := olderModels(ids)
	for i := range list {
		list[i].Older = older[list[i].ID]
		list[i].Enabled = !blocked[list[i].ID] && (!list[i].Older || allowed[list[i].ID]) && allowsModelEffort(list[i], options)
	}
	sort.Slice(list, func(i, j int) bool {
		if list[i].Provider != list[j].Provider {
			return list[i].Provider < list[j].Provider
		}
		return list[i].ID < list[j].ID
	})
	return list, nil
}

func allowsModelEffort(model adminModel, options helps.PoolOptions) bool {
	return len(model.levels) == 0 || slices.ContainsFunc(model.levels, func(level string) bool { return options.Efforts[model.Provider][level] })
}

// applyBlockedModels switches off the saved models and older generations.
func (p *Pool) applyBlockedModels(ctx context.Context) error {
	list, err := p.models(ctx)
	if err != nil {
		return err
	}
	blocked, err := p.modelSet(ctx, `SELECT model FROM blocked_models`)
	if err != nil {
		return err
	}
	ids := make([]string, 0, len(blocked))
	for id := range blocked {
		ids = append(ids, id)
	}
	for _, model := range list {
		if !model.Enabled && !blocked[model.ID] {
			ids = append(ids, model.ID)
		}
	}
	registry.GetGlobalRegistry().SetBlockedModels(ids)
	return nil
}

// watchModels reapplies the switches when accounts register models, so a new generation
// replaces the old one without the owner stepping in.
func (p *Pool) watchModels(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	var seen uint64
	for {
		if epoch := registry.GetGlobalRegistry().RegistrationEpoch(); epoch != seen {
			if err := p.applyBlockedModels(ctx); err != nil {
				log.WithError(err).Warn("pool: could not apply model switches")
			} else {
				seen = epoch
			}
		}
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (p *Pool) updateModel(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID      string `json:"id"`
		Enabled bool   `json:"enabled"`
	}
	if !decodeBody(w, r, &body) {
		return
	}
	list, err := p.models(r.Context())
	if err != nil {
		jsonReply(w, 503, map[string]string{"error": "Models unavailable"})
		return
	}
	index := slices.IndexFunc(list, func(model adminModel) bool { return model.ID == body.ID })
	if index < 0 {
		jsonReply(w, 404, map[string]string{"error": "Unknown model"})
		return
	}
	if body.Enabled {
		options, err := p.poolOptions(r.Context())
		if err != nil {
			jsonReply(w, 503, map[string]string{"error": "Options unavailable"})
			return
		}
		if !allowsModelEffort(list[index], options) {
			jsonReply(w, 400, map[string]string{"error": "Allow an effort level supported by " + body.ID})
			return
		}
	}
	// Older models are off by default, so their switch saves an exception rather than a block.
	table := `blocked_models`
	if list[index].Older {
		table = `allowed_models`
	}
	if body.Enabled == list[index].Older {
		_, err = p.db.ExecContext(r.Context(), `INSERT OR IGNORE INTO `+table+`(model) VALUES(?)`, body.ID)
	} else {
		_, err = p.db.ExecContext(r.Context(), `DELETE FROM `+table+` WHERE model=?`, body.ID)
	}
	if err == nil && body.Enabled {
		_, err = p.db.ExecContext(r.Context(), `DELETE FROM blocked_models WHERE model=?`, body.ID)
	}
	if err == nil {
		err = p.applyBlockedModels(r.Context())
	}
	if err != nil {
		jsonReply(w, 503, map[string]string{"error": "Could not update models"})
		return
	}
	jsonReply(w, 200, map[string]string{})
}

type adminEffort struct {
	Level   string `json:"level"`
	Enabled bool   `json:"enabled"`
}

type adminOptions struct {
	Efforts []adminEffort `json:"efforts"`
	Fast    bool          `json:"fast"`
}

// poolOptions returns the saved effort and fast mode switches over the defaults.
func (p *Pool) poolOptions(ctx context.Context) (helps.PoolOptions, error) {
	options := helps.DefaultPoolOptions()
	rows, err := p.db.QueryContext(ctx, `SELECT provider, setting, enabled FROM options`)
	if err != nil {
		return options, err
	}
	defer rows.Close()
	for rows.Next() {
		var provider, setting string
		var enabled bool
		if err = rows.Scan(&provider, &setting, &enabled); err != nil {
			return options, err
		}
		if setting == "fast" {
			options.Fast[provider] = enabled
		} else if efforts := options.Efforts[provider]; efforts != nil {
			efforts[setting] = enabled
		}
	}
	return options, rows.Err()
}

// adminOptionList lists each provider's switches, with the effort levels its current models offer.
func adminOptionList(options helps.PoolOptions, models []adminModel) map[string]adminOptions {
	list := map[string]adminOptions{}
	for _, provider := range []string{"claude", "codex"} {
		offered := map[string]bool{}
		for _, model := range models {
			if model.Provider == provider && !model.Older {
				for _, level := range model.levels {
					offered[level] = true
				}
			}
		}
		item := adminOptions{Efforts: []adminEffort{}, Fast: options.Fast[provider]}
		for _, level := range helps.EffortLevels {
			if offered[level] {
				item.Efforts = append(item.Efforts, adminEffort{Level: level, Enabled: options.Efforts[provider][level]})
			}
		}
		list[provider] = item
	}
	return list
}

func (p *Pool) applyOptions(ctx context.Context) error {
	options, err := p.poolOptions(ctx)
	if err == nil {
		helps.SetPoolOptions(options)
		err = p.applyBlockedModels(ctx)
	}
	return err
}

func (p *Pool) updateOption(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Provider string `json:"provider"`
		Setting  string `json:"setting"`
		Enabled  bool   `json:"enabled"`
	}
	if !decodeBody(w, r, &body) {
		return
	}
	options, err := p.poolOptions(r.Context())
	if err != nil {
		jsonReply(w, 503, map[string]string{"error": "Options unavailable"})
		return
	}
	efforts := options.Efforts[body.Provider]
	if efforts == nil || body.Setting != "fast" && !slices.Contains(helps.EffortLevels, body.Setting) {
		jsonReply(w, 404, map[string]string{"error": "Unknown option"})
		return
	}
	if body.Setting != "fast" && !body.Enabled {
		efforts[body.Setting] = false
		if !slices.ContainsFunc(helps.EffortLevels, func(level string) bool { return efforts[level] }) {
			jsonReply(w, 400, map[string]string{"error": "Keep at least one effort level"})
			return
		}
		models, err := p.models(r.Context())
		if err != nil {
			jsonReply(w, 503, map[string]string{"error": "Models unavailable"})
			return
		}
		for _, model := range models {
			if model.Provider == body.Provider && model.Enabled && !allowsModelEffort(model, options) {
				jsonReply(w, 400, map[string]string{"error": "Keep an effort level supported by " + model.ID})
				return
			}
		}
	}
	_, err = p.db.ExecContext(r.Context(), `INSERT INTO options VALUES (?,?,?) ON CONFLICT(provider,setting) DO UPDATE SET enabled=excluded.enabled`,
		body.Provider, body.Setting, body.Enabled)
	if err == nil {
		err = p.applyOptions(r.Context())
	}
	if err != nil {
		jsonReply(w, 503, map[string]string{"error": "Could not update options"})
		return
	}
	jsonReply(w, 200, map[string]string{})
}
