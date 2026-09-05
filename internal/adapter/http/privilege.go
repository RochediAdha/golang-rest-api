package httpadapter

import (
	"net/http"
	"strconv"

	"golang-rest-api/internal/domain"
	"golang-rest-api/internal/usecase"
)

type PrivilegeHandler struct {
	uc *usecase.PrivilegeUseCase
}

func NewPrivilegeHandler(uc *usecase.PrivilegeUseCase) *PrivilegeHandler {
	return &PrivilegeHandler{uc: uc}
}

func (h *PrivilegeHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/privileges", h.list)
	mux.HandleFunc("POST /api/v1/privileges", h.create)
	mux.HandleFunc("GET /api/v1/privileges/{id}", h.get)
	mux.HandleFunc("PUT /api/v1/privileges/{id}", h.update)
	mux.HandleFunc("DELETE /api/v1/privileges/{id}", h.delete)
}

func (h *PrivilegeHandler) list(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	items, total, err := h.uc.List(r.Context(), domain.ListFilter{
		Query:  r.URL.Query().Get("q"),
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		writeError(w, err)
		return
	}
	if items == nil {
		items = []domain.Privilege{}
	}

	if limit <= 0 {
		limit = 20
	}

	writeJSON(w, http.StatusOK, listResponse[domain.Privilege]{
		Data: items,
		Meta: listMeta{Total: total, Limit: limit, Offset: offset},
	})
}

func (h *PrivilegeHandler) get(w http.ResponseWriter, r *http.Request) {
	item, err := h.uc.Get(r.Context(), r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *PrivilegeHandler) create(w http.ResponseWriter, r *http.Request) {
	var input domain.CreatePrivilegeInput
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

func (h *PrivilegeHandler) update(w http.ResponseWriter, r *http.Request) {
	var input domain.UpdatePrivilegeInput
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, domain.ErrInvalidInput)
		return
	}

	item, err := h.uc.Update(r.Context(), r.PathValue("id"), input)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *PrivilegeHandler) delete(w http.ResponseWriter, r *http.Request) {
	if err := h.uc.Delete(r.Context(), r.PathValue("id")); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
