package domain

import "time"

type Role struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	IsActive    bool       `json:"isActive"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty"`
	CreatedBy   *string    `json:"createdBy,omitempty"`
	UpdatedBy   *string    `json:"updatedBy,omitempty"`
	DeletedBy   *string    `json:"deletedBy,omitempty"`
}

type CreateRoleInput struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	IsActive    *bool   `json:"isActive"`
	CreatedBy   *string `json:"createdBy"`
}

type UpdateRoleInput struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	IsActive    *bool   `json:"isActive"`
	UpdatedBy   *string `json:"updatedBy"`
}

type DeleteRoleInput struct {
	DeletedBy *string `json:"deletedBy"`
}
