package httpadapter

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"golang-rest-api/internal/adapter/memory"
	"golang-rest-api/internal/usecase"
)

func TestMenuAPI(t *testing.T) {
	handler := NewRouter(usecase.NewBookUseCase(memory.NewBookRepository()), usecase.NewUserUseCase(memory.NewUserRepository()), usecase.NewRoleUseCase(memory.NewRoleRepository()), usecase.NewMenuUseCase(memory.NewMenuRepository()), nil)

	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/menus", strings.NewReader(`{"code":"dashboard","name":"Dashboard","path":"/dashboard","icon":"home"}`))
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
	if id == "" {
		t.Fatal("missing id")
	}
	if created["type"] != "ITEM" {
		t.Fatalf("expected default type ITEM, got %v", created["type"])
	}

	getRec := httptest.NewRecorder()
	handler.ServeHTTP(getRec, httptest.NewRequest(http.MethodGet, "/api/v1/menus/"+id, nil))
	if getRec.Code != http.StatusOK {
		t.Fatalf("get status = %d", getRec.Code)
	}

	listRec := httptest.NewRecorder()
	handler.ServeHTTP(listRec, httptest.NewRequest(http.MethodGet, "/api/v1/menus?parentId=root", nil))
	if listRec.Code != http.StatusOK {
		t.Fatalf("list status = %d", listRec.Code)
	}

	updateReq := httptest.NewRequest(http.MethodPut, "/api/v1/menus/"+id, strings.NewReader(`{"code":"dashboard","name":"Home","path":"/home","sortOrder":1}`))
	updateReq.Header.Set("Content-Type", "application/json")
	updateRec := httptest.NewRecorder()
	handler.ServeHTTP(updateRec, updateReq)
	if updateRec.Code != http.StatusOK {
		t.Fatalf("update status = %d body=%s", updateRec.Code, updateRec.Body.String())
	}

	deleteRec := httptest.NewRecorder()
	handler.ServeHTTP(deleteRec, httptest.NewRequest(http.MethodDelete, "/api/v1/menus/"+id, nil))
	if deleteRec.Code != http.StatusOK {
		t.Fatalf("delete status = %d body=%s", deleteRec.Code, deleteRec.Body.String())
	}
	var deleted map[string]any
	if err := json.Unmarshal(deleteRec.Body.Bytes(), &deleted); err != nil {
		t.Fatalf("decode delete: %v", err)
	}
	if deleted["deletedAt"] == nil {
		t.Fatalf("expected deletedAt: %s", deleteRec.Body.String())
	}
}
