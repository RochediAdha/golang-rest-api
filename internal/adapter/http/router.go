package httpadapter

import (
	"context"
	"net/http"

	"golang-rest-api/internal/usecase"
)

type HealthChecker interface {
	Ping(ctx context.Context) error
}

func NewRouter(userUC *usecase.UserUseCase, roleUC *usecase.RoleUseCase, menuUC *usecase.MenuUseCase, userRoleUC *usecase.UserRoleUseCase, privilegeUC *usecase.PrivilegeUseCase, rolePrivilegeUC *usecase.RolePrivilegeUseCase, health HealthChecker) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handleHealth(health))
	NewUserHandler(userUC).Register(mux)
	NewRoleHandler(roleUC).Register(mux)
	NewMenuHandler(menuUC).Register(mux)
	NewUserRoleHandler(userRoleUC).Register(mux)
	NewPrivilegeHandler(privilegeUC).Register(mux)
	NewRolePrivilegeHandler(rolePrivilegeUC).Register(mux)

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
				"health":         "GET /health",
				"users":          "GET /api/v1/users",
				"roles":          "GET /api/v1/roles",
				"menus":          "GET /api/v1/menus",
				"userRoles":      "GET /api/v1/user-roles",
				"privileges":     "GET /api/v1/privileges",
				"rolePrivileges": "GET /api/v1/role-privileges",
			},
		})
	})

	return chain(mux, withRecover, withLogging, withCORS)
}
