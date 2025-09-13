package common

import "time"

var (
	Version = "1.0.1"
	Commit  = "dev"
)

// Configuration file for http request configuration
type ServerConfig struct {
	Example         string
	Addr            string
	TlsCertFile     string
	TlsKeyFile      string
	EnableH2C       bool
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	MaxHeaderBytes  int
	ShutdownTimeout time.Duration
	LogLevel        string
}
