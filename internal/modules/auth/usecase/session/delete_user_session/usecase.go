package delete_user_session

import (
	"context"

	"go-enterprise-blueprint/internal/modules/auth/domain"
	"go-enterprise-blueprint/internal/modules/auth/domain/session"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/ucdef"
)

const (
	OperationID = "delete-user-session"

	CodeSessionNotFound  = "SESSION_NOT_FOUND"
	CodePermissionDenied = "PERMISSION_DENIED"
)

type Input struct {
	SessionID    int64  `json:"session_id" validate:"required"`
	ActorType    string `json:"-"`
	ActorID      string `json:"-"`
	IsSuperadmin bool   `json:"-"`
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
	sess, err := uc.dc.SessionRepo().FirstOrNil(ctx, session.SessionFilter{ID: &input.SessionID})
	if err != nil {
		return nil, errx.Wrap(err)
	}
	if sess == nil {
		return nil, errx.New("session not found", errx.WithCode(CodeSessionNotFound))
	}

	if !input.IsSuperadmin {
		if sess.ActorType != input.ActorType || sess.ActorID != input.ActorID {
			return nil, errx.New("not authorized to delete this session", errx.WithCode(CodePermissionDenied))
		}
	}

	err = uc.dc.SessionRepo().Delete(ctx, sess)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	return &Output{Success: true}, nil
}
