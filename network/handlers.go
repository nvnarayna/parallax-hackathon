package network

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

func (s *Server) handleStatus(
	w http.ResponseWriter,
	r *http.Request,
) {
	counts, err := s.db.Counts()
	if err != nil {
		s.writeError(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	s.writeJSON(w, http.StatusOK, map[string]any{
		"running": s.engine.Running(),
		"counts":  counts,
	})
}

func (s *Server) handleLatestHealth(
	w http.ResponseWriter,
	r *http.Request,
) {
	result, err := s.db.LatestHealth()
	if err != nil {
		s.handleDatabaseError(w, err)
		return
	}

	s.writeJSON(w, http.StatusOK, result)
}

func (s *Server) handleHealthHistory(
	w http.ResponseWriter,
	r *http.Request,
) {
	limit := parseLimit(r)

	results, err := s.db.HealthHistory(limit)
	if err != nil {
		s.writeError(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	s.writeJSON(w, http.StatusOK, results)
}

func (s *Server) handleLatestWiFi(
	w http.ResponseWriter,
	r *http.Request,
) {
	result, err := s.db.LatestWiFi()
	if err != nil {
		s.handleDatabaseError(w, err)
		return
	}

	s.writeJSON(w, http.StatusOK, result)
}

func (s *Server) handleWiFiHistory(
	w http.ResponseWriter,
	r *http.Request,
) {
	limit := parseLimit(r)

	results, err := s.db.WiFiHistory(limit)
	if err != nil {
		s.writeError(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	s.writeJSON(w, http.StatusOK, results)
}

func (s *Server) handleLatestSpeedTest(
	w http.ResponseWriter,
	r *http.Request,
) {
	result, err := s.db.LatestSpeedTest()
	if err != nil {
		s.handleDatabaseError(w, err)
		return
	}

	s.writeJSON(w, http.StatusOK, result)
}

func (s *Server) handleSpeedTestHistory(
	w http.ResponseWriter,
	r *http.Request,
) {
	limit := parseLimit(r)

	results, err := s.db.SpeedTestHistory(limit)
	if err != nil {
		s.writeError(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	s.writeJSON(w, http.StatusOK, results)
}

func (s *Server) handleDiagnosticHistory(
	w http.ResponseWriter,
	r *http.Request,
) {
	limit := parseLimit(r)

	results, err := s.db.DiagnosticHistory(limit)
	if err != nil {
		s.writeError(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	s.writeJSON(w, http.StatusOK, results)
}

func (s *Server) handleLatestDiagnostic(
	w http.ResponseWriter,
	r *http.Request,
) {
	diagnostic, err := s.db.LatestDiagnostic()
	if err != nil {
		s.handleDatabaseError(w, err)
		return
	}

	bundle, err := s.db.DiagnosticBundle(diagnostic.ID)
	if err != nil {
		s.writeError(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	s.writeJSON(w, http.StatusOK, bundle)
}

func (s *Server) handleDiagnostic(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := strconv.ParseInt(
		r.PathValue("id"),
		10,
		64,
	)
	if err != nil || id <= 0 {
		s.writeError(
			w,
			http.StatusBadRequest,
			"invalid diagnostic id",
		)
		return
	}

	bundle, err := s.db.DiagnosticBundle(id)
	if err != nil {
		s.handleDatabaseError(w, err)
		return
	}

	s.writeJSON(w, http.StatusOK, bundle)
}

type diagnosticRequest struct {
	Reason   string `json:"reason"`
	Severity string `json:"severity"`
}

func (s *Server) handleRunDiagnostic(
	w http.ResponseWriter,
	r *http.Request,
) {
	var request diagnosticRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		s.writeError(
			w,
			http.StatusBadRequest,
			"invalid json body",
		)
		return
	}

	request.Reason = strings.TrimSpace(request.Reason)
	request.Severity = strings.TrimSpace(request.Severity)

	if request.Reason == "" {
		request.Reason = "manual diagnostic"
	}

	if request.Severity == "" {
		request.Severity = "unknown"
	}

	err := s.engine.RunDiagnostic(
		"manual",
		request.Reason,
		request.Severity,
	)
	if err != nil {
		s.writeError(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	diagnostic, err := s.db.LatestDiagnostic()
	if err != nil {
		s.writeError(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	bundle, err := s.db.DiagnosticBundle(diagnostic.ID)
	if err != nil {
		s.writeError(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	s.writeJSON(w, http.StatusOK, bundle)
}

func (s *Server) handleDatabaseError(
	w http.ResponseWriter,
	err error,
) {
	if errors.Is(err, sql.ErrNoRows) {
		s.writeError(
			w,
			http.StatusNotFound,
			"no data available",
		)
		return
	}

	s.writeError(
		w,
		http.StatusInternalServerError,
		err.Error(),
	)
}

func parseLimit(r *http.Request) int {
	limit := 50

	value := r.URL.Query().Get("limit")
	if value == "" {
		return limit
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return limit
	}

	if parsed < 1 {
		return 1
	}

	if parsed > 1000 {
		return 1000
	}

	return parsed
}