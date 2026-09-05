package domain

import "time"

type Privilege struct {
	ID          string    `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	IsActive    bool      `json:"isActive"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	CreatedBy   *string   `json:"createdBy,omitempty"`
	UpdatedBy   *string   `json:"updatedBy,omitempty"`
}

type CreatePrivilegeInput struct {
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	IsActive    *bool   `json:"isActive"`
	CreatedBy   *string `json:"createdBy"`
}

type UpdatePrivilegeInput struct {
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	IsActive    *bool   `json:"isActive"`
	UpdatedBy   *string `json:"updatedBy"`
}
