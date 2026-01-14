package get_role_permissions

import (
	"context"

	"go-enterprise-blueprint/internal/modules/auth/domain"
	"go-enterprise-blueprint/internal/modules/auth/domain/rbac"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/ucdef"
)

const (
	OperationID = "get-role-permissions"

	CodeRoleNotFound = "ROLE_NOT_FOUND"
)

type Input struct {
	RoleID int64 `query:"role_id" validate:"required"`
}

type Output struct {
	RoleID      int64    `json:"role_id"`
	RoleName    string   `json:"role_name"`
	Permissions []string `json:"permissions"`
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
	role, err := uc.dc.RoleRepo().FirstOrNil(ctx, rbac.RoleFilter{ID: &input.RoleID})
	if err != nil {
		return nil, errx.Wrap(err)
	}
	if role == nil {
		return nil, errx.New("role not found", errx.WithCode(CodeRoleNotFound))
	}

	perms, err := uc.dc.RolePermissionRepo().List(ctx, rbac.RolePermissionFilter{RoleID: &input.RoleID})
	if err != nil {
		return nil, errx.Wrap(err)
	}

	permissions := make([]string, len(perms))
	for i, p := range perms {
		permissions[i] = p.Permission
	}

	return &Output{
		RoleID:      role.ID,
		RoleName:    role.Name,
		Permissions: permissions,
	}, nil
}
