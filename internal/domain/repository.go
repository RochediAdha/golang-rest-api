package domain

import "context"

type MenuRepository interface {
	Create(ctx context.Context, menu Menu) (Menu, error)
	GetByID(ctx context.Context, id string) (Menu, error)
	GetByCode(ctx context.Context, code string) (Menu, error)
	List(ctx context.Context, filter MenuListFilter) ([]Menu, int, error)
	CountChildren(ctx context.Context, parentID string) (int, error)
	Update(ctx context.Context, menu Menu) (Menu, error)
	Delete(ctx context.Context, menu Menu) error
}

type RoleRepository interface {
	Create(ctx context.Context, role Role) (Role, error)
	GetByID(ctx context.Context, id string) (Role, error)
	GetByName(ctx context.Context, name string) (Role, error)
	List(ctx context.Context, filter ListFilter) ([]Role, int, error)
	Update(ctx context.Context, role Role) (Role, error)
	Delete(ctx context.Context, role Role) error
}

type UserRepository interface {
	Create(ctx context.Context, user User) (User, error)
	GetByID(ctx context.Context, id string) (User, error)
	GetByUsername(ctx context.Context, username string) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	List(ctx context.Context, filter ListFilter) ([]User, int, error)
	Update(ctx context.Context, user User) (User, error)
	Delete(ctx context.Context, id string) error
}

type BookRepository interface {
	Create(ctx context.Context, book Book) (Book, error)
	GetByID(ctx context.Context, id string) (Book, error)
	GetByISBN(ctx context.Context, isbn string) (Book, error)
	List(ctx context.Context, filter ListFilter) ([]Book, int, error)
	Update(ctx context.Context, book Book) (Book, error)
	Delete(ctx context.Context, id string) error
}
