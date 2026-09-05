package httpadapter

import (
	"errors"
	"io"
	"net/http"
	"strconv"

	"golang-rest-api/internal/domain"
	"golang-rest-api/internal/usecase"
)

type RoleHandler struct {
	uc *usecase.RoleUseCase
}

func NewRoleHandler(uc *usecase.RoleUseCase) *RoleHandler {
	return &RoleHandler{uc: uc}
}

func (h *RoleHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/roles", h.list)
	mux.HandleFunc("POST /api/v1/roles", h.create)
	mux.HandleFunc("GET /api/v1/roles/{id}", h.get)
	mux.HandleFunc("PUT /api/v1/roles/{id}", h.update)
	mux.HandleFunc("DELETE /api/v1/roles/{id}", h.delete)
}

func (h *RoleHandler) list(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	roles, total, err := h.uc.List(r.Context(), domain.ListFilter{
		Query:  r.URL.Query().Get("q"),
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		writeError(w, err)
		return
	}
	if roles == nil {
		roles = []domain.Role{}
	}

	if limit <= 0 {
		limit = 20
	}

	writeJSON(w, http.StatusOK, listResponse[domain.Role]{
		Data: roles,
		Meta: listMeta{Total: total, Limit: limit, Offset: offset},
	})
}

func (h *RoleHandler) get(w http.ResponseWriter, r *http.Request) {
	role, err := h.uc.Get(r.Context(), r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, role)
}

func (h *RoleHandler) create(w http.ResponseWriter, r *http.Request) {
	var input domain.CreateRoleInput
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, domain.ErrInvalidInput)
		return
	}

	role, err := h.uc.Create(r.Context(), input)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, role)
}

func (h *RoleHandler) update(w http.ResponseWriter, r *http.Request) {
	var input domain.UpdateRoleInput
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, domain.ErrInvalidInput)
		return
	}

	role, err := h.uc.Update(r.Context(), r.PathValue("id"), input)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, role)
}

func (h *RoleHandler) delete(w http.ResponseWriter, r *http.Request) {
	var input domain.DeleteRoleInput
	if r.Body != nil && r.Body != http.NoBody {
		if err := decodeJSON(r, &input); err != nil && !errors.Is(err, io.EOF) {
			writeError(w, domain.ErrInvalidInput)
			return
		}
	}

	role, err := h.uc.Delete(r.Context(), r.PathValue("id"), input)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, role)
}
