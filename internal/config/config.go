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

// ─── ROUTE CONFIG ───────────────────────────────────────────────────────────────
type RouteConfig struct {
	Path         string   `yaml:"path"`
	Upstream     string   `yaml:"upstream"`
	Methods      []string `yaml:"methods"`
	AuthRequired bool     `yaml:"auth_required"`
}

// ─── LOGGING CONFIG ───────────────────────────────────────────────────────────────
type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

// ─── AUTH CONFIG ───────────────────────────────────────────────────────────────
type AuthConfig struct {
	Mode         string   `yaml:"mode"`
	ApiKeyHeader string   `yaml:"api_key_header"`
	ValidKeys    []string `yaml:"valid_keys"`
}
