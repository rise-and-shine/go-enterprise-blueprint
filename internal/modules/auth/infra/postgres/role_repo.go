package postgres

import (
	"go-enterprise-blueprint/internal/modules/auth/domain/rbac"

	"github.com/rise-and-shine/pkg/repogen"
	"github.com/uptrace/bun"
)

const (
	RoleNotFoundCode     = "ROLE_NOT_FOUND"
	RoleNameConflictCode = "ROLE_NAME_EXISTS"
)

func NewRoleRepo(idb bun.IDB) rbac.RoleRepo {
	return repogen.NewPgRepo[rbac.Role, rbac.RoleFilter](
		idb,
		"role",
		RoleNotFoundCode,
		map[string]string{
			"roles_name_key": RoleNameConflictCode,
		},
		roleFilterFunc,
	)
}

func roleFilterFunc(q *bun.SelectQuery, f rbac.RoleFilter) *bun.SelectQuery {
	if f.ID != nil {
		q = q.Where("id = ?", *f.ID)
	}
	if f.Name != nil {
		q = q.Where("name = ?", *f.Name)
	}
	if len(f.IDs) > 0 {
		q = q.Where("id IN (?)", bun.In(f.IDs))
	}
	if f.Limit > 0 {
		q = q.Limit(f.Limit)
	}
	if f.Offset > 0 {
		q = q.Offset(f.Offset)
	}
	return q
}
