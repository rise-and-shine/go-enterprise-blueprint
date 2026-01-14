package set_role_permission

import (
	"context"

	"go-enterprise-blueprint/internal/modules/auth/domain"
	"go-enterprise-blueprint/internal/modules/auth/domain/rbac"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/ucdef"
)

const (
	OperationID = "set-role-permission"

	CodeRoleNotFound = "ROLE_NOT_FOUND"
)

type Input struct {
	RoleID      int64    `json:"role_id"     validate:"required"`
	Permissions []string `json:"permissions" validate:"required"`
}

type Output struct {
	RoleID      int64    `json:"role_id"`
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
	exists, err := uc.dc.RoleRepo().Exists(ctx, rbac.RoleFilter{ID: &input.RoleID})
	if err != nil {
		return nil, errx.Wrap(err)
	}
	if !exists {
		return nil, errx.New("role not found", errx.WithCode(CodeRoleNotFound))
	}

	existing, err := uc.dc.RolePermissionRepo().List(ctx, rbac.RolePermissionFilter{RoleID: &input.RoleID})
	if err != nil {
		return nil, errx.Wrap(err)
	}

	if len(existing) > 0 {
		err = uc.dc.RolePermissionRepo().BulkDelete(ctx, existing)
		if err != nil {
			return nil, errx.Wrap(err)
		}
	}

	if len(input.Permissions) > 0 {
		newPerms := make([]rbac.RolePermission, len(input.Permissions))
		for i, p := range input.Permissions {
			newPerms[i] = rbac.RolePermission{
				RoleID:     input.RoleID,
				Permission: p,
			}
		}

		err = uc.dc.RolePermissionRepo().BulkCreate(ctx, newPerms)
		if err != nil {
			return nil, errx.Wrap(err)
		}
	}

	return &Output{
		RoleID:      input.RoleID,
		Permissions: input.Permissions,
	}, nil
}
