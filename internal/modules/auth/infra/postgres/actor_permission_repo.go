package postgres

import (
	"go-enterprise-blueprint/internal/modules/auth/domain/rbac"

	"github.com/rise-and-shine/pkg/repogen"
	"github.com/uptrace/bun"
)

func NewActorPermissionRepo(idb bun.IDB) rbac.ActorPermissionRepo {
	return repogen.NewPgRepoBuilder[rbac.ActorPermission, rbac.ActorPermissionFilter](idb).
		WithSchemaName(schemaName).
		WithNotFoundCode(rbac.CodeActorPermissionNotFound).
		WithFilterFunc(actorPermissionFilterFunc).
		Build()
}

func actorPermissionFilterFunc(q *bun.SelectQuery, f rbac.ActorPermissionFilter) *bun.SelectQuery {
	if f.ID != nil {
		q = q.Where("id = ?", *f.ID)
	}
	if f.ActorType != nil {
		q = q.Where("actor_type = ?", *f.ActorType)
	}
	if f.ActorID != nil {
		q = q.Where("actor_id = ?", *f.ActorID)
	}
	if f.Limit > 0 {
		q = q.Limit(f.Limit)
	}
	if f.Offset > 0 {
		q = q.Offset(f.Offset)
	}
	return q
}
