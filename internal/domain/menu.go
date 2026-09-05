package domain

import "time"

const (
	MenuTypeItem     = "ITEM"
	MenuTypeGroup    = "GROUP"
	MenuTypeCollapse = "COLLAPSE"
)

type Menu struct {
	ID          string     `json:"id"`
	ParentID    *string    `json:"parentId,omitempty"`
	Code        string     `json:"code"`
	Name        string     `json:"name"`
	Path        string     `json:"path,omitempty"`
	Icon        string     `json:"icon,omitempty"`
	Description string     `json:"description,omitempty"`
	SortOrder   int        `json:"sortOrder"`
	Type        string     `json:"type"`
	IsActive    bool       `json:"isActive"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty"`
	CreatedBy   *string    `json:"createdBy,omitempty"`
	UpdatedBy   *string    `json:"updatedBy,omitempty"`
	DeletedBy   *string    `json:"deletedBy,omitempty"`
}

type CreateMenuInput struct {
	ParentID    *string `json:"parentId"`
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Path        string  `json:"path"`
	Icon        string  `json:"icon"`
	Description string  `json:"description"`
	SortOrder   *int    `json:"sortOrder"`
	Type        string  `json:"type"`
	IsActive    *bool   `json:"isActive"`
	CreatedBy   *string `json:"createdBy"`
}

type UpdateMenuInput struct {
	ParentID    *string `json:"parentId"`
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Path        string  `json:"path"`
	Icon        string  `json:"icon"`
	Description string  `json:"description"`
	SortOrder   *int    `json:"sortOrder"`
	Type        string  `json:"type"`
	IsActive    *bool   `json:"isActive"`
	UpdatedBy   *string `json:"updatedBy"`
}

type DeleteMenuInput struct {
	DeletedBy *string `json:"deletedBy"`
}

type MenuListFilter struct {
	Query    string
	ParentID string
	RootOnly bool
	Limit    int
	Offset   int
}
