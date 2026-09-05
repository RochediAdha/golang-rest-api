package domain

import "time"

type UserRole struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	RoleID    string    `json:"roleId"`
	Number    string    `json:"number,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	CreatedBy *string   `json:"createdBy,omitempty"`
}

type CreateUserRoleInput struct {
	UserID    string  `json:"userId"`
	RoleID    string  `json:"roleId"`
	Number    string  `json:"number"`
	CreatedBy *string `json:"createdBy"`
}

type UserRoleListFilter struct {
	UserID string
	RoleID string
	Number string
	Limit  int
	Offset int
}
