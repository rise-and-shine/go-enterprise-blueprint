package rbac

import (
	"slices"

	"github.com/rise-and-shine/pkg/pg"
)

const (
	CodeRoleNotFound            = "ROLE_NOT_FOUND"
	CodeRoleNameConflict        = "ROLE_NAME_CONFLICT"
	CodeRolePermissionNotFound  = "ROLE_PERMISSION_NOT_FOUND"
	CodeActorRoleNotFound       = "ACTOR_ROLE_NOT_FOUND"
	CodeActorPermissionNotFound = "ACTOR_PERMISSION_NOT_FOUND"
)

type ActorType string

func (at ActorType) IsValid() bool {
	valids := []ActorType{ActorTypeUser, ActorTypeAdmin, ActorTypeServiceAcc}

	return slices.Contains(valids, at)
}

const (
	ActorTypeUser       ActorType = "user"
	ActorTypeAdmin      ActorType = "admin"
	ActorTypeServiceAcc ActorType = "service_acc"
)

type Role struct {
	pg.BaseModel

	ID int64 `json:"id" bun:"id,pk,autoincrement"`

	ActorType ActorType `json:"actor_type"`

	// Name is a unique name of the role
	Name string `json:"name"`
}

type RolePermission struct {
	pg.BaseModel

	ID int64 `json:"id" bun:"id,pk,autoincrement"`

	RoleID     int64  `json:"role_id"`
	Permission string `json:"permission"`
}

type ActorRole struct {
	pg.BaseModel

	ID int64 `json:"id" bun:"id,pk,autoincrement"`

	ActorType ActorType `json:"actor_type"`
	ActorID   string    `json:"actor_id"`

	RoleID int64 `json:"role_id"`
}

type ActorPermission struct {
	pg.BaseModel

	ID int64 `json:"id" bun:"id,pk,autoincrement"`

	ActorType ActorType `json:"actor_type"`
	ActorID   string    `json:"actor_id"`

	Permission string `json:"permission"`
}
