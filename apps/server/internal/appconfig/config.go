package appconfig

import (
	"fmt"
	"net/netip"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Mode                string
	Origin              string
	Port                string
	IdleTTL             time.Duration
	AbsoluteTTL         time.Duration
	CreationLimit       int
	CreationWindow      time.Duration
	MaxSourceKeys       int
	OutboxCapacity      int
	RetryMin            time.Duration
	RetryMax            time.Duration
	TrustedProxyCIDRs   []netip.Prefix
	TrustedProxyMaxHops int
}

// Load reads explicit application configuration through lookup, which keeps
// validation deterministic and straightforward to test.
func Load(lookup func(string) string) (Config, error) {
	if lookup == nil {
		return Config{}, fmt.Errorf("environment lookup is required")
	}
	mode := strings.TrimSpace(lookup("APP_ENV"))
	if mode != "development" && mode != "production" {
		return Config{}, fmt.Errorf("APP_ENV must be development or production")
	}
	origin := strings.TrimSpace(lookup("APP_ORIGIN"))
	if origin == "" {
		return Config{}, fmt.Errorf("APP_ORIGIN is required")
	}
	port := lookup("PORT")
	if port == "" {
		port = "8080"
	}
	portNumber, err := strconv.Atoi(port)
	if err != nil || portNumber < 1 || portNumber > 65535 {
		return Config{}, fmt.Errorf("PORT must be between 1 and 65535")
	}
	cfg := Config{Mode: mode, Origin: origin, Port: strconv.Itoa(portNumber)}
	if cfg.IdleTTL, err = duration(lookup, "SESSION_IDLE_TTL", 24*time.Hour); err != nil {
		return Config{}, err
	}
	if cfg.AbsoluteTTL, err = duration(lookup, "SESSION_ABSOLUTE_TTL", 7*24*time.Hour); err != nil {
		return Config{}, err
	}
	if cfg.CreationLimit, err = integer(lookup, "SESSION_CREATION_LIMIT", 10, 1, 1000000); err != nil {
		return Config{}, err
	}
	if cfg.CreationWindow, err = duration(lookup, "SESSION_CREATION_WINDOW", 10*time.Minute); err != nil {
		return Config{}, err
	}
	if cfg.MaxSourceKeys, err = integer(lookup, "SESSION_MAX_SOURCE_KEYS", 65536, 1, 10000000); err != nil {
		return Config{}, err
	}
	if cfg.OutboxCapacity, err = integer(lookup, "SESSION_EXPIRY_OUTBOX_CAPACITY", 1024, 1, 1000000); err != nil {
		return Config{}, err
	}
	if cfg.RetryMin, err = duration(lookup, "SESSION_EXPIRY_RETRY_MIN", 100*time.Millisecond); err != nil {
		return Config{}, err
	}
	if cfg.RetryMax, err = duration(lookup, "SESSION_EXPIRY_RETRY_MAX", 30*time.Second); err != nil {
		return Config{}, err
	}
	if cfg.RetryMax < cfg.RetryMin {
		return Config{}, fmt.Errorf("SESSION_EXPIRY_RETRY_MAX must be at least SESSION_EXPIRY_RETRY_MIN")
	}
	if cfg.TrustedProxyMaxHops, err = integer(lookup, "TRUSTED_PROXY_MAX_HOPS", 32, 1, 256); err != nil {
		return Config{}, err
	}
	for _, value := range strings.Split(lookup("TRUSTED_PROXY_CIDRS"), ",") {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		prefix, parseErr := netip.ParsePrefix(value)
		if parseErr != nil || prefix != prefix.Masked() {
			return Config{}, fmt.Errorf("TRUSTED_PROXY_CIDRS contains an invalid or non-canonical prefix")
		}
		cfg.TrustedProxyCIDRs = append(cfg.TrustedProxyCIDRs, prefix)
	}
	return cfg, nil
}

func duration(lookup func(string) string, key string, fallback time.Duration) (time.Duration, error) {
	value := lookup(key)
	if value == "" {
		return fallback, nil
	}
	parsed, err := time.ParseDuration(value)
	if err != nil || parsed <= 0 {
		return 0, fmt.Errorf("%s must be a positive Go duration", key)
	}
	return parsed, nil
}

func integer(lookup func(string) string, key string, fallback, min, max int) (int, error) {
	value := lookup(key)
	if value == "" {
		return fallback, nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < min || parsed > max {
		return 0, fmt.Errorf("%s must be between %d and %d", key, min, max)
	}
	return parsed, nil
}
