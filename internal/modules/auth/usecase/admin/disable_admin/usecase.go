package disable_admin

import (
	"context"

	"go-enterprise-blueprint/internal/modules/auth/domain"
	"go-enterprise-blueprint/internal/modules/auth/domain/rbac"
	"go-enterprise-blueprint/internal/modules/auth/domain/session"
	"go-enterprise-blueprint/internal/modules/auth/domain/user"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/ucdef"
)

const (
	OperationID = "disable-admin"

	CodeAdminNotFound               = "ADMIN_NOT_FOUND"
	CodeAdminAlreadyDisabled        = "ADMIN_ALREADY_DISABLED"
	CodeCannotDisableLastSuperadmin = "CANNOT_DISABLE_LAST_SUPERADMIN"
)

type Input struct {
	ID string `json:"id" validate:"required,uuid"`
}

type Output struct {
	ID                 string `json:"id"`
	Username           string `json:"username"`
	IsActive           bool   `json:"is_active"`
	SessionsTerminated int    `json:"sessions_terminated"`
}

type UseCase = ucdef.UserAction[*Input, *Output]

type usecase struct {
	dc *domain.Container
}

func New(dc *domain.Container) UseCase {
	return &usecase{dc: dc}
}

func (uc *usecase) OperationID() string { return OperationID }

func (uc *usecase) Execute(ctx context.Context, input *Input) (*Output, error) {
	admin, err := uc.dc.AdminRepo().FirstOrNil(ctx, user.AdminFilter{ID: &input.ID})
	if err != nil {
		return nil, errx.Wrap(err)
	}
	if admin == nil {
		return nil, errx.New("admin not found", errx.WithCode(CodeAdminNotFound))
	}

	if !admin.IsActive {
		return nil, errx.New("admin is already disabled", errx.WithCode(CodeAdminAlreadyDisabled))
	}

	if admin.IsSuperadmin {
		isSuper := true
		isActive := true
		count, err := uc.dc.AdminRepo().Count(ctx, user.AdminFilter{IsSuperadmin: &isSuper, IsActive: &isActive})
		if err != nil {
			return nil, errx.Wrap(err)
		}
		if count <= 1 {
			return nil, errx.New(
				"cannot disable the last active superadmin",
				errx.WithCode(CodeCannotDisableLastSuperadmin),
			)
		}
	}

	admin.IsActive = false
	_, err = uc.dc.AdminRepo().Update(ctx, admin)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	actorType := string(rbac.ActorTypeAdmin)
	sessions, err := uc.dc.SessionRepo().List(ctx, session.SessionFilter{
		ActorType: &actorType,
		ActorID:   &admin.ID,
	})
	if err != nil {
		return nil, errx.Wrap(err)
	}

	terminatedCount := 0
	if len(sessions) > 0 {
		err = uc.dc.SessionRepo().BulkDelete(ctx, sessions)
		if err != nil {
			return nil, errx.Wrap(err)
		}
		terminatedCount = len(sessions)
	}

	return &Output{
		ID:                 admin.ID,
		Username:           admin.Username,
		IsActive:           false,
		SessionsTerminated: terminatedCount,
	}, nil
}
