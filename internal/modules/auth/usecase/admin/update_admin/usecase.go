package update_admin

import (
	"context"
	"time"

	"go-enterprise-blueprint/internal/modules/auth/domain"
	"go-enterprise-blueprint/internal/modules/auth/domain/user"
	"go-enterprise-blueprint/internal/modules/auth/service"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/ucdef"
)

const (
	OperationID = "update-admin"

	CodeAdminNotFound              = "ADMIN_NOT_FOUND"
	CodeUsernameExists             = "USERNAME_EXISTS"
	CodeCannotDemoteLastSuperadmin = "CANNOT_DEMOTE_LAST_SUPERADMIN"
)

type Input struct {
	ID           string  `json:"id"                      validate:"required,uuid"`
	Username     *string `json:"username,omitempty"      validate:"omitempty,min=3,max=50"`
	Password     *string `json:"password,omitempty"      validate:"omitempty,min=8"`
	IsSuperadmin *bool   `json:"is_superadmin,omitempty"`
}

type Output struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	IsSuperadmin bool      `json:"is_superadmin"`
	IsActive     bool      `json:"is_active"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UseCase = ucdef.UserAction[*Input, *Output]

type usecase struct {
	dc *domain.Container
	sc *service.Container
}

func New(dc *domain.Container, sc *service.Container) UseCase {
	return &usecase{dc: dc, sc: sc}
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

	if input.Username != nil && *input.Username != admin.Username {
		existing, err := uc.dc.AdminRepo().FirstOrNil(ctx, user.AdminFilter{Username: input.Username})
		if err != nil {
			return nil, errx.Wrap(err)
		}
		if existing != nil && existing.ID != admin.ID {
			return nil, errx.New("another admin with this username already exists", errx.WithCode(CodeUsernameExists))
		}
		admin.Username = *input.Username
	}

	if input.Password != nil {
		passwordHash, err := uc.sc.Hasher().Hash(*input.Password)
		if err != nil {
			return nil, errx.Wrap(err)
		}
		admin.PasswordHash = passwordHash
	}

	if input.IsSuperadmin != nil && *input.IsSuperadmin != admin.IsSuperadmin {
		if admin.IsSuperadmin && !*input.IsSuperadmin {
			isSuper := true
			isActive := true
			count, err := uc.dc.AdminRepo().Count(ctx, user.AdminFilter{IsSuperadmin: &isSuper, IsActive: &isActive})
			if err != nil {
				return nil, errx.Wrap(err)
			}
			if count <= 1 {
				return nil, errx.New("cannot demote the last superadmin", errx.WithCode(CodeCannotDemoteLastSuperadmin))
			}
		}
		admin.IsSuperadmin = *input.IsSuperadmin
	}

	updated, err := uc.dc.AdminRepo().Update(ctx, admin)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	return &Output{
		ID:           updated.ID,
		Username:     updated.Username,
		IsSuperadmin: updated.IsSuperadmin,
		IsActive:     updated.IsActive,
		UpdatedAt:    updated.UpdatedAt,
	}, nil
}
