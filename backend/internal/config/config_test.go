package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	cfg := Load()

	if cfg.Port != "8080" || cfg.StaticDir != "" {
		t.Errorf("got %+v, want port 8080 and no static dir", cfg)
	}
}

func TestLoadFromEnv(t *testing.T) {
	t.Setenv("PORT", "9090")
	t.Setenv("STATIC_DIR", "/srv/web")

	cfg := Load()

	if cfg.Port != "9090" || cfg.StaticDir != "/srv/web" {
		t.Errorf("got %+v, want env values", cfg)
	}
}
