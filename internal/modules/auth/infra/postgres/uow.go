package postgres

import (
	"context"
	"database/sql"
	"errors"

	"go-enterprise-blueprint/internal/modules/auth/domain/rbac"
	"go-enterprise-blueprint/internal/modules/auth/domain/session"
	"go-enterprise-blueprint/internal/modules/auth/domain/uow"
	"go-enterprise-blueprint/internal/modules/auth/domain/user"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/observability/logger"
	"github.com/uptrace/bun"
)

func NewUOWFactory(
	db *bun.DB,
) uow.Factory {
	return &factory{
		db,
	}
}

type factory struct {
	db *bun.DB
}

func (f *factory) NewUOW(ctx context.Context) (uow.UnitOfWork, error) {
	tx, err := f.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	// Repositories will be lazily initialized when accessed
	return &pgUOW{
		tx,
	}, nil
}

type pgUOW struct {
	tx bun.Tx
}

func (u *pgUOW) ApplyChanges() error {
	return errx.Wrap(u.tx.Commit())
}

func (u *pgUOW) DiscardUnapplied() {
	err := errx.Wrap(u.tx.Rollback())
	if err == nil || errors.Is(err, sql.ErrTxDone) {
		return
	}
	logger.Named("auth_uow").With("method", "DiscardUnapplied").Warnx(err)
}

func (u *pgUOW) Role() rbac.RoleRepo {
	return NewRoleRepo(u.tx)
}

func (u *pgUOW) RolePermission() rbac.RolePermissionRepo {
	return NewRolePermissionRepo(u.tx)
}

func (u *pgUOW) ActorRole() rbac.ActorRoleRepo {
	return NewActorRoleRepo(u.tx)
}

func (u *pgUOW) ActorPermission() rbac.ActorPermissionRepo {
	return NewActorPermissionRepo(u.tx)
}

func (u *pgUOW) Session() session.Repo {
	return NewSessionRepo(u.tx)
}

func (u *pgUOW) Admin() user.AdminRepo {
	return NewAdminRepo(u.tx)
}
