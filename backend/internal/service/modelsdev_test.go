package service

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

const testCatalogJSON = `{
  "anthropic": {"models": {
    "claude-test": {"cost": {"input": 3, "output": 15, "cache_read": 0.3, "cache_write": 3.75}},
    "claude-new": {"cost": {"input": 1, "output": 5}}
  }},
  "openai": {"models": {
    "gpt-test": {"cost": {"input": 1.25, "output": 10}}
  }},
  "mystery-inc": {"models": {"mystery-model": {}}
  }
}`

func catalogWithHandler(t *testing.T, h http.HandlerFunc) (*ModelsDevCatalog, *atomic.Int32) {
	t.Helper()
	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		h(w, r)
	}))
	t.Cleanup(srv.Close)
	return newModelsDevCatalog(srv.URL, time.Hour), &hits
}

func catalogServing(t *testing.T, body string) (*ModelsDevCatalog, *atomic.Int32) {
	t.Helper()
	return catalogWithHandler(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(body))
	})
}

func catalogFailing(t *testing.T) (*ModelsDevCatalog, *atomic.Int32) {
	t.Helper()
	return catalogWithHandler(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
}

func TestModelsDevLookupByProviderKey(t *testing.T) {
	c, _ := catalogServing(t, testCatalogJSON)

	p := c.Lookup("anthropic", "anthropic", "claude-test")
	if p == nil {
		t.Fatal("Lookup returned nil, want pricing for claude-test")
	}
	if p.Input == nil || *p.Input != 3 {
		t.Errorf("Input = %v, want 3", p.Input)
	}
	if p.Output == nil || *p.Output != 15 {
		t.Errorf("Output = %v, want 15", p.Output)
	}
	if p.CacheRead == nil || *p.CacheRead != 0.3 {
		t.Errorf("CacheRead = %v, want 0.3", p.CacheRead)
	}
	if p.CacheWrite5m == nil || *p.CacheWrite5m != 3.75 {
		t.Errorf("CacheWrite5m = %v, want 3.75", p.CacheWrite5m)
	}
}

func TestModelsDevLookupFallsBackToProviderType(t *testing.T) {
	c, _ := catalogServing(t, testCatalogJSON)

	p := c.Lookup("my-anthropic", "anthropic", "claude-test")
	if p == nil {
		t.Fatal("Lookup returned nil, want provider_type fallback to match")
	}
	if p.Input == nil || *p.Input != 3 {
		t.Errorf("Input = %v, want 3", p.Input)
	}
}

func TestModelsDevLookupFallsBackToGlobalModelSearch(t *testing.T) {
	c, _ := catalogServing(t, testCatalogJSON)

	p := c.Lookup("my-router", "custom", "gpt-test")
	if p == nil {
		t.Fatal("Lookup returned nil, want global model search to match")
	}
	if p.Input == nil || *p.Input != 1.25 {
		t.Errorf("Input = %v, want 1.25", p.Input)
	}
	if p.Output == nil || *p.Output != 10 {
		t.Errorf("Output = %v, want 10", p.Output)
	}
	if p.CacheRead != nil {
		t.Errorf("CacheRead = %v, want nil (absent in catalog)", *p.CacheRead)
	}
	if p.CacheWrite5m != nil {
		t.Errorf("CacheWrite5m = %v, want nil (absent in catalog)", *p.CacheWrite5m)
	}
}

func TestModelsDevLookupMissReturnsNil(t *testing.T) {
	c, _ := catalogServing(t, testCatalogJSON)

	if p := c.Lookup("openai", "openai", "unknown-model"); p != nil {
		t.Errorf("unknown model = %+v, want nil", p)
	}
	if p := c.Lookup("mystery-inc", "custom", "mystery-model"); p != nil {
		t.Errorf("model without cost data = %+v, want nil", p)
	}
}

func TestModelsDevCatalogFetchesOnce(t *testing.T) {
	c, hits := catalogServing(t, testCatalogJSON)

	if p := c.Lookup("anthropic", "anthropic", "claude-test"); p == nil {
		t.Fatal("first Lookup returned nil")
	}
	if p := c.Lookup("openai", "openai", "gpt-test"); p == nil {
		t.Fatal("second Lookup returned nil")
	}
	if n := hits.Load(); n != 1 {
		t.Errorf("catalog fetched %d times, want 1 (cached)", n)
	}
}

func TestModelsDevLookupFailureReturnsNilAndCachesFailure(t *testing.T) {
	c, hits := catalogFailing(t)

	if p := c.Lookup("anthropic", "anthropic", "claude-test"); p != nil {
		t.Errorf("Lookup on failing catalog = %+v, want nil", p)
	}
	if p := c.Lookup("openai", "openai", "gpt-test"); p != nil {
		t.Errorf("second Lookup on failing catalog = %+v, want nil", p)
	}
	if n := hits.Load(); n != 1 {
		t.Errorf("catalog fetched %d times, want 1 (failure cached)", n)
	}
}
