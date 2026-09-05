package httpadapter

import (
	"errors"
	"io"
	"net/http"
	"strconv"

	"golang-rest-api/internal/domain"
	"golang-rest-api/internal/usecase"
)

type MenuHandler struct {
	uc *usecase.MenuUseCase
}

func NewMenuHandler(uc *usecase.MenuUseCase) *MenuHandler {
	return &MenuHandler{uc: uc}
}

func (h *MenuHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/menus", h.list)
	mux.HandleFunc("POST /api/v1/menus", h.create)
	mux.HandleFunc("GET /api/v1/menus/{id}", h.get)
	mux.HandleFunc("PUT /api/v1/menus/{id}", h.update)
	mux.HandleFunc("DELETE /api/v1/menus/{id}", h.delete)
}

func (h *MenuHandler) list(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	filter := domain.MenuListFilter{
		Query:  r.URL.Query().Get("q"),
		Limit:  limit,
		Offset: offset,
	}
	if values, ok := r.URL.Query()["parentId"]; ok {
		parentID := values[0]
		if parentID == "" || parentID == "null" || parentID == "root" {
			filter.RootOnly = true
		} else {
			filter.ParentID = parentID
		}
	}

	menus, total, err := h.uc.List(r.Context(), filter)
	if err != nil {
		writeError(w, err)
		return
	}
	if menus == nil {
		menus = []domain.Menu{}
	}

	if limit <= 0 {
		limit = 20
	}

	writeJSON(w, http.StatusOK, listResponse[domain.Menu]{
		Data: menus,
		Meta: listMeta{Total: total, Limit: limit, Offset: offset},
	})
}

func (h *MenuHandler) get(w http.ResponseWriter, r *http.Request) {
	menu, err := h.uc.Get(r.Context(), r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, menu)
}

func (h *MenuHandler) create(w http.ResponseWriter, r *http.Request) {
	var input domain.CreateMenuInput
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, domain.ErrInvalidInput)
		return
	}

	menu, err := h.uc.Create(r.Context(), input)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, menu)
}

func (h *MenuHandler) update(w http.ResponseWriter, r *http.Request) {
	var input domain.UpdateMenuInput
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, domain.ErrInvalidInput)
		return
	}

	menu, err := h.uc.Update(r.Context(), r.PathValue("id"), input)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, menu)
}

func (h *MenuHandler) delete(w http.ResponseWriter, r *http.Request) {
	var input domain.DeleteMenuInput
	if r.Body != nil && r.Body != http.NoBody {
		if err := decodeJSON(r, &input); err != nil && !errors.Is(err, io.EOF) {
			writeError(w, domain.ErrInvalidInput)
			return
		}
	}

	menu, err := h.uc.Delete(r.Context(), r.PathValue("id"), input)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, menu)
}
