package create_role

import (
	"context"
	"time"

	"go-enterprise-blueprint/internal/modules/auth/domain"
	"go-enterprise-blueprint/internal/modules/auth/domain/rbac"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/ucdef"
)

const (
	OperationID = "create-role"

	CodeRoleNameExists = "ROLE_NAME_EXISTS"
)

type Input struct {
	Name string `json:"name" validate:"required,min=2,max=100"`
}

type Output struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
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
	exists, err := uc.dc.RoleRepo().Exists(ctx, rbac.RoleFilter{Name: &input.Name})
	if err != nil {
		return nil, errx.Wrap(err)
	}
	if exists {
		return nil, errx.New("role with this name already exists", errx.WithCode(CodeRoleNameExists))
	}

	role := &rbac.Role{Name: input.Name}
	created, err := uc.dc.RoleRepo().Create(ctx, role)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	return &Output{
		ID:        created.ID,
		Name:      created.Name,
		CreatedAt: created.CreatedAt,
	}, nil
}
