package domain

import (
	"go-enterprise-blueprint/internal/modules/auth/domain/rbac"
	"go-enterprise-blueprint/internal/modules/auth/domain/session"
	"go-enterprise-blueprint/internal/modules/auth/domain/user"
)

// Container holds domain interfaces.
// It acts as a dependency injection container for the domain layer.
type Container struct {
	adminRepo           user.AdminRepo
	sessionRepo         session.SessionRepo
	roleRepo            rbac.RoleRepo
	rolePermissionRepo  rbac.RolePermissionRepo
	actorRoleRepo       rbac.ActorRoleRepo
	actorPermissionRepo rbac.ActorPermissionRepo
}

func NewContainer(
	adminRepo user.AdminRepo,
	sessionRepo session.SessionRepo,
	roleRepo rbac.RoleRepo,
	rolePermissionRepo rbac.RolePermissionRepo,
	actorRoleRepo rbac.ActorRoleRepo,
	actorPermissionRepo rbac.ActorPermissionRepo,
) *Container {
	return &Container{
		adminRepo:           adminRepo,
		sessionRepo:         sessionRepo,
		roleRepo:            roleRepo,
		rolePermissionRepo:  rolePermissionRepo,
		actorRoleRepo:       actorRoleRepo,
		actorPermissionRepo: actorPermissionRepo,
	}
}

func (c *Container) AdminRepo() user.AdminRepo {
	return c.adminRepo
}

func (c *Container) SessionRepo() session.SessionRepo {
	return c.sessionRepo
}

func (c *Container) RoleRepo() rbac.RoleRepo {
	return c.roleRepo
}

func (c *Container) RolePermissionRepo() rbac.RolePermissionRepo {
	return c.rolePermissionRepo
}

func (c *Container) ActorRoleRepo() rbac.ActorRoleRepo {
	return c.actorRoleRepo
}

func (c *Container) ActorPermissionRepo() rbac.ActorPermissionRepo {
	return c.actorPermissionRepo
}
