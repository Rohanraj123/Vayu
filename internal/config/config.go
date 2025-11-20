package config

import "time"

// Root Configuration struct
type Config struct {
	Server  ServerConfig  `yaml:"server"`
	Routes  []RouteConfig `yaml:"routes"`
	Logging LoggingConfig `yaml:"logging"`
	Auth    AuthConfig    `yaml:"auth"`
}

// ─── SERVER CONFIG ───────────────────────────────────────────────────────────────
type ServerConfig struct {
	Port          int           `yaml:"port"`
	ReadTimeout   time.Duration `yaml:"read_timeout"`
	WriteTimeount time.Duration `yaml:"write_timeout"`
	IdleTimeount  time.Duration `yaml:"idle_timeount"`
}

// ─── RATE-LIMIT CONFIG ───────────────────────────────────────────────────────────────
// rules defined by admin
type RateLimitConfig struct {
	Enabled          bool   `yaml:"enabled"`
	RequestPerMinute int    `yaml:"request_per_minute"`
	Burst            int    `yaml:"burst"`
	KeyScope         string `yaml:"key_scope"`
}

// ─── ROUTE CONFIG ───────────────────────────────────────────────────────────────
type RouteConfig struct {
	Path         string          `yaml:"path"`
	Upstream     string          `yaml:"upstream"`
	Methods      []string        `yaml:"methods"`
	AuthRequired bool            `yaml:"auth_required"`
	Service      string          `yaml:"service"`
	RateLimit    RateLimitConfig `yaml:"rate_limit"`
}

// ─── LOGGING CONFIG ───────────────────────────────────────────────────────────────
type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

// ─── AUTH CONFIG ───────────────────────────────────────────────────────────────
type AuthConfig struct {
	Mode    string `yaml:"mode"`
	Enabled bool   `yaml:"enabled"`
}
