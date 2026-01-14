package usecase

import (
	"go-enterprise-blueprint/internal/modules/auth/usecase/admin/admin_login"
	"go-enterprise-blueprint/internal/modules/auth/usecase/admin/admin_logout"
	"go-enterprise-blueprint/internal/modules/auth/usecase/admin/admin_refresh_token"
	"go-enterprise-blueprint/internal/modules/auth/usecase/admin/create_admin"
	"go-enterprise-blueprint/internal/modules/auth/usecase/admin/disable_admin"
	"go-enterprise-blueprint/internal/modules/auth/usecase/admin/get_admins"
	"go-enterprise-blueprint/internal/modules/auth/usecase/admin/update_admin"
	"go-enterprise-blueprint/internal/modules/auth/usecase/rbac/get_actor_permissions"
	"go-enterprise-blueprint/internal/modules/auth/usecase/rbac/get_actor_roles"
	"go-enterprise-blueprint/internal/modules/auth/usecase/rbac/get_role_permissions"
	"go-enterprise-blueprint/internal/modules/auth/usecase/rbac/set_actor_permission"
	"go-enterprise-blueprint/internal/modules/auth/usecase/rbac/set_actor_role"
	"go-enterprise-blueprint/internal/modules/auth/usecase/rbac/set_role_permission"
	"go-enterprise-blueprint/internal/modules/auth/usecase/role/create_role"
	"go-enterprise-blueprint/internal/modules/auth/usecase/role/delete_role"
	"go-enterprise-blueprint/internal/modules/auth/usecase/role/get_roles"
	"go-enterprise-blueprint/internal/modules/auth/usecase/role/update_role"
	"go-enterprise-blueprint/internal/modules/auth/usecase/session/delete_user_all_sessions"
	"go-enterprise-blueprint/internal/modules/auth/usecase/session/delete_user_session"
	"go-enterprise-blueprint/internal/modules/auth/usecase/user/create_superadmin"
)

type Container struct {
	createSuperuser       create_superadmin.UseCase
	adminLogin            admin_login.UseCase
	adminRefreshToken     admin_refresh_token.UseCase
	adminLogout           admin_logout.UseCase
	deleteUserSession     delete_user_session.UseCase
	deleteUserAllSessions delete_user_all_sessions.UseCase
	createRole            create_role.UseCase
	updateRole            update_role.UseCase
	deleteRole            delete_role.UseCase
	getRoles              get_roles.UseCase
	setRolePermission     set_role_permission.UseCase
	getRolePermissions    get_role_permissions.UseCase
	setActorRole          set_actor_role.UseCase
	getActorRoles         get_actor_roles.UseCase
	setActorPermission    set_actor_permission.UseCase
	getActorPermissions   get_actor_permissions.UseCase
	createAdmin           create_admin.UseCase
	updateAdmin           update_admin.UseCase
	disableAdmin          disable_admin.UseCase
	getAdmins             get_admins.UseCase
}

func NewContainer(
	createSuperuser create_superadmin.UseCase,
	adminLogin admin_login.UseCase,
	adminRefreshToken admin_refresh_token.UseCase,
	adminLogout admin_logout.UseCase,
	deleteUserSession delete_user_session.UseCase,
	deleteUserAllSessions delete_user_all_sessions.UseCase,
	createRole create_role.UseCase,
	updateRole update_role.UseCase,
	deleteRole delete_role.UseCase,
	getRoles get_roles.UseCase,
	setRolePermission set_role_permission.UseCase,
	getRolePermissions get_role_permissions.UseCase,
	setActorRole set_actor_role.UseCase,
	getActorRoles get_actor_roles.UseCase,
	setActorPermission set_actor_permission.UseCase,
	getActorPermissions get_actor_permissions.UseCase,
	createAdmin create_admin.UseCase,
	updateAdmin update_admin.UseCase,
	disableAdmin disable_admin.UseCase,
	getAdmins get_admins.UseCase,
) *Container {
	return &Container{
		createSuperuser:       createSuperuser,
		adminLogin:            adminLogin,
		adminRefreshToken:     adminRefreshToken,
		adminLogout:           adminLogout,
		deleteUserSession:     deleteUserSession,
		deleteUserAllSessions: deleteUserAllSessions,
		createRole:            createRole,
		updateRole:            updateRole,
		deleteRole:            deleteRole,
		getRoles:              getRoles,
		setRolePermission:     setRolePermission,
		getRolePermissions:    getRolePermissions,
		setActorRole:          setActorRole,
		getActorRoles:         getActorRoles,
		setActorPermission:    setActorPermission,
		getActorPermissions:   getActorPermissions,
		createAdmin:           createAdmin,
		updateAdmin:           updateAdmin,
		disableAdmin:          disableAdmin,
		getAdmins:             getAdmins,
	}
}

func (c *Container) CreateSuperuser() create_superadmin.UseCase {
	return c.createSuperuser
}

func (c *Container) AdminLogin() admin_login.UseCase {
	return c.adminLogin
}

func (c *Container) AdminRefreshToken() admin_refresh_token.UseCase {
	return c.adminRefreshToken
}

func (c *Container) AdminLogout() admin_logout.UseCase {
	return c.adminLogout
}

func (c *Container) DeleteUserSession() delete_user_session.UseCase {
	return c.deleteUserSession
}

func (c *Container) DeleteUserAllSessions() delete_user_all_sessions.UseCase {
	return c.deleteUserAllSessions
}

func (c *Container) CreateRole() create_role.UseCase {
	return c.createRole
}

func (c *Container) UpdateRole() update_role.UseCase {
	return c.updateRole
}

func (c *Container) DeleteRole() delete_role.UseCase {
	return c.deleteRole
}

func (c *Container) GetRoles() get_roles.UseCase {
	return c.getRoles
}

func (c *Container) SetRolePermission() set_role_permission.UseCase {
	return c.setRolePermission
}

func (c *Container) GetRolePermissions() get_role_permissions.UseCase {
	return c.getRolePermissions
}

func (c *Container) SetActorRole() set_actor_role.UseCase {
	return c.setActorRole
}

func (c *Container) GetActorRoles() get_actor_roles.UseCase {
	return c.getActorRoles
}

func (c *Container) SetActorPermission() set_actor_permission.UseCase {
	return c.setActorPermission
}

func (c *Container) GetActorPermissions() get_actor_permissions.UseCase {
	return c.getActorPermissions
}

func (c *Container) CreateAdmin() create_admin.UseCase {
	return c.createAdmin
}

func (c *Container) UpdateAdmin() update_admin.UseCase {
	return c.updateAdmin
}

func (c *Container) DisableAdmin() disable_admin.UseCase {
	return c.disableAdmin
}

func (c *Container) GetAdmins() get_admins.UseCase {
	return c.getAdmins
}
