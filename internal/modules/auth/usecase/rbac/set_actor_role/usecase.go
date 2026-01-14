package set_actor_role

import (
	"context"

	"go-enterprise-blueprint/internal/modules/auth/domain"
	"go-enterprise-blueprint/internal/modules/auth/domain/rbac"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/ucdef"
)

const (
	OperationID = "set-actor-role"

	CodeInvalidActorType = "INVALID_ACTOR_TYPE"
	CodeRoleNotFound     = "ROLE_NOT_FOUND"
)

type Input struct {
	ActorType string  `json:"actor_type" validate:"required,oneof=user admin service_acc"`
	ActorID   string  `json:"actor_id"   validate:"required,uuid"`
	RoleIDs   []int64 `json:"role_ids"   validate:"required"`
}

type Output struct {
	ActorType string       `json:"actor_type"`
	ActorID   string       `json:"actor_id"`
	Roles     []RoleOutput `json:"roles"`
}

type RoleOutput struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
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
	actorType := rbac.ActorType(input.ActorType)
	if !actorType.IsValid() {
		return nil, errx.New("invalid actor type", errx.WithCode(CodeInvalidActorType))
	}

	if len(input.RoleIDs) > 0 {
		roles, err := uc.dc.RoleRepo().List(ctx, rbac.RoleFilter{IDs: input.RoleIDs})
		if err != nil {
			return nil, errx.Wrap(err)
		}
		if len(roles) != len(input.RoleIDs) {
			return nil, errx.New("one or more roles not found", errx.WithCode(CodeRoleNotFound))
		}
	}

	existing, err := uc.dc.ActorRoleRepo().List(ctx, rbac.ActorRoleFilter{
		ActorType: &actorType,
		ActorID:   &input.ActorID,
	})
	if err != nil {
		return nil, errx.Wrap(err)
	}

	if len(existing) > 0 {
		err = uc.dc.ActorRoleRepo().BulkDelete(ctx, existing)
		if err != nil {
			return nil, errx.Wrap(err)
		}
	}

	if len(input.RoleIDs) > 0 {
		newActorRoles := make([]rbac.ActorRole, len(input.RoleIDs))
		for i, roleID := range input.RoleIDs {
			newActorRoles[i] = rbac.ActorRole{
				ActorType: actorType,
				ActorID:   input.ActorID,
				RoleID:    roleID,
			}
		}

		err = uc.dc.ActorRoleRepo().BulkCreate(ctx, newActorRoles)
		if err != nil {
			return nil, errx.Wrap(err)
		}
	}

	roles, err := uc.dc.RoleRepo().List(ctx, rbac.RoleFilter{IDs: input.RoleIDs})
	if err != nil {
		return nil, errx.Wrap(err)
	}

	roleOutputs := make([]RoleOutput, len(roles))
	for i, r := range roles {
		roleOutputs[i] = RoleOutput{ID: r.ID, Name: r.Name}
	}

	return &Output{
		ActorType: input.ActorType,
		ActorID:   input.ActorID,
		Roles:     roleOutputs,
	}, nil
}
