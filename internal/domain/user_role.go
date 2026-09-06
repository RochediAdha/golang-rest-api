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

type UserRoleItem struct {
	ID        string    `json:"id"`
	RoleID    string    `json:"roleId"`
	Name      string    `json:"name"`
	Number    string    `json:"number"`
	CreatedAt time.Time `json:"createdAt"`
	CreatedBy *string   `json:"createdBy,omitempty"`
}

type UserRoleView struct {
	ID     string         `json:"id"`
	UserID string         `json:"userId"`
	Name   string         `json:"name"`
	Roles  []UserRoleItem `json:"roles"`
}

type ActorRef struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type UserRoleListItem struct {
	ID              string    `json:"id"`
	UserID          string    `json:"userId"`
	UserName        string    `json:"userName"`
	RoleID          string    `json:"roleId"`
	RoleName        string    `json:"roleName"`
	RoleDescription string    `json:"roleDescription,omitempty"`
	Number          string    `json:"number"`
	CreatedAt       time.Time `json:"createdAt"`
	CreatedBy       *ActorRef `json:"createdBy,omitempty"`
}
