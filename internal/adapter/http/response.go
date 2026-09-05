package httpadapter

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"golang-rest-api/internal/domain"
)

type errorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type errorResponse struct {
	Error errorBody `json:"error"`
}

type listMeta struct {
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type listResponse[T any] struct {
	Data []T      `json:"data"`
	Meta listMeta `json:"meta"`
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if payload == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		slog.Error("encode json response", "err", err)
	}
}

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: errorBody{
			Code:    "invalid_input",
			Message: "request body or parameters are invalid",
		}})
	case errors.Is(err, domain.ErrNotFound):
		writeJSON(w, http.StatusNotFound, errorResponse{Error: errorBody{
			Code:    "not_found",
			Message: "resource not found",
		}})
	case errors.Is(err, domain.ErrDuplicateISBN):
		writeJSON(w, http.StatusConflict, errorResponse{Error: errorBody{
			Code:    "duplicate_isbn",
			Message: "a book with this isbn already exists",
		}})
	case errors.Is(err, domain.ErrDuplicateUsername):
		writeJSON(w, http.StatusConflict, errorResponse{Error: errorBody{
			Code:    "duplicate_username",
			Message: "a user with this username already exists",
		}})
	case errors.Is(err, domain.ErrDuplicateEmail):
		writeJSON(w, http.StatusConflict, errorResponse{Error: errorBody{
			Code:    "duplicate_email",
			Message: "a user with this email already exists",
		}})
	case errors.Is(err, domain.ErrDuplicateRoleName):
		writeJSON(w, http.StatusConflict, errorResponse{Error: errorBody{
			Code:    "duplicate_role_name",
			Message: "a role with this name already exists",
		}})
	case errors.Is(err, domain.ErrDuplicateMenuCode):
		writeJSON(w, http.StatusConflict, errorResponse{Error: errorBody{
			Code:    "duplicate_menu_code",
			Message: "a menu with this code already exists",
		}})
	case errors.Is(err, domain.ErrInvalidParent):
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: errorBody{
			Code:    "invalid_parent",
			Message: "parent menu is invalid",
		}})
	case errors.Is(err, domain.ErrMenuHasChildren):
		writeJSON(w, http.StatusConflict, errorResponse{Error: errorBody{
			Code:    "menu_has_children",
			Message: "menu still has child menus",
		}})
	default:
		slog.Error("unhandled error", "err", err)
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: errorBody{
			Code:    "internal_error",
			Message: "internal server error",
		}})
	}
}
