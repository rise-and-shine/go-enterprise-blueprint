package get_actor_permissions

import (
	"context"

	"go-enterprise-blueprint/internal/modules/auth/domain"
	"go-enterprise-blueprint/internal/modules/auth/domain/rbac"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/ucdef"
)

const (
	OperationID = "get-actor-permissions"

	CodeInvalidActorType = "INVALID_ACTOR_TYPE"
)

type Input struct {
	ActorType string `query:"actor_type" validate:"required,oneof=user admin service_acc"`
	ActorID   string `query:"actor_id"   validate:"required,uuid"`
}

type Output struct {
	ActorType   string            `json:"actor_type"`
	ActorID     string            `json:"actor_id"`
	Permissions PermissionsOutput `json:"permissions"`
}

type PermissionsOutput struct {
	FromRoles []string `json:"from_roles"`
	Direct    []string `json:"direct"`
	Effective []string `json:"effective"`
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

	directPerms, err := uc.dc.ActorPermissionRepo().List(ctx, rbac.ActorPermissionFilter{
		ActorType: &actorType,
		ActorID:   &input.ActorID,
	})
	if err != nil {
		return nil, errx.Wrap(err)
	}

	actorRoles, err := uc.dc.ActorRoleRepo().List(ctx, rbac.ActorRoleFilter{
		ActorType: &actorType,
		ActorID:   &input.ActorID,
	})
	if err != nil {
		return nil, errx.Wrap(err)
	}

	var rolePerms []string
	if len(actorRoles) > 0 {
		roleIDs := make([]int64, len(actorRoles))
		for i, ar := range actorRoles {
			roleIDs[i] = ar.RoleID
		}

		for _, roleID := range roleIDs {
			perms, err := uc.dc.RolePermissionRepo().List(ctx, rbac.RolePermissionFilter{RoleID: &roleID})
			if err != nil {
				return nil, errx.Wrap(err)
			}
			for _, p := range perms {
				rolePerms = append(rolePerms, p.Permission)
			}
		}
	}

	direct := make([]string, len(directPerms))
	for i, p := range directPerms {
		direct[i] = p.Permission
	}

	effectiveMap := make(map[string]struct{})
	for _, p := range rolePerms {
		effectiveMap[p] = struct{}{}
	}
	for _, p := range direct {
		effectiveMap[p] = struct{}{}
	}

	effective := make([]string, 0, len(effectiveMap))
	for p := range effectiveMap {
		effective = append(effective, p)
	}

	return &Output{
		ActorType: input.ActorType,
		ActorID:   input.ActorID,
		Permissions: PermissionsOutput{
			FromRoles: rolePerms,
			Direct:    direct,
			Effective: effective,
		},
	}, nil
}
