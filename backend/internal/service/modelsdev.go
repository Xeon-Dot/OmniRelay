package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"sync"
	"time"

	"omnirelay/internal/models"
)

// ModelsDevCatalog caches model pricing from https://models.dev/api.json.
// Prices are USD per 1M tokens, the same unit as the models table columns.
type ModelsDevCatalog struct {
	url    string
	ttl    time.Duration
	client *http.Client

	mu        sync.Mutex
	providers map[string]mdProvider
	order     []string
	nextFetch time.Time
}

type mdProvider struct {
	Models map[string]mdModel `json:"models"`
}

type mdModel struct {
	Cost *mdCost `json:"cost"`
}

type mdCost struct {
	Input      *float64 `json:"input"`
	Output     *float64 `json:"output"`
	CacheRead  *float64 `json:"cache_read"`
	CacheWrite *float64 `json:"cache_write"`
}

// ModelPricing holds models.dev prices in USD per 1M tokens. A nil field
// means models.dev does not publish that price (e.g. the 1h cache write).
type ModelPricing struct {
	Input        *float64
	Output       *float64
	CacheWrite5m *float64
	CacheRead    *float64
}

const (
	modelsDevAPIURL    = "https://models.dev/api.json"
	modelsDevTTL       = time.Hour
	modelsDevRetryWait = time.Minute
	modelsDevTimeout   = 10 * time.Second
)

// providerTypeDefaults maps OmniRelay provider types to models.dev provider IDs.
var providerTypeDefaults = map[string]string{
	"openai":    "openai",
	"anthropic": "anthropic",
	"gemini":    "google",
}

// NewModelsDevCatalog returns a catalog backed by the live models.dev API.
func NewModelsDevCatalog() *ModelsDevCatalog {
	return newModelsDevCatalog(modelsDevAPIURL, modelsDevTTL)
}

func newModelsDevCatalog(url string, ttl time.Duration) *ModelsDevCatalog {
	return &ModelsDevCatalog{
		url:    url,
		ttl:    ttl,
		client: &http.Client{Timeout: modelsDevTimeout},
	}
}

// Lookup finds pricing for modelID. Stage 1 tries the provider key, then the
// provider_type default; stage 2 scans every provider for the model ID alone
// so routing services and third-party endpoints still get prices. Returns nil
// when nothing is published or fetching fails, so callers keep zero prices.
func (c *ModelsDevCatalog) Lookup(providerKey, providerType, modelID string) *ModelPricing {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.ensureLoaded()
	if c.providers == nil {
		return nil
	}

	for _, key := range stageOneKeys(providerKey, providerType) {
		if p := c.providers[key].lookupModel(modelID); p != nil {
			return p
		}
	}
	for _, key := range c.order {
		if p := c.providers[key].lookupModel(modelID); p != nil {
			return p
		}
	}
	return nil
}

func stageOneKeys(providerKey, providerType string) []string {
	keys := []string{providerKey}
	if def, ok := providerTypeDefaults[providerType]; ok && def != providerKey {
		keys = append(keys, def)
	}
	return keys
}

func (p mdProvider) lookupModel(modelID string) *ModelPricing {
	m, ok := p.Models[modelID]
	if !ok {
		return nil
	}
	return m.pricing()
}

func (m mdModel) pricing() *ModelPricing {
	if m.Cost == nil {
		return nil
	}
	p := &ModelPricing{
		Input:        m.Cost.Input,
		Output:       m.Cost.Output,
		CacheWrite5m: m.Cost.CacheWrite,
		CacheRead:    m.Cost.CacheRead,
	}
	if p.Input == nil && p.Output == nil && p.CacheWrite5m == nil && p.CacheRead == nil {
		return nil
	}
	return p
}

// ensureLoaded refreshes the catalog when stale. Callers hold c.mu. A failed
// fetch is negative-cached for a short window so outages do not stall adds;
// the last good data (if any) stays usable.
func (c *ModelsDevCatalog) ensureLoaded() {
	now := time.Now()
	if !c.nextFetch.IsZero() && now.Before(c.nextFetch) {
		return
	}
	if err := c.fetch(); err != nil {
		log.Printf("models.dev pricing fetch failed: %v", err)
		c.nextFetch = now.Add(modelsDevRetryWait)
		return
	}
	c.nextFetch = now.Add(c.ttl)
}

func (c *ModelsDevCatalog) fetch() error {
	resp, err := c.client.Get(c.url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %s", resp.Status)
	}

	providers := map[string]mdProvider{}
	if err := json.NewDecoder(resp.Body).Decode(&providers); err != nil {
		return err
	}

	order := make([]string, 0, len(providers))
	for key := range providers {
		order = append(order, key)
	}
	sort.Strings(order)

	c.providers = providers
	c.order = order
	return nil
}

// fillFromCatalog fills only the still-zero price fields of a model being
// created. It is a no-op without a catalog or when every field is set.
func (s *ModelService) fillFromCatalog(req *models.CreateModelRequest, providerKey, providerType string) {
	if s.pricing == nil {
		return
	}
	if req.InputPricePer1MTok != 0 && req.OutputPricePer1MTok != 0 &&
		req.CacheWrite5mPricePer1MTok != 0 && req.CacheReadPricePer1MTok != 0 {
		return
	}

	p := s.pricing.Lookup(providerKey, providerType, req.ModelID)
	if p == nil {
		return
	}
	req.InputPricePer1MTok = fillZeroPrice(req.InputPricePer1MTok, p.Input)
	req.OutputPricePer1MTok = fillZeroPrice(req.OutputPricePer1MTok, p.Output)
	req.CacheWrite5mPricePer1MTok = fillZeroPrice(req.CacheWrite5mPricePer1MTok, p.CacheWrite5m)
	req.CacheReadPricePer1MTok = fillZeroPrice(req.CacheReadPricePer1MTok, p.CacheRead)
}
