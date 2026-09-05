package domain

import "errors"

var (
	ErrNotFound                = errors.New("not found")
	ErrInvalidInput            = errors.New("invalid input")
	ErrDuplicateUsername       = errors.New("username already exists")
	ErrDuplicateEmail          = errors.New("email already exists")
	ErrDuplicateRoleName       = errors.New("role name already exists")
	ErrDuplicateMenuCode       = errors.New("menu code already exists")
	ErrInvalidParent           = errors.New("invalid parent")
	ErrMenuHasChildren         = errors.New("menu has children")
	ErrInvalidUser             = errors.New("invalid user")
	ErrInvalidRole             = errors.New("invalid role")
	ErrDuplicateUserRole       = errors.New("user role already exists")
	ErrDuplicateUserRoleNumber = errors.New("user role number already exists")
	ErrDuplicatePrivilegeCode  = errors.New("privilege code already exists")
	ErrInvalidMenu             = errors.New("invalid menu")
	ErrInvalidPrivilege        = errors.New("invalid privilege")
	ErrDuplicateRolePrivilege  = errors.New("role privilege already exists")
)
