package user

import (
	"time"

	"github.com/rise-and-shine/pkg/pg"
)

type Admin struct {
	pg.BaseModel

	ID           string `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`

	IsSuperadmin bool       `json:"is_superadmin"`
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
