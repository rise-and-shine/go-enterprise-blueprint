package create_admin

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
	OperationID = "create-admin"

	CodeUsernameExists = "USERNAME_EXISTS"
)

type Input struct {
	Username     string `json:"username"      validate:"required,min=3,max=50"`
	Password     string `json:"password"      validate:"required,min=8"`
	IsSuperadmin bool   `json:"is_superadmin"`
}

type Output struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	IsSuperadmin bool      `json:"is_superadmin"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
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
	exists, err := uc.dc.AdminRepo().Exists(ctx, user.AdminFilter{Username: &input.Username})
	if err != nil {
		return nil, errx.Wrap(err)
	}
	if exists {
		return nil, errx.New("admin with this username already exists", errx.WithCode(CodeUsernameExists))
	}

	passwordHash, err := uc.sc.Hasher().Hash(input.Password)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	admin := &user.Admin{
		Username:     input.Username,
		PasswordHash: passwordHash,
		IsSuperadmin: input.IsSuperadmin,
		IsActive:     true,
	}

	created, err := uc.dc.AdminRepo().Create(ctx, admin)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	return &Output{
		ID:           created.ID,
		Username:     created.Username,
		IsSuperadmin: created.IsSuperadmin,
		IsActive:     created.IsActive,
		CreatedAt:    created.CreatedAt,
	}, nil
}
