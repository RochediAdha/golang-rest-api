package httpadapter

import (
	"net/http"
	"strconv"

	"golang-rest-api/internal/domain"
	"golang-rest-api/internal/usecase"
)

type UserHandler struct {
	uc *usecase.UserUseCase
}

func NewUserHandler(uc *usecase.UserUseCase) *UserHandler {
	return &UserHandler{uc: uc}
}

func (h *UserHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/users", h.list)
	mux.HandleFunc("POST /api/v1/users", h.create)
	mux.HandleFunc("GET /api/v1/users/{id}", h.get)
	mux.HandleFunc("PUT /api/v1/users/{id}", h.update)
	mux.HandleFunc("DELETE /api/v1/users/{id}", h.delete)
}

func (h *UserHandler) list(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	users, total, err := h.uc.List(r.Context(), domain.ListFilter{
		Query:  r.URL.Query().Get("q"),
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		writeError(w, err)
		return
	}
	if users == nil {
		users = []domain.User{}
	}

	if limit <= 0 {
		limit = 20
	}

	writeJSON(w, http.StatusOK, listResponse[domain.User]{
		Data: users,
		Meta: listMeta{Total: total, Limit: limit, Offset: offset},
	})
}

func (h *UserHandler) get(w http.ResponseWriter, r *http.Request) {
	user, err := h.uc.Get(r.Context(), r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func (h *UserHandler) create(w http.ResponseWriter, r *http.Request) {
	var input domain.CreateUserInput
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, domain.ErrInvalidInput)
		return
	}

	user, err := h.uc.Create(r.Context(), input)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, user)
}

func (h *UserHandler) update(w http.ResponseWriter, r *http.Request) {
	var input domain.UpdateUserInput
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, domain.ErrInvalidInput)
		return
	}

	user, err := h.uc.Update(r.Context(), r.PathValue("id"), input)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func (h *UserHandler) delete(w http.ResponseWriter, r *http.Request) {
	if err := h.uc.Delete(r.Context(), r.PathValue("id")); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
