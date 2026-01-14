package update_role

import (
	"context"
	"time"

	"go-enterprise-blueprint/internal/modules/auth/domain"
	"go-enterprise-blueprint/internal/modules/auth/domain/rbac"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/ucdef"
)

const (
	OperationID = "update-role"

	CodeRoleNotFound   = "ROLE_NOT_FOUND"
	CodeRoleNameExists = "ROLE_NAME_EXISTS"
)

type Input struct {
	ID   int64  `json:"id"   validate:"required"`
	Name string `json:"name" validate:"required,min=2,max=100"`
}

type Output struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	UpdatedAt time.Time `json:"updated_at"`
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
	role, err := uc.dc.RoleRepo().FirstOrNil(ctx, rbac.RoleFilter{ID: &input.ID})
	if err != nil {
		return nil, errx.Wrap(err)
	}
	if role == nil {
		return nil, errx.New("role not found", errx.WithCode(CodeRoleNotFound))
	}

	if role.Name != input.Name {
		existing, err := uc.dc.RoleRepo().FirstOrNil(ctx, rbac.RoleFilter{Name: &input.Name})
		if err != nil {
			return nil, errx.Wrap(err)
		}
		if existing != nil && existing.ID != role.ID {
			return nil, errx.New("another role with this name already exists", errx.WithCode(CodeRoleNameExists))
		}
	}

	role.Name = input.Name
	updated, err := uc.dc.RoleRepo().Update(ctx, role)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	return &Output{
		ID:        updated.ID,
		Name:      updated.Name,
		UpdatedAt: updated.UpdatedAt,
	}, nil
}
