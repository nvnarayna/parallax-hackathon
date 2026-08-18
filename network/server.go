package network

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"wifi-diagnostic/database"
	"wifi-diagnostic/engine"
)

type Server struct {
	httpServer *http.Server
	db         *database.Database
	engine     *engine.Engine
	logger     *log.Logger
}

func New(
	addr string,
	db *database.Database,
	engine *engine.Engine,
	logger *log.Logger,
) *Server {
	if logger == nil {
		logger = log.Default()
	}

	s := &Server{
		db:     db,
		engine: engine,
		logger: logger,
	}

	mux := http.NewServeMux()

	s.registerRoutes(mux)

	s.httpServer = &http.Server{
		Addr:              addr,
		Handler:           loggingMiddleware(logger, mux),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	return s
}

func (s *Server) Start() error {
	s.logger.Printf("http server listening on %s", s.httpServer.Addr)

	err := s.httpServer.ListenAndServe()

	if err == http.ErrServerClosed {
		return nil
	}

	return err
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/status", s.handleStatus)

	mux.HandleFunc("GET /api/health", s.handleLatestHealth)
	mux.HandleFunc("GET /api/health/history", s.handleHealthHistory)

	mux.HandleFunc("GET /api/wifi", s.handleLatestWiFi)
	mux.HandleFunc("GET /api/wifi/history", s.handleWiFiHistory)

	mux.HandleFunc("GET /api/speedtest", s.handleLatestSpeedTest)
	mux.HandleFunc("GET /api/speedtest/history", s.handleSpeedTestHistory)

	mux.HandleFunc("GET /api/diagnostics", s.handleDiagnosticHistory)
	mux.HandleFunc("GET /api/diagnostics/latest", s.handleLatestDiagnostic)
	mux.HandleFunc("GET /api/diagnostics/{id}", s.handleDiagnostic)

	mux.HandleFunc("POST /api/diagnostic", s.handleRunDiagnostic)
}

func (s *Server) writeJSON(
	w http.ResponseWriter,
	status int,
	value any,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(value); err != nil {
		s.logger.Printf("json response error: %v", err)
	}
}

func (s *Server) writeError(
	w http.ResponseWriter,
	status int,
	message string,
) {
	s.writeJSON(w, status, map[string]string{
		"error": message,
	})
}

func loggingMiddleware(
	logger *log.Logger,
	next http.Handler,
) http.Handler {
	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		start := time.Now()

		next.ServeHTTP(w, r)

		logger.Printf(
			"%s %s %s",
			r.Method,
			r.URL.Path,
			time.Since(start),
		)
	})
}
