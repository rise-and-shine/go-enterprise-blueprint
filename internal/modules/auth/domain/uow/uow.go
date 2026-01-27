package uow

import (
	"context"
	"go-enterprise-blueprint/internal/modules/auth/domain/rbac"
	"go-enterprise-blueprint/internal/modules/auth/domain/session"
	"go-enterprise-blueprint/internal/modules/auth/domain/user"
)

// Factory defines an interface for creating new instances of the UnitOfWork.
type Factory interface {
	// NewUOW creates and returns a new instance of the UnitOfWork.
	NewUOW(ctx context.Context) (UnitOfWork, error)
}

// UnitOfWork represents a single unit of work, typically mapping to a database transaction.
// It provides access to various repositories and methods to finalize or discard changes.
type UnitOfWork interface {
	// Repository accessors
	Role() rbac.RoleRepo
	RolePermission() rbac.RolePermissionRepo
	ActorRole() rbac.ActorRoleRepo
	ActorPermission() rbac.ActorPermissionRepo
	Session() session.Repo
	Admin() user.AdminRepo

	// ApplyChanges finalizes the unit of work, typically committing the underlying transaction.
	// This method doesn't take context.Context, instead should be used context which is used in unit of work creation
	ApplyChanges() error

	// DiscardUnapplied rolls back any pending changes in the unit of work if any error occured until call to Apply method,
	// typically rolling back the transaction.
	DiscardUnapplied()
}
