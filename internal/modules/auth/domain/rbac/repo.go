package rbac

import "github.com/rise-and-shine/pkg/repogen"

type RoleFilter struct{}

type RolePermissionFilter struct{}

type ActorRoleFilter struct{}

type ActorPermissionFilter struct{}

type RoleRepo interface {
	repogen.Repo[RoleFilter, Role]
}

type RolePermissionRepo interface {
	repogen.Repo[RolePermissionFilter, RolePermission]
}

type ActorRoleRepo interface {
	repogen.Repo[ActorRoleFilter, ActorRole]
}

type ActorPermissionRepo interface {
	repogen.Repo[ActorPermissionFilter, ActorPermission]
}
