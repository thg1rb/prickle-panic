package appconfig

import (
	"net/netip"
	"reflect"
	"testing"
	"time"
)

func TestLoadRequiresExplicitModeAndOrigin(t *testing.T) {
	for _, values := range []map[string]string{
		{"APP_ORIGIN": "http://localhost:5173"},
		{"APP_ENV": "development"},
		{"APP_ENV": "staging", "APP_ORIGIN": "https://example.test"},
	} {
		if _, err := Load(func(key string) string { return values[key] }); err == nil {
			t.Fatalf("accepted incomplete/invalid config: %v", values)
		}
	}
}

func TestLoadDefaultsAndApprovedOverrides(t *testing.T) {
	values := map[string]string{"APP_ENV": "development", "APP_ORIGIN": "http://localhost:5173", "SESSION_IDLE_TTL": "12h", "SESSION_CREATION_LIMIT": "8", "TRUSTED_PROXY_CIDRS": "10.0.0.0/24, 192.0.2.1/32"}
	cfg, err := Load(func(key string) string { return values[key] })
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Port != "8080" || cfg.IdleTTL != 12*time.Hour || cfg.AbsoluteTTL != 7*24*time.Hour || cfg.CreationLimit != 8 || cfg.CreationWindow != 10*time.Minute || cfg.MaxSourceKeys != 65536 || cfg.TrustedProxyMaxHops != 32 {
		t.Fatalf("unexpected effective config: %+v", cfg)
	}
	if !reflect.DeepEqual(cfg.TrustedProxyCIDRs, []netip.Prefix{netip.MustParsePrefix("10.0.0.0/24"), netip.MustParsePrefix("192.0.2.1/32")}) {
		t.Fatalf("unexpected trusted proxies: %v", cfg.TrustedProxyCIDRs)
	}
}

func TestLoadRejectsUnsafeValues(t *testing.T) {
	base := map[string]string{"APP_ENV": "production", "APP_ORIGIN": "https://game.example.test"}
	for key, value := range map[string]string{
		"PORT":                   "70000",
		"SESSION_IDLE_TTL":       "0s",
		"SESSION_CREATION_LIMIT": "0",
		"TRUSTED_PROXY_CIDRS":    "10.0.0.1/24",
		"TRUSTED_PROXY_MAX_HOPS": "257",
	} {
		values := map[string]string{}
		for k, v := range base {
			values[k] = v
		}
		values[key] = value
		if _, err := Load(func(k string) string { return values[k] }); err == nil {
			t.Errorf("accepted %s=%q", key, value)
		}
	}
	values := map[string]string{"APP_ENV": "production", "APP_ORIGIN": "https://game.example.test", "SESSION_EXPIRY_RETRY_MIN": "2s", "SESSION_EXPIRY_RETRY_MAX": "1s"}
	if _, err := Load(func(k string) string { return values[k] }); err == nil {
		t.Fatal("accepted retry maximum below minimum")
	}
}
