package uc1001createuser

import (
	"context"
	"go-enterprise-blueprint/internal/modules/auth/domain"
	"go-enterprise-blueprint/internal/modules/auth/domain/user"
	"go-enterprise-blueprint/internal/modules/auth/service"
	"go-enterprise-blueprint/internal/portal"
	"go-enterprise-blueprint/pkg/ucdef"
)

type CreateUserIn struct {
	Body struct {
		Username string
		Password string
	}
}

type CreateUserOut struct {
	Body user.User
}

type UseCase ucdef.UserAction[*CreateUserIn, *CreateUserOut]

func New(dc domain.Container, sc service.Container, pc portal.Container) UseCase {
	return &useCase{dc, sc, pc}
}

type useCase struct {
	dc domain.Container
	sc service.Container
	pc portal.Container
}

func (uc useCase) OperationID() string {
	return "uc-1001-create-user"
}

func (uc useCase) Execute(_ context.Context, in *CreateUserIn) (*CreateUserOut, error) {
	res := &CreateUserOut{
		Body: user.User{
			ID:       "wqer",
			Username: in.Body.Username,
		},
	}

	return res, nil
}
