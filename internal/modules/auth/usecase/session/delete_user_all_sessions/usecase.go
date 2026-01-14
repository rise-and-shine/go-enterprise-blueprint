package delete_user_all_sessions

import (
	"context"

	"go-enterprise-blueprint/internal/modules/auth/domain"
	"go-enterprise-blueprint/internal/modules/auth/domain/rbac"
	"go-enterprise-blueprint/internal/modules/auth/domain/session"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/ucdef"
)

const (
	OperationID = "delete-user-all-sessions"

	CodeInvalidActorType = "INVALID_ACTOR_TYPE"
)

type Input struct {
	ActorType string `json:"actor_type" validate:"required,oneof=user admin service_acc"`
	ActorID   string `json:"actor_id"   validate:"required,uuid"`
}

type Output struct {
	DeletedCount int `json:"deleted_count"`
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
	if !rbac.ActorType(input.ActorType).IsValid() {
		return nil, errx.New("invalid actor type", errx.WithCode(CodeInvalidActorType))
	}

	sessions, err := uc.dc.SessionRepo().List(ctx, session.SessionFilter{
		ActorType: &input.ActorType,
		ActorID:   &input.ActorID,
	})
	if err != nil {
		return nil, errx.Wrap(err)
	}

	if len(sessions) > 0 {
		err = uc.dc.SessionRepo().BulkDelete(ctx, sessions)
		if err != nil {
			return nil, errx.Wrap(err)
		}
	}

	return &Output{DeletedCount: len(sessions)}, nil
}
