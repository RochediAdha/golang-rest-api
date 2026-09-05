package httpadapter

import (
	"context"
	"net/http"

	"golang-rest-api/internal/usecase"
)

type HealthChecker interface {
	Ping(ctx context.Context) error
}

func NewRouter(bookUC *usecase.BookUseCase, userUC *usecase.UserUseCase, roleUC *usecase.RoleUseCase, menuUC *usecase.MenuUseCase, health HealthChecker) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handleHealth(health))
	NewBookHandler(bookUC).Register(mux)
	NewUserHandler(userUC).Register(mux)
	NewRoleHandler(roleUC).Register(mux)
	NewMenuHandler(menuUC).Register(mux)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			writeJSON(w, http.StatusNotFound, errorResponse{Error: errorBody{
				Code:    "not_found",
				Message: "resource not found",
			}})
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"name":    "golang-rest-api",
			"version": "v1",
			"docs": map[string]string{
				"health": "GET /health",
				"books":  "GET /api/v1/books",
				"users":  "GET /api/v1/users",
				"roles":  "GET /api/v1/roles",
				"menus":  "GET /api/v1/menus",
			},
		})
	})

	return chain(mux, withRecover, withLogging, withCORS)
}
