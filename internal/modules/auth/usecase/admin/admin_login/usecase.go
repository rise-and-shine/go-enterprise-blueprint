package admin_login

import (
	"context"
	"time"

	"go-enterprise-blueprint/internal/modules/auth/domain"
	"go-enterprise-blueprint/internal/modules/auth/domain/rbac"
	"go-enterprise-blueprint/internal/modules/auth/domain/session"
	"go-enterprise-blueprint/internal/modules/auth/domain/user"
	"go-enterprise-blueprint/internal/modules/auth/service"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/ucdef"
)

const (
	OperationID = "admin-login"

	CodeInvalidCredentials = "INVALID_CREDENTIALS"
	CodeAdminDisabled      = "ADMIN_DISABLED"
)

type Input struct {
	Username  string `json:"username" validate:"required,min=3,max=50"`
	Password  string `json:"password" validate:"required,min=8"`
	IPAddress string `json:"-"`
	UserAgent string `json:"-"`
}

type Output struct {
	Admin   AdminOutput   `json:"admin"`
	Session SessionOutput `json:"session"`
}

type AdminOutput struct {
	ID           string     `json:"id"`
	Username     string     `json:"username"`
	IsSuperadmin bool       `json:"is_superadmin"`
	IsActive     bool       `json:"is_active"`
	LastActiveAt *time.Time `json:"last_active_at,omitempty"`
}

type SessionOutput struct {
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
	admin, err := uc.dc.AdminRepo().FirstOrNil(ctx, user.AdminFilter{Username: &input.Username})
	if err != nil {
		return nil, errx.Wrap(err)
	}
	if admin == nil {
		return nil, errx.New("invalid username or password", errx.WithCode(CodeInvalidCredentials))
	}

	if !admin.IsActive {
		return nil, errx.New("admin account is disabled", errx.WithCode(CodeAdminDisabled))
	}

	if !uc.sc.Hasher().Compare(input.Password, admin.PasswordHash) {
		return nil, errx.New("invalid username or password", errx.WithCode(CodeInvalidCredentials))
	}

	refreshToken, refreshExpiresAt, err := uc.sc.JWT().GenerateRefreshToken()
	if err != nil {
		return nil, errx.Wrap(err)
	}

	sess := &session.Session{
		ActorType:             string(rbac.ActorTypeAdmin),
		ActorID:               admin.ID,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshExpiresAt,
		IPAddress:             input.IPAddress,
		UserAgent:             input.UserAgent,
		LastUsedAt:            time.Now(),
	}

	createdSession, err := uc.dc.SessionRepo().Create(ctx, sess)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	accessToken, accessExpiresAt, err := uc.sc.JWT().GenerateAccessToken(
		createdSession.ID,
		string(rbac.ActorTypeAdmin),
		admin.ID,
	)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	createdSession.AccessToken = accessToken
	createdSession.AccessTokenExpiresAt = accessExpiresAt
	_, err = uc.dc.SessionRepo().Update(ctx, createdSession)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	admin.LastActiveAt = time.Now()
	_, err = uc.dc.AdminRepo().Update(ctx, admin)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	var lastActiveAt *time.Time
	if !admin.LastActiveAt.IsZero() {
		lastActiveAt = &admin.LastActiveAt
	}

	return &Output{
		Admin: AdminOutput{
			ID:           admin.ID,
			Username:     admin.Username,
			IsSuperadmin: admin.IsSuperadmin,
			IsActive:     admin.IsActive,
			LastActiveAt: lastActiveAt,
		},
		Session: SessionOutput{
			AccessToken:           accessToken,
			AccessTokenExpiresAt:  accessExpiresAt,
			RefreshToken:          refreshToken,
			RefreshTokenExpiresAt: refreshExpiresAt,
		},
	}, nil
}
