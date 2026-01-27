package postgres

import (
	"go-enterprise-blueprint/internal/modules/auth/domain/rbac"

	"github.com/rise-and-shine/pkg/repogen"
	"github.com/uptrace/bun"
)

func NewRoleRepo(idb bun.IDB) rbac.RoleRepo {
	return repogen.NewPgRepoBuilder[rbac.Role, rbac.RoleFilter](idb).
		WithSchemaName(schemaName).
		WithNotFoundCode(rbac.CodeRoleNotFound).
		WithConflictCodesMap(map[string]string{
			"roles_name_key": rbac.CodeRoleNameConflict,
		}).
		WithFilterFunc(roleFilterFunc).
		Build()
}

func roleFilterFunc(q *bun.SelectQuery, f rbac.RoleFilter) *bun.SelectQuery {
	if f.ID != nil {
		q = q.Where("id = ?", *f.ID)
	}
	if f.ActorType != nil {
		q = q.Where("actor_type = ?", *f.ActorType)
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
