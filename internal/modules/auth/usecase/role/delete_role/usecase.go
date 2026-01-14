package delete_role

import (
	"context"

	"go-enterprise-blueprint/internal/modules/auth/domain"
	"go-enterprise-blueprint/internal/modules/auth/domain/rbac"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/ucdef"
)

const (
	OperationID = "delete-role"

	CodeRoleNotFound = "ROLE_NOT_FOUND"
)

type Input struct {
	ID int64 `json:"id" validate:"required"`
}

type Output struct {
	Success bool `json:"success"`
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

	err = uc.dc.RoleRepo().Delete(ctx, role)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	return &Output{Success: true}, nil
}
