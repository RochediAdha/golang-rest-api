package domain

import "time"

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedBy *string   `json:"createdBy,omitempty"`
	UpdatedBy *string   `json:"updatedBy,omitempty"`
}

type CreateUserInput struct {
	Username  string  `json:"username"`
	Email     string  `json:"email"`
	Name      string  `json:"name"`
	IsActive  *bool   `json:"isActive"`
	CreatedBy *string `json:"createdBy"`
}

type UpdateUserInput struct {
	Username  string  `json:"username"`
	Email     string  `json:"email"`
	Name      string  `json:"name"`
	IsActive  *bool   `json:"isActive"`
	UpdatedBy *string `json:"updatedBy"`
}
