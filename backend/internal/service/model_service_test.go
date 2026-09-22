package service

import (
	"database/sql"
	"path/filepath"
	"testing"

	"omnirelay/internal/database"
	"omnirelay/internal/models"
)

func newTestModelService(t *testing.T) *ModelService {
	t.Helper()
	db, err := database.Init(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("init db: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return NewModelService(db)
}

func seedProvider(t *testing.T, db *sql.DB, key string, userID *int64) int64 {
	t.Helper()
	var id int64
	if userID == nil {
		err := db.QueryRow(
			`INSERT INTO providers (provider_key, name, provider_type, api_base_url, api_key_encrypted, is_active, user_id)
			 VALUES (?, ?, 'openai', 'https://api.example.com', 'enc', 1, NULL) RETURNING id`,
			key, key,
		).Scan(&id)
		if err != nil {
			t.Fatalf("insert provider: %v", err)
		}
		return id
	}
	err := db.QueryRow(
		`INSERT INTO providers (provider_key, name, provider_type, api_base_url, api_key_encrypted, is_active, user_id)
		 VALUES (?, ?, 'openai', 'https://api.example.com', 'enc', 1, ?) RETURNING id`,
		key, key, *userID,
	).Scan(&id)
	if err != nil {
		t.Fatalf("insert provider: %v", err)
	}
	return id
}

func TestCountActiveScopesByUser(t *testing.T) {
	svc := newTestModelService(t)
	db := svc.db

	pShared := seedProvider(t, db, "shared", nil)
	pUser1 := seedProvider(t, db, "u1", int64Ptr(1))

	_, err := db.Exec(
		`INSERT INTO models (provider_id, model_id, provider_key, is_manual, user_id) VALUES (?, 'm-shared', 'shared', 1, NULL)`,
		pShared,
	)
	if err != nil {
		t.Fatalf("insert shared model: %v", err)
	}
	_, err = db.Exec(
		`INSERT INTO models (provider_id, model_id, provider_key, is_manual, user_id) VALUES (?, 'm-u1', 'u1', 1, 1)`,
		pUser1,
	)
	if err != nil {
		t.Fatalf("insert user model: %v", err)
	}

	count, err := svc.CountActive(1)
	if err != nil {
		t.Fatalf("CountActive: %v", err)
	}
	if count != 2 {
		t.Errorf("user 1 count = %d, want 2 (shared + own)", count)
	}

	count, err = svc.CountActive(2)
	if err != nil {
		t.Fatalf("CountActive user 2: %v", err)
	}
	if count != 1 {
		t.Errorf("user 2 count = %d, want 1 (shared only)", count)
	}
}

func int64Ptr(v int64) *int64 {
	return &v
}

func seedProviderOf(t *testing.T, db *sql.DB, key, ptype string, userID *int64) int64 {
	t.Helper()
	var id int64
	q := `INSERT INTO providers (provider_key, name, provider_type, api_base_url, api_key_encrypted, is_active, user_id)
		 VALUES (?, ?, ?, 'https://api.example.com', 'enc', 1, `
	var err error
	if userID == nil {
		err = db.QueryRow(q+`NULL) RETURNING id`, key, key, ptype).Scan(&id)
	} else {
		err = db.QueryRow(q+`?) RETURNING id`, key, key, ptype, *userID).Scan(&id)
	}
	if err != nil {
		t.Fatalf("insert provider: %v", err)
	}
	return id
}

func TestCreateFillsZeroPricesFromCatalog(t *testing.T) {
	svc := newTestModelService(t)
	pid := seedProviderOf(t, svc.db, "anthropic", "anthropic", nil)
	c, _ := catalogServing(t, testCatalogJSON)
	svc.pricing = c

	m, err := svc.Create(models.CreateModelRequest{ProviderID: pid, ModelID: "claude-test"}, 0)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if m.InputPricePer1MTok != 3 {
		t.Errorf("input price = %v, want 3", m.InputPricePer1MTok)
	}
	if m.OutputPricePer1MTok != 15 {
		t.Errorf("output price = %v, want 15", m.OutputPricePer1MTok)
	}
	if m.CacheWrite5mPricePer1MTok != 3.75 {
		t.Errorf("cache write 5m price = %v, want 3.75", m.CacheWrite5mPricePer1MTok)
	}
	if m.CacheReadPricePer1MTok != 0.3 {
		t.Errorf("cache read price = %v, want 0.3", m.CacheReadPricePer1MTok)
	}
	if m.CacheWrite1hPricePer1MTok != 0 {
		t.Errorf("cache write 1h price = %v, want 0 (not in models.dev)", m.CacheWrite1hPricePer1MTok)
	}
}

func TestCreateKeepsUserEnteredPrices(t *testing.T) {
	svc := newTestModelService(t)
	pid := seedProviderOf(t, svc.db, "anthropic", "anthropic", nil)
	c, _ := catalogServing(t, testCatalogJSON)
	svc.pricing = c

	m, err := svc.Create(models.CreateModelRequest{
		ProviderID: pid, ModelID: "claude-test",
		InputPricePer1MTok: 9,
	}, 0)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if m.InputPricePer1MTok != 9 {
		t.Errorf("input price = %v, want 9 (user value kept)", m.InputPricePer1MTok)
	}
	if m.OutputPricePer1MTok != 15 {
		t.Errorf("output price = %v, want 15 (zero field filled)", m.OutputPricePer1MTok)
	}
}

func TestCreateSkipsCatalogWhenAllPricesEntered(t *testing.T) {
	svc := newTestModelService(t)
	pid := seedProviderOf(t, svc.db, "anthropic", "anthropic", nil)
	c, hits := catalogServing(t, testCatalogJSON)
	svc.pricing = c

	m, err := svc.Create(models.CreateModelRequest{
		ProviderID: pid, ModelID: "claude-test",
		InputPricePer1MTok: 1, OutputPricePer1MTok: 1,
		CacheWrite5mPricePer1MTok: 1, CacheReadPricePer1MTok: 1,
	}, 0)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if m.InputPricePer1MTok != 1 {
		t.Errorf("input price = %v, want 1", m.InputPricePer1MTok)
	}
	if n := hits.Load(); n != 0 {
		t.Errorf("catalog fetched %d times, want 0 (no zero price fields)", n)
	}
}

func TestCreateWithoutCatalogKeepsZeroPrices(t *testing.T) {
	svc := newTestModelService(t)
	pid := seedProviderOf(t, svc.db, "openai", "openai", nil)

	m, err := svc.Create(models.CreateModelRequest{ProviderID: pid, ModelID: "gpt-test"}, 0)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if m.InputPricePer1MTok != 0 || m.OutputPricePer1MTok != 0 {
		t.Errorf("prices = %v/%v, want 0/0 without a pricing catalog", m.InputPricePer1MTok, m.OutputPricePer1MTok)
	}
}

func TestCreateSurvivesCatalogFailure(t *testing.T) {
	svc := newTestModelService(t)
	pid := seedProviderOf(t, svc.db, "anthropic", "anthropic", nil)
	c, _ := catalogFailing(t)
	svc.pricing = c

	m, err := svc.Create(models.CreateModelRequest{ProviderID: pid, ModelID: "claude-test"}, 0)
	if err != nil {
		t.Fatalf("Create with failing catalog: %v", err)
	}
	if m.InputPricePer1MTok != 0 {
		t.Errorf("input price = %v, want 0 on catalog failure", m.InputPricePer1MTok)
	}
}

func TestSyncFillsNewModelAndKeepsExistingPrices(t *testing.T) {
	svc := newTestModelService(t)
	pid := seedProviderOf(t, svc.db, "anthropic", "anthropic", int64Ptr(1))
	if _, err := svc.db.Exec(
		`INSERT INTO models (provider_id, model_id, display_name, provider_key, is_manual, input_price_per_1mtok, output_price_per_1mtok, user_id)
		 VALUES (?, 'claude-test', 'claude-test', 'anthropic', 0, 9, 9, 1)`,
		pid,
	); err != nil {
		t.Fatalf("seed existing model: %v", err)
	}
	c, _ := catalogServing(t, testCatalogJSON)
	svc.pricing = c

	if err := svc.SyncFromProvider(pid, "anthropic", []string{"claude-test", "claude-new"}, 1); err != nil {
		t.Fatalf("SyncFromProvider: %v", err)
	}

	existing, err := svc.List("anthropic", 1)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	byID := map[string]models.Model{}
	for _, m := range existing {
		byID[m.ModelID] = m
	}

	old, ok := byID["claude-test"]
	if !ok {
		t.Fatal("claude-test missing after sync")
	}
	if old.InputPricePer1MTok != 9 || old.OutputPricePer1MTok != 9 {
		t.Errorf("existing prices = %v/%v, want 9/9 (preserved)", old.InputPricePer1MTok, old.OutputPricePer1MTok)
	}

	fresh, ok := byID["claude-new"]
	if !ok {
		t.Fatal("claude-new missing after sync")
	}
	if fresh.InputPricePer1MTok != 1 {
		t.Errorf("new model input price = %v, want 1 (filled from catalog)", fresh.InputPricePer1MTok)
	}
	if fresh.OutputPricePer1MTok != 5 {
		t.Errorf("new model output price = %v, want 5 (filled from catalog)", fresh.OutputPricePer1MTok)
	}
}
