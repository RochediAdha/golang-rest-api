package domain

import "time"

type RolePrivilege struct {
	ID          string    `json:"id"`
	RoleID      string    `json:"roleId"`
	MenuID      string    `json:"menuId"`
	PrivilegeID string    `json:"privilegeId"`
	CreatedAt   time.Time `json:"createdAt"`
	CreatedBy   *string   `json:"createdBy,omitempty"`
}

type CreateRolePrivilegeInput struct {
	RoleID      string  `json:"roleId"`
	MenuID      string  `json:"menuId"`
	PrivilegeID string  `json:"privilegeId"`
	CreatedBy   *string `json:"createdBy"`
}

type RolePrivilegeListFilter struct {
	RoleID      string
	MenuID      string
	PrivilegeID string
	Limit       int
	Offset      int
}
