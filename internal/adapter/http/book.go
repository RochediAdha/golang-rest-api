package httpadapter

import (
	"encoding/json"
	"net/http"
	"strconv"

	"golang-rest-api/internal/domain"
	"golang-rest-api/internal/usecase"
)

type BookHandler struct {
	uc *usecase.BookUseCase
}

func NewBookHandler(uc *usecase.BookUseCase) *BookHandler {
	return &BookHandler{uc: uc}
}

func (h *BookHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/books", h.list)
	mux.HandleFunc("POST /api/v1/books", h.create)
	mux.HandleFunc("GET /api/v1/books/{id}", h.get)
	mux.HandleFunc("PUT /api/v1/books/{id}", h.update)
	mux.HandleFunc("DELETE /api/v1/books/{id}", h.delete)
}

func (h *BookHandler) list(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	books, total, err := h.uc.List(r.Context(), domain.ListFilter{
		Query:  r.URL.Query().Get("q"),
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		writeError(w, err)
		return
	}
	if books == nil {
		books = []domain.Book{}
	}

	if limit <= 0 {
		limit = 20
	}

	writeJSON(w, http.StatusOK, listResponse[domain.Book]{
		Data: books,
		Meta: listMeta{Total: total, Limit: limit, Offset: offset},
	})
}

func (h *BookHandler) get(w http.ResponseWriter, r *http.Request) {
	book, err := h.uc.Get(r.Context(), r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, book)
}

func (h *BookHandler) create(w http.ResponseWriter, r *http.Request) {
	var input domain.CreateBookInput
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, domain.ErrInvalidInput)
		return
	}

	book, err := h.uc.Create(r.Context(), input)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, book)
}

func (h *BookHandler) update(w http.ResponseWriter, r *http.Request) {
	var input domain.UpdateBookInput
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, domain.ErrInvalidInput)
		return
	}

	book, err := h.uc.Update(r.Context(), r.PathValue("id"), input)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, book)
}

func (h *BookHandler) delete(w http.ResponseWriter, r *http.Request) {
	if err := h.uc.Delete(r.Context(), r.PathValue("id")); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func decodeJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}
