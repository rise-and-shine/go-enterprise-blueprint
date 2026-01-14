package admin_refresh_token

import (
	"context"
	"time"

	"go-enterprise-blueprint/internal/modules/auth/domain"
	"go-enterprise-blueprint/internal/modules/auth/domain/session"
	"go-enterprise-blueprint/internal/modules/auth/domain/user"
	"go-enterprise-blueprint/internal/modules/auth/service"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/ucdef"
)

const (
	OperationID = "admin-refresh-token"

	CodeInvalidRefreshToken = "INVALID_REFRESH_TOKEN"
	CodeRefreshTokenExpired = "REFRESH_TOKEN_EXPIRED"
	CodeAdminDisabled       = "ADMIN_DISABLED"
)

type Input struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type Output struct {
	AccessToken           string    `json:"access_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
}

type UseCase = ucdef.UserAction[*Input, *Output]

type usecase struct {
	dc *domain.Container
	sc *service.Container
}

func New(dc *domain.Container, sc *service.Container) UseCase {
	return &usecase{dc: dc, sc: sc}
}

func (uc *usecase) OperationID() string { return OperationID }

func (uc *usecase) Execute(ctx context.Context, input *Input) (*Output, error) {
	sess, err := uc.dc.SessionRepo().FirstOrNil(ctx, session.SessionFilter{RefreshToken: &input.RefreshToken})
	if err != nil {
		return nil, errx.Wrap(err)
	}
	if sess == nil {
		return nil, errx.New("refresh token not found", errx.WithCode(CodeInvalidRefreshToken))
	}

	if time.Now().After(sess.RefreshTokenExpiresAt) {
		return nil, errx.New("refresh token has expired", errx.WithCode(CodeRefreshTokenExpired))
	}

	admin, err := uc.dc.AdminRepo().FirstOrNil(ctx, user.AdminFilter{ID: &sess.ActorID})
	if err != nil {
		return nil, errx.Wrap(err)
	}
	if admin == nil || !admin.IsActive {
		return nil, errx.New("admin account is disabled", errx.WithCode(CodeAdminDisabled))
	}

	accessToken, accessExpiresAt, err := uc.sc.JWT().GenerateAccessToken(
		sess.ID,
		sess.ActorType,
		sess.ActorID,
	)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	refreshToken, refreshExpiresAt, err := uc.sc.JWT().GenerateRefreshToken()
	if err != nil {
		return nil, errx.Wrap(err)
	}

	sess.AccessToken = accessToken
	sess.AccessTokenExpiresAt = accessExpiresAt
	sess.RefreshToken = refreshToken
	sess.RefreshTokenExpiresAt = refreshExpiresAt
	sess.LastUsedAt = time.Now()

	_, err = uc.dc.SessionRepo().Update(ctx, sess)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	return &Output{
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessExpiresAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshExpiresAt,
	}, nil
}
