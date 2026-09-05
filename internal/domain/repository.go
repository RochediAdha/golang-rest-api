package domain

import "context"

type RolePrivilegeRepository interface {
	Create(ctx context.Context, rolePrivilege RolePrivilege) (RolePrivilege, error)
	GetByID(ctx context.Context, id string) (RolePrivilege, error)
	GetByRoleMenuPrivilege(ctx context.Context, roleID, menuID, privilegeID string) (RolePrivilege, error)
	List(ctx context.Context, filter RolePrivilegeListFilter) ([]RolePrivilege, int, error)
	Delete(ctx context.Context, id string) error
}

type PrivilegeRepository interface {
	Create(ctx context.Context, privilege Privilege) (Privilege, error)
	GetByID(ctx context.Context, id string) (Privilege, error)
	GetByCode(ctx context.Context, code string) (Privilege, error)
	List(ctx context.Context, filter ListFilter) ([]Privilege, int, error)
	Update(ctx context.Context, privilege Privilege) (Privilege, error)
	Delete(ctx context.Context, id string) error
}

type UserRoleRepository interface {
	Create(ctx context.Context, userRole UserRole) (UserRole, error)
	GetByID(ctx context.Context, id string) (UserRole, error)
	GetByUserAndRole(ctx context.Context, userID, roleID string) (UserRole, error)
	GetByRoleAndNumber(ctx context.Context, roleID, number string) (UserRole, error)
	List(ctx context.Context, filter UserRoleListFilter) ([]UserRole, int, error)
	Delete(ctx context.Context, id string) error
}

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
