package httpadapter

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"golang-rest-api/internal/adapter/memory"
	"golang-rest-api/internal/usecase"
)

func testRouter() http.Handler {
	userRepo := memory.NewUserRepository()
	roleRepo := memory.NewRoleRepository()
	menuRepo := memory.NewMenuRepository()
	privilegeRepo := memory.NewPrivilegeRepository()
	return NewRouter(
		usecase.NewUserUseCase(userRepo),
		usecase.NewRoleUseCase(roleRepo),
		usecase.NewMenuUseCase(menuRepo),
		usecase.NewUserRoleUseCase(memory.NewUserRoleRepository(), userRepo, roleRepo),
		usecase.NewPrivilegeUseCase(privilegeRepo),
		usecase.NewRolePrivilegeUseCase(memory.NewRolePrivilegeRepository(), roleRepo, menuRepo, privilegeRepo),
		nil,
	)
}

func TestHealth(t *testing.T) {
	handler := testRouter()
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/health", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("health status = %d", rec.Code)
	}
}
