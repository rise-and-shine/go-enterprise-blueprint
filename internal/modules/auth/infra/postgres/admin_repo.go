package postgres

import (
	"go-enterprise-blueprint/internal/modules/auth/domain/user"

	"github.com/rise-and-shine/pkg/repogen"
	"github.com/uptrace/bun"
)

const (
	AdminNotFoundCode     = "ADMIN_NOT_FOUND"
	AdminUsernameConflict = "USERNAME_EXISTS"
)

func NewAdminRepo(idb bun.IDB) user.AdminRepo {
	return repogen.NewPgRepo[user.Admin, user.AdminFilter](
		idb,
		"admin",
		AdminNotFoundCode,
		map[string]string{
			"admins_username_key": AdminUsernameConflict,
		},
		adminFilterFunc,
	)
}

func adminFilterFunc(q *bun.SelectQuery, f user.AdminFilter) *bun.SelectQuery {
	if f.ID != nil {
		q = q.Where("id = ?", *f.ID)
	}
	if f.Username != nil {
		q = q.Where("username = ?", *f.Username)
	}
	if f.IsSuperadmin != nil {
		q = q.Where("is_superadmin = ?", *f.IsSuperadmin)
	}
	if f.IsActive != nil {
		q = q.Where("is_active = ?", *f.IsActive)
	}
	if f.Limit > 0 {
		q = q.Limit(f.Limit)
	}
	if f.Offset > 0 {
		q = q.Offset(f.Offset)
	}
	return q
}
