package httpadapter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"golang-rest-api/internal/adapter/memory"
	"golang-rest-api/internal/usecase"
)

func TestUserRoleAPI(t *testing.T) {
	userRepo := memory.NewUserRepository()
	roleRepo := memory.NewRoleRepository()
	menuRepo := memory.NewMenuRepository()
	privilegeRepo := memory.NewPrivilegeRepository()
	handler := NewRouter(
		usecase.NewUserUseCase(userRepo),
		usecase.NewRoleUseCase(roleRepo),
		usecase.NewMenuUseCase(menuRepo),
		usecase.NewUserRoleUseCase(memory.NewUserRoleRepository(), userRepo, roleRepo),
		usecase.NewPrivilegeUseCase(privilegeRepo),
		usecase.NewRolePrivilegeUseCase(memory.NewRolePrivilegeRepository(), roleRepo, menuRepo, privilegeRepo),
		nil,
	)

	userID := createJSON(t, handler, "/api/v1/users", `{"username":"rochedi","email":"rochedi@example.com","name":"Rochedi"}`)
	roleID := createJSON(t, handler, "/api/v1/roles", `{"name":"admin"}`)

	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/user-roles", strings.NewReader(fmt.Sprintf(`{"userId":%q,"roleId":%q}`, userID, roleID)))
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()
	handler.ServeHTTP(createRec, createReq)
	if createRec.Code != http.StatusCreated {
		t.Fatalf("create status = %d body=%s", createRec.Code, createRec.Body.String())
	}

	var created map[string]any
	if err := json.Unmarshal(createRec.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode create: %v", err)
	}
	id, _ := created["id"].(string)
	if id == "" || created["userId"] != userID || created["roleId"] != roleID {
		t.Fatalf("unexpected create body: %s", createRec.Body.String())
	}
	if _, ok := created["number"]; ok {
		t.Fatalf("admin assignment should omit number: %s", createRec.Body.String())
	}

	studentID := createJSON(t, handler, "/api/v1/roles", `{"name":"mahasiswa"}`)
	studentReq := httptest.NewRequest(http.MethodPost, "/api/v1/user-roles", strings.NewReader(fmt.Sprintf(`{"userId":%q,"roleId":%q,"number":"2301001"}`, userID, studentID)))
	studentReq.Header.Set("Content-Type", "application/json")
	studentRec := httptest.NewRecorder()
	handler.ServeHTTP(studentRec, studentReq)
	if studentRec.Code != http.StatusCreated {
		t.Fatalf("student create status = %d body=%s", studentRec.Code, studentRec.Body.String())
	}
	numberRec := httptest.NewRecorder()
	handler.ServeHTTP(numberRec, httptest.NewRequest(http.MethodGet, "/api/v1/user-roles?number=2301001", nil))
	if numberRec.Code != http.StatusOK {
		t.Fatalf("list by number status = %d", numberRec.Code)
	}

	getRec := httptest.NewRecorder()
	handler.ServeHTTP(getRec, httptest.NewRequest(http.MethodGet, "/api/v1/user-roles/"+userID, nil))
	if getRec.Code != http.StatusOK {
		t.Fatalf("get status = %d body=%s", getRec.Code, getRec.Body.String())
	}
	var shown map[string]any
	if err := json.Unmarshal(getRec.Body.Bytes(), &shown); err != nil {
		t.Fatalf("decode get: %v", err)
	}
	if shown["id"] != userID || shown["userId"] != userID || shown["name"] != "Rochedi" {
		t.Fatalf("unexpected show body: %s", getRec.Body.String())
	}
	roles, _ := shown["roles"].([]any)
	if len(roles) != 2 {
		t.Fatalf("expected 2 roles, got %s", getRec.Body.String())
	}

	listRec := httptest.NewRecorder()
	handler.ServeHTTP(listRec, httptest.NewRequest(http.MethodGet, "/api/v1/user-roles?userId="+userID, nil))
	if listRec.Code != http.StatusOK {
		t.Fatalf("list status = %d", listRec.Code)
	}

	dupReq := httptest.NewRequest(http.MethodPost, "/api/v1/user-roles", strings.NewReader(fmt.Sprintf(`{"userId":%q,"roleId":%q}`, userID, roleID)))
	dupReq.Header.Set("Content-Type", "application/json")
	dupRec := httptest.NewRecorder()
	handler.ServeHTTP(dupRec, dupReq)
	if dupRec.Code != http.StatusConflict {
		t.Fatalf("duplicate status = %d body=%s", dupRec.Code, dupRec.Body.String())
	}

	deleteRec := httptest.NewRecorder()
	handler.ServeHTTP(deleteRec, httptest.NewRequest(http.MethodDelete, "/api/v1/user-roles/"+id, nil))
	if deleteRec.Code != http.StatusNoContent {
		t.Fatalf("delete status = %d body=%s", deleteRec.Code, deleteRec.Body.String())
	}

	missingRec := httptest.NewRecorder()
	handler.ServeHTTP(missingRec, httptest.NewRequest(http.MethodGet, "/api/v1/user-roles/"+id, nil))
	if missingRec.Code != http.StatusNotFound {
		t.Fatalf("missing status = %d", missingRec.Code)
	}
}

func createJSON(t *testing.T, handler http.Handler, path, body string) string {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create %s status = %d body=%s", path, rec.Code, rec.Body.String())
	}
	var payload map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode %s: %v", path, err)
	}
	id, _ := payload["id"].(string)
	if id == "" {
		t.Fatalf("missing id from %s", path)
	}
	return id
}
