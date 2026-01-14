package set_actor_permission

import (
	"context"

	"go-enterprise-blueprint/internal/modules/auth/domain"
	"go-enterprise-blueprint/internal/modules/auth/domain/rbac"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/ucdef"
)

const (
	OperationID = "set-actor-permission"

	CodeInvalidActorType = "INVALID_ACTOR_TYPE"
)

type Input struct {
	ActorType   string   `json:"actor_type"  validate:"required,oneof=user admin service_acc"`
	ActorID     string   `json:"actor_id"    validate:"required,uuid"`
	Permissions []string `json:"permissions" validate:"required"`
}

type Output struct {
	ActorType   string   `json:"actor_type"`
	ActorID     string   `json:"actor_id"`
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
	actorType := rbac.ActorType(input.ActorType)
	if !actorType.IsValid() {
		return nil, errx.New("invalid actor type", errx.WithCode(CodeInvalidActorType))
	}

	existing, err := uc.dc.ActorPermissionRepo().List(ctx, rbac.ActorPermissionFilter{
		ActorType: &actorType,
		ActorID:   &input.ActorID,
	})
	if err != nil {
		return nil, errx.Wrap(err)
	}

	if len(existing) > 0 {
		err = uc.dc.ActorPermissionRepo().BulkDelete(ctx, existing)
		if err != nil {
			return nil, errx.Wrap(err)
		}
	}

	if len(input.Permissions) > 0 {
		newPerms := make([]rbac.ActorPermission, len(input.Permissions))
		for i, p := range input.Permissions {
			newPerms[i] = rbac.ActorPermission{
				ActorType:  actorType,
				ActorID:    input.ActorID,
				Permission: p,
			}
		}

		err = uc.dc.ActorPermissionRepo().BulkCreate(ctx, newPerms)
		if err != nil {
			return nil, errx.Wrap(err)
		}
	}

	return &Output{
		ActorType:   input.ActorType,
		ActorID:     input.ActorID,
		Permissions: input.Permissions,
	}, nil
}
