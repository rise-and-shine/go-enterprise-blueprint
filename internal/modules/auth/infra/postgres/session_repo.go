package postgres

import (
	"go-enterprise-blueprint/internal/modules/auth/domain/session"

	"github.com/rise-and-shine/pkg/repogen"
	"github.com/uptrace/bun"
)

const (
	SessionNotFoundCode = "SESSION_NOT_FOUND"
)

func NewSessionRepo(idb bun.IDB) session.Repo {
	return repogen.NewPgRepo[session.Session, session.Filter](
		idb,
		"session",
		SessionNotFoundCode,
		nil,
		sessionFilterFunc,
	)
}

func sessionFilterFunc(q *bun.SelectQuery, f session.Filter) *bun.SelectQuery {
	if f.ID != nil {
		q = q.Where("id = ?", *f.ID)
	}
	if f.ActorType != nil {
		q = q.Where("actor_type = ?", *f.ActorType)
	}
	if f.ActorID != nil {
		q = q.Where("actor_id = ?", *f.ActorID)
	}
	if f.AccessToken != nil {
		q = q.Where("access_token = ?", *f.AccessToken)
	}
	if f.RefreshToken != nil {
		q = q.Where("refresh_token = ?", *f.RefreshToken)
	}
	if f.Limit > 0 {
		q = q.Limit(f.Limit)
	}
	if f.Offset > 0 {
		q = q.Offset(f.Offset)
	}
	return q
}
