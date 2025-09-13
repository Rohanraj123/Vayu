package httpserver

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/Rohanraj123/vayu/common"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type atomicBool struct {
	v uint32
}

func (a *atomicBool) set(val bool) {
	if val {
		atomic.StoreUint32(&a.v, 1)
	} else {
		atomic.StoreUint32(&a.v, 0)
	}
}

func (a *atomicBool) get() bool {
	return atomic.LoadUint32(&a.v) == 1
}

type Server struct {
	cfg        common.ServerConfig
	httpServer *http.Server
	ready      atomicBool
	logger     *slog.Logger
}

func New(cfg common.ServerConfig, lgr *slog.Logger) *Server {
	server := &Server{cfg: cfg, logger: lgr}
	server.ready.set(true) // ready once listening; flips to false if shutdown
	return server
}

func (s *Server) Start() error {
	// Base handler mux
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.root)
	mux.HandleFunc("/healthz", s.healthz)
	mux.HandleFunc("/readyz", s.readyz)
	mux.HandleFunc("/version", s.version)

	// middleware creation
	handler := s.accessLog(mux)

	// Enable h2c (HTTP/2 cleartext) if present
	if s.cfg.EnableH2C {
		h2s := &http2.Server{}
		handler = h2c.NewHandler(handler, h2s)
	}

	// populating the server
	s.httpServer = &http.Server{
		Addr:              s.cfg.Addr,
		Handler:           handler,
		ReadTimeout:       s.cfg.ReadTimeout,
		WriteTimeout:      s.cfg.WriteTimeout,
		IdleTimeout:       s.cfg.IdleTimeout,
		MaxHeaderBytes:    s.cfg.MaxHeaderBytes,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// If TLS provided, configure TLS + HTTP/2
	if s.cfg.TlsCertFile != "" && s.cfg.TlsKeyFile != "" {
		tlsConfig := &tls.Config{
			MinVersion: tls.VersionTLS13,
			NextProtos: []string{"h2", "http/1.1"},
		}
		s.httpServer.TLSConfig = tlsConfig

		if err := http2.ConfigureServer(s.httpServer, &http2.Server{}); err != nil {
			return fmt.Errorf("configure http2: %w", err)
		}
		s.logger.Info("vayu listening (TLS)...", slog.String("Addr", s.cfg.Addr), slog.Bool("h2", true))
		return s.httpServer.ListenAndServeTLS(s.cfg.TlsCertFile, s.cfg.TlsKeyFile)
	}

	s.logger.Info("vayu listening", slog.String("addr", s.cfg.Addr), slog.Bool("h2c", s.cfg.EnableH2C))
	return s.httpServer.ListenAndServe()
}

// shuts down the server.
func (s *Server) Shutdown(ctx context.Context) error {
	s.ready.set(false)
	return s.httpServer.Shutdown(ctx)
}

// handlers
func (s *Server) root(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "vayu")
	fmt.Fprintln(w, "vayu is breathing. See /healthz, /root, /version")
}

func (s *Server) healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func (s *Server) readyz(w http.ResponseWriter, r *http.Request) {
	if s.ready.get() {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ready"))
		return
	}
	w.WriteHeader(http.StatusServiceUnavailable)
	w.Write([]byte("draining"))
}

func (s *Server) version(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(common.Version))
}

// helpers
func (s *Server) accessLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		aw := &statusWriter{ResponseWriter: w, status: 200}
		next.ServeHTTP(aw, r)
		remote := remoteIp(r)
		s.logger.Info("request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", 200),
			slog.Duration("duration", time.Since(start)),
			slog.String("remote_ip", remote),
			slog.String("proto", r.Proto),
			slog.String("ua", r.UserAgent()),
		)
	})
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func remoteIp(r *http.Request) string {
	// Best-Effor: Honor X-Forwarded-For first
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
