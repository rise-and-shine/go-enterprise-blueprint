package postgres

import (
	"go-enterprise-blueprint/internal/modules/auth/domain/rbac"

	"github.com/rise-and-shine/pkg/repogen"
	"github.com/uptrace/bun"
)

const (
	ActorRoleNotFoundCode = "ACTOR_ROLE_NOT_FOUND"
)

func NewActorRoleRepo(idb bun.IDB) rbac.ActorRoleRepo {
	return repogen.NewPgRepo[rbac.ActorRole, rbac.ActorRoleFilter](
		idb,
		"actor_role",
		ActorRoleNotFoundCode,
		nil,
		actorRoleFilterFunc,
	)
}

func actorRoleFilterFunc(q *bun.SelectQuery, f rbac.ActorRoleFilter) *bun.SelectQuery {
	if f.ID != nil {
		q = q.Where("id = ?", *f.ID)
	}
	if f.ActorType != nil {
		q = q.Where("actor_type = ?", *f.ActorType)
	}
	if f.ActorID != nil {
		q = q.Where("actor_id = ?", *f.ActorID)
	}
	if f.RoleID != nil {
		q = q.Where("role_id = ?", *f.RoleID)
	}
	if f.Limit > 0 {
		q = q.Limit(f.Limit)
	}
	if f.Offset > 0 {
		q = q.Offset(f.Offset)
	}
	return q
}
