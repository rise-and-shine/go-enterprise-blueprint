package createsuperadmin

import (
	"context"
	"go-enterprise-blueprint/internal/modules/auth/domain"
	"go-enterprise-blueprint/internal/modules/auth/domain/rbac"
	"go-enterprise-blueprint/internal/modules/auth/domain/user"
	"go-enterprise-blueprint/internal/portal/auth"

	"github.com/code19m/errx"
	"github.com/google/uuid"
	"github.com/rise-and-shine/pkg/hasher"
	"github.com/rise-and-shine/pkg/ucdef"
)

type Input struct {
	Username string
	Password string
}

type UseCase = ucdef.ManualCommand[*Input]

type usecase struct {
	domainContainer *domain.Container
}

func New(domainContainer *domain.Container) UseCase {
	return &usecase{
		domainContainer,
	}
}

func (uc *usecase) OperationID() string { return "create-superadmin" }

func (uc *usecase) Execute(ctx context.Context, input *Input) error {
	// Hash the password
	passwordHash, err := hasher.Hash(input.Password)
	if err != nil {
		return errx.Wrap(err)
	}

	// Start UOW
	uow, err := uc.domainContainer.UOWFactory().NewUOW(ctx)
	if err != nil {
		return errx.Wrap(err)
	}
	defer uow.DiscardUnapplied()

	// Create admin
	a, err := uow.Admin().Create(ctx, &user.Admin{
		ID:           uuid.NewString(),
		Username:     input.Username,
		PasswordHash: passwordHash,
		IsActive:     true,
	})
	if err != nil {
		return errx.Wrap(err)
	}

	// Create actor permission with superadmin permission
	_, err = uow.ActorPermission().Create(ctx, &rbac.ActorPermission{
		ActorType:  rbac.ActorTypeAdmin,
		ActorID:    a.ID,
		Permission: auth.PermissionSuperadmin,
	})
	if err != nil {
		return errx.Wrap(err)
	}

	// Apply UOW
	err = uow.ApplyChanges()
	return errx.Wrap(err)
}
