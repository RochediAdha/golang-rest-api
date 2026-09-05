package httpadapter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRolePrivilegeAPI(t *testing.T) {
	handler := testRouter()

	roleID := createJSON(t, handler, "/api/v1/roles", `{"name":"admin"}`)
	menuID := createJSON(t, handler, "/api/v1/menus", `{"code":"users","name":"Users"}`)
	privilegeID := createJSON(t, handler, "/api/v1/privileges", `{"code":"VIEW","name":"View"}`)

	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/role-privileges", strings.NewReader(fmt.Sprintf(`{"roleId":%q,"menuId":%q,"privilegeId":%q}`, roleID, menuID, privilegeID)))
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
	if id == "" || created["roleId"] != roleID || created["menuId"] != menuID || created["privilegeId"] != privilegeID {
		t.Fatalf("unexpected create body: %s", createRec.Body.String())
	}

	getRec := httptest.NewRecorder()
	handler.ServeHTTP(getRec, httptest.NewRequest(http.MethodGet, "/api/v1/role-privileges/"+id, nil))
	if getRec.Code != http.StatusOK {
		t.Fatalf("get status = %d", getRec.Code)
	}

	listRec := httptest.NewRecorder()
	handler.ServeHTTP(listRec, httptest.NewRequest(http.MethodGet, "/api/v1/role-privileges?roleId="+roleID, nil))
	if listRec.Code != http.StatusOK {
		t.Fatalf("list status = %d", listRec.Code)
	}

	dupReq := httptest.NewRequest(http.MethodPost, "/api/v1/role-privileges", strings.NewReader(fmt.Sprintf(`{"roleId":%q,"menuId":%q,"privilegeId":%q}`, roleID, menuID, privilegeID)))
	dupReq.Header.Set("Content-Type", "application/json")
	dupRec := httptest.NewRecorder()
	handler.ServeHTTP(dupRec, dupReq)
	if dupRec.Code != http.StatusConflict {
		t.Fatalf("duplicate status = %d body=%s", dupRec.Code, dupRec.Body.String())
	}

	deleteRec := httptest.NewRecorder()
	handler.ServeHTTP(deleteRec, httptest.NewRequest(http.MethodDelete, "/api/v1/role-privileges/"+id, nil))
	if deleteRec.Code != http.StatusNoContent {
		t.Fatalf("delete status = %d body=%s", deleteRec.Code, deleteRec.Body.String())
	}

	missingRec := httptest.NewRecorder()
	handler.ServeHTTP(missingRec, httptest.NewRequest(http.MethodGet, "/api/v1/role-privileges/"+id, nil))
	if missingRec.Code != http.StatusNotFound {
		t.Fatalf("missing status = %d", missingRec.Code)
	}
}
