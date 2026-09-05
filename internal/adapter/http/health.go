package httpadapter

import (
	"context"
	"net/http"
	"time"
)

func handleHealth(checker HealthChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := http.StatusOK
		dbStatus := "ok"

		if checker != nil {
			ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
			defer cancel()
			if err := checker.Ping(ctx); err != nil {
				status = http.StatusServiceUnavailable
				dbStatus = "unavailable"
			}
		}

		appStatus := "ok"
		if status != http.StatusOK {
			appStatus = "degraded"
		}

		writeJSON(w, status, map[string]any{
			"status":   appStatus,
			"database": dbStatus,
			"time":     time.Now().UTC().Format(time.RFC3339),
		})
	}
}
