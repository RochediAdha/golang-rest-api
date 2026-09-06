package httpadapter

import (
	"net/http"
	"strconv"

	"golang-rest-api/internal/domain"
	"golang-rest-api/internal/usecase"
)

type UserRoleHandler struct {
	uc *usecase.UserRoleUseCase
}

func NewUserRoleHandler(uc *usecase.UserRoleUseCase) *UserRoleHandler {
	return &UserRoleHandler{uc: uc}
}

func (h *UserRoleHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/user-roles", h.list)
	mux.HandleFunc("POST /api/v1/user-roles", h.create)
	mux.HandleFunc("GET /api/v1/user-roles/{id}", h.get)
	mux.HandleFunc("DELETE /api/v1/user-roles/{id}", h.delete)
}

func (h *UserRoleHandler) list(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	items, total, err := h.uc.List(r.Context(), domain.UserRoleListFilter{
		UserID: r.URL.Query().Get("userId"),
		RoleID: r.URL.Query().Get("roleId"),
		Number: r.URL.Query().Get("number"),
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		writeError(w, err)
		return
	}
	if items == nil {
		items = []domain.UserRoleListItem{}
	}

	if limit <= 0 {
		limit = 20
	}

	writeJSON(w, http.StatusOK, listResponse[domain.UserRoleListItem]{
		Data: items,
		Meta: listMeta{Total: total, Limit: limit, Offset: offset},
	})
}

func (h *UserRoleHandler) get(w http.ResponseWriter, r *http.Request) {
	item, err := h.uc.Get(r.Context(), r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *UserRoleHandler) create(w http.ResponseWriter, r *http.Request) {
	var input domain.CreateUserRoleInput
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

func (h *UserRoleHandler) delete(w http.ResponseWriter, r *http.Request) {
	if err := h.uc.Delete(r.Context(), r.PathValue("id")); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
