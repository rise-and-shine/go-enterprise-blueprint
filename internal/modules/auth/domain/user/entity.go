package user

import (
	"time"

	"github.com/rise-and-shine/pkg/pg"
)

const (
	CodeAdminNotFound         = "ADMIN_NOT_FOUND"
	CodeAdminUsernameConflict = "USERNAME_CONFLICT"
)

type Admin struct {
	pg.BaseModel

	ID string `json:"id" bun:"id,pk,autoincrement"`

	Username     string `json:"username"`
	PasswordHash string `json:"-"`

	IsActive     bool       `json:"is_active"`
	LastActiveAt *time.Time `json:"last_active_at"`
}

type ServiceAccount struct {
	pg.BaseModel

	// TODO: implement service account entity if needed in your project
}

type User struct {
	pg.BaseModel

	// TODO: implement user entity if needed in your project
}
