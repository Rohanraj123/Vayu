package listen

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Rohanraj123/vayu/common"
	"github.com/Rohanraj123/vayu/internal/httpserver"
	"github.com/Rohanraj123/vayu/internal/logging"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	serverConfig := &common.ServerConfig{}
	cmd := &cobra.Command{
		Use:          "listen",
		Short:        "subcommand for vayu",
		Long:         "subcommand for vayu that helps listen at configured specifications.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Starting logs
			logr := logging.New(serverConfig.LogLevel)
			logr.Info("Starting vayu",
				slog.String("version", common.Version),
				slog.String("commit", common.Commit),
				slog.String("addr", serverConfig.Addr),
			)

			// new server
			server := httpserver.New(*serverConfig, logr)

			// start
			errCh := make(chan error, 1)
			go func() {
				if err := server.Start(); err != nil {
					errCh <- err
				}
			}()

			// signals
			sigCh := make(chan os.Signal, 1)
			signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

			select {
			case sig := <-sigCh:
				logr.Info("signal received, initiating graceful shutdown", slog.String("signal", sig.String()))
				ctx, cancel := context.WithTimeout(context.Background(), serverConfig.ShutdownTimeout)
				defer cancel()
				if err := server.Shutdown(ctx); err != nil {
					logr.Error("graceful shutdown failed", slog.Any("error", err))
					os.Exit(1)
				}
				logr.Info("shutdown complete")
			case err := <-errCh:
				log.Fatalf("server error: %v", err)
			}

			return nil
		},
	}

	// Add more flags as required
	cmd.Flags().StringVarP(&serverConfig.Addr, "addr", "a", ":8080", "listen address")
	cmd.Flags().StringVar(&serverConfig.TlsCertFile, "tlsCert", "", "path to tls certificate file")
	cmd.Flags().StringVar(&serverConfig.TlsKeyFile, "tlsKey", "", "path to tls key file")
	cmd.Flags().BoolVar(&serverConfig.EnableH2C, "h2c", false, "enable HTTP/2 over clean text (no TLS)")
	cmd.Flags().DurationVar(&serverConfig.ReadTimeout, "read-timeout", 5*time.Second, "read timeout")
	cmd.Flags().DurationVar(&serverConfig.WriteTimeout, "write-timeout", 30*time.Second, "write timeout")
	cmd.Flags().DurationVar(&serverConfig.IdleTimeout, "idle-timeout", 60*time.Second, "idle timeout")
	cmd.Flags().IntVarP(&serverConfig.MaxHeaderBytes, "max-header-bytes", "", 1<<20, "max header bytes")
	cmd.Flags().DurationVar(&serverConfig.ShutdownTimeout, "shutdown-timeout", 20*time.Second, "graceful shutdown timeout")
	cmd.Flags().StringVarP(&serverConfig.LogLevel, "log-level", "", "info", "log level: debug|info|warn|error")

	return cmd

}
