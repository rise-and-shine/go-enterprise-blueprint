package create_superadmin

import (
	"context"

	"go-enterprise-blueprint/internal/modules/auth/domain"
	"go-enterprise-blueprint/internal/modules/auth/domain/user"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/hasher"
	"github.com/rise-and-shine/pkg/ucdef"
)

type Input struct {
	Username string
	Password string
}

type UseCase = ucdef.ManualCommand[Input]

type usecase struct {
	domainContainer *domain.Container
}

func New(domainContainer *domain.Container) UseCase {
	return &usecase{
		domainContainer,
	}
}

func (uc *usecase) OperationID() string { return "create-superadmin" }

func (uc *usecase) Execute(ctx context.Context, input Input) error {
	passwordHash, err := hasher.Hash(input.Password)
	if err != nil {
		return errx.Wrap(err)
	}

	_, err = uc.domainContainer.AdminRepo().Create(ctx, &user.Admin{
		Username:     input.Username,
		PasswordHash: passwordHash,
		IsSuperadmin: true,
		IsActive:     true,
	})
	return errx.Wrap(err)
}
