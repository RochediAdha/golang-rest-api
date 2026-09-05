package domain

import "errors"

var (
	ErrNotFound          = errors.New("not found")
	ErrInvalidInput      = errors.New("invalid input")
	ErrDuplicateISBN     = errors.New("isbn already exists")
	ErrDuplicateUsername = errors.New("username already exists")
	ErrDuplicateEmail    = errors.New("email already exists")
	ErrDuplicateRoleName = errors.New("role name already exists")
	ErrDuplicateMenuCode = errors.New("menu code already exists")
	ErrInvalidParent     = errors.New("invalid parent")
	ErrMenuHasChildren   = errors.New("menu has children")
)
