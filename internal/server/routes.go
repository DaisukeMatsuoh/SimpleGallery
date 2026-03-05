package server

import (
	"net/http"

	json "github.com/goccy/go-json"
)

type HealthzResponse struct {
	Status string `json:"status"`
}

// registerRoutes registers all routes for the server
func registerRoutes(s *Server) {
	// Health check endpoint
	s.mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		resp := HealthzResponse{Status: "ok"}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			s.log.Error("failed to encode healthz response", "error", err)
		}
	})
}
