package rbac

import "github.com/rise-and-shine/pkg/repogen"

type RoleFilter struct {
	ID   *int64
	Name *string
	IDs  []int64

	Limit  int
	Offset int
}

type RolePermissionFilter struct {
	ID     *int64
	RoleID *int64

	Limit  int
	Offset int
}

type ActorRoleFilter struct {
	ID        *int64
	ActorType *ActorType
	ActorID   *string
	RoleID    *int64

	Limit  int
	Offset int
}

type ActorPermissionFilter struct {
	ID        *int64
	ActorType *ActorType
	ActorID   *string

	Limit  int
	Offset int
}

type RoleRepo interface {
	repogen.Repo[Role, RoleFilter]
}

type RolePermissionRepo interface {
	repogen.Repo[RolePermission, RolePermissionFilter]
}

type ActorRoleRepo interface {
	repogen.Repo[ActorRole, ActorRoleFilter]
}

type ActorPermissionRepo interface {
	repogen.Repo[ActorPermission, ActorPermissionFilter]
}
