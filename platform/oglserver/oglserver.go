// shared/platform/server.go
package oglserver

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/rotisserie/eris"
)

const (
	readHeaderTimeout = 5 * time.Second
	idleTimeout       = 120 * time.Second
	shutdownTimeout   = 30 * time.Second
)

// HTTPServer is a generic worker that manages the lifecycle of an HTTP server.
type HTTPServer struct {
	server  *http.Server
	logger  *slog.Logger
	appName string
}

// NewHTTPServer creates a pre-configured server ready to be started.
// TODO: pass config instead of add, to be able to overwrite default value readerTimout, idleTimeout, etc.
// TODO: add middleware recovery and middlewareLogin and cors (domains from config)
// TODO: add routes for monitoring
func NewHTTPServer(appName, addr string, handler http.Handler, logger *slog.Logger) *HTTPServer {
	return &HTTPServer{
		server: &http.Server{
			Addr:              addr,
			Handler:           handler,
			ReadHeaderTimeout: readHeaderTimeout,
			IdleTimeout:       idleTimeout,
		},
		logger:  logger,
		appName: appName,
	}
}

// Start blocks until the context is canceled, then performs a graceful shutdown.
func (s *HTTPServer) Start(ctx context.Context) error {
	errChan := make(chan error, 1)

	// 1. Start the server in the background
	go func() {
		s.logger.Info("starting http server", "name", s.appName, "addr", s.server.Addr)
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- eris.Wrapf(err, "server %s crashed", s.appName)
		}
		close(errChan)
	}()

	// 2. Wait for either a fatal crash OR a shutdown signal from the Runner
	select {
	case err := <-errChan:
		// The server crashed on startup (e.g., port already in use)
		return err

	case <-ctx.Done():
		// The global context was canceled (Ctrl+C). Time to shut down cleanly.
		s.logger.Info("initiating graceful shutdown", "name", s.appName)

		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		if err := s.server.Shutdown(shutdownCtx); err != nil {
			return eris.Wrapf(err, "failed to shutdown server %s cleanly", s.appName)
		}

		s.logger.Info("server stopped gracefully", "name", s.appName)

		return nil
	}
}
