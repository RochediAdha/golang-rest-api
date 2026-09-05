package httpadapter

import (
	"net/http"
	"strconv"

	"golang-rest-api/internal/domain"
	"golang-rest-api/internal/usecase"
)

type RolePrivilegeHandler struct {
	uc *usecase.RolePrivilegeUseCase
}

func NewRolePrivilegeHandler(uc *usecase.RolePrivilegeUseCase) *RolePrivilegeHandler {
	return &RolePrivilegeHandler{uc: uc}
}

func (h *RolePrivilegeHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/role-privileges", h.list)
	mux.HandleFunc("POST /api/v1/role-privileges", h.create)
	mux.HandleFunc("GET /api/v1/role-privileges/{id}", h.get)
	mux.HandleFunc("DELETE /api/v1/role-privileges/{id}", h.delete)
}

func (h *RolePrivilegeHandler) list(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	items, total, err := h.uc.List(r.Context(), domain.RolePrivilegeListFilter{
		RoleID:      r.URL.Query().Get("roleId"),
		MenuID:      r.URL.Query().Get("menuId"),
		PrivilegeID: r.URL.Query().Get("privilegeId"),
		Limit:       limit,
		Offset:      offset,
	})
	if err != nil {
		writeError(w, err)
		return
	}
	if items == nil {
		items = []domain.RolePrivilege{}
	}

	if limit <= 0 {
		limit = 20
	}

	writeJSON(w, http.StatusOK, listResponse[domain.RolePrivilege]{
		Data: items,
		Meta: listMeta{Total: total, Limit: limit, Offset: offset},
	})
}

func (h *RolePrivilegeHandler) get(w http.ResponseWriter, r *http.Request) {
	item, err := h.uc.Get(r.Context(), r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *RolePrivilegeHandler) create(w http.ResponseWriter, r *http.Request) {
	var input domain.CreateRolePrivilegeInput
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, domain.ErrInvalidInput)
		return
	}

	item, err := h.uc.Create(r.Context(), input)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (h *RolePrivilegeHandler) delete(w http.ResponseWriter, r *http.Request) {
	if err := h.uc.Delete(r.Context(), r.PathValue("id")); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
