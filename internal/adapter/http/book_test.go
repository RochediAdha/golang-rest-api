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

func TestBookAPI(t *testing.T) {
	handler := NewRouter(usecase.NewBookUseCase(memory.NewBookRepository()), usecase.NewUserUseCase(memory.NewUserRepository()), usecase.NewRoleUseCase(memory.NewRoleRepository()), usecase.NewMenuUseCase(memory.NewMenuRepository()), nil)

	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/books", strings.NewReader(`{"title":"Go in Action","author":"William Kennedy","year":2015}`))
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

	getRec := httptest.NewRecorder()
	handler.ServeHTTP(getRec, httptest.NewRequest(http.MethodGet, "/api/v1/books/"+id, nil))
	if getRec.Code != http.StatusOK {
		t.Fatalf("get status = %d", getRec.Code)
	}

	listRec := httptest.NewRecorder()
	handler.ServeHTTP(listRec, httptest.NewRequest(http.MethodGet, "/api/v1/books?q=go", nil))
	if listRec.Code != http.StatusOK {
		t.Fatalf("list status = %d", listRec.Code)
	}

	updateReq := httptest.NewRequest(http.MethodPut, "/api/v1/books/"+id, strings.NewReader(`{"title":"Go in Action","author":"William Kennedy","year":2016}`))
	updateReq.Header.Set("Content-Type", "application/json")
	updateRec := httptest.NewRecorder()
	handler.ServeHTTP(updateRec, updateReq)
	if updateRec.Code != http.StatusOK {
		t.Fatalf("update status = %d body=%s", updateRec.Code, updateRec.Body.String())
	}

	deleteRec := httptest.NewRecorder()
	handler.ServeHTTP(deleteRec, httptest.NewRequest(http.MethodDelete, "/api/v1/books/"+id, nil))
	if deleteRec.Code != http.StatusNoContent {
		t.Fatalf("delete status = %d", deleteRec.Code)
	}

	missingRec := httptest.NewRecorder()
	handler.ServeHTTP(missingRec, httptest.NewRequest(http.MethodGet, "/api/v1/books/"+id, nil))
	if missingRec.Code != http.StatusNotFound {
		t.Fatalf("missing status = %d", missingRec.Code)
	}
}

func TestHealth(t *testing.T) {
	handler := NewRouter(usecase.NewBookUseCase(memory.NewBookRepository()), usecase.NewUserUseCase(memory.NewUserRepository()), usecase.NewRoleUseCase(memory.NewRoleRepository()), usecase.NewMenuUseCase(memory.NewMenuRepository()), nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/health", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("health status = %d", rec.Code)
	}
}
