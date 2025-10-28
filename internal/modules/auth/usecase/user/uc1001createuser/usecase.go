package uc1001createuser

import (
	"context"
	"go-enterprise-blueprint/internal/modules/auth/domain"
	"go-enterprise-blueprint/internal/modules/auth/domain/user"
	"go-enterprise-blueprint/internal/modules/auth/service"
	"go-enterprise-blueprint/internal/portal"
	"go-enterprise-blueprint/pkg/ucdef"
)

type CreateUserRequest struct {
	Body struct {
		Username string
		Password string
	}
}

type CreateUserResponse struct {
	Body user.User
}

type state struct {
	req  CreateUserRequest
	user user.User
}

type UseCase ucdef.UserWriteAction[CreateUserRequest, *CreateUserResponse, *state]

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

func (uc useCase) Validate(_ context.Context, req CreateUserRequest) (*state, error) {
	// validate input
	return &state{req, user.User{}}, nil
}

func (uc useCase) Execute(_ context.Context, s *state) (*CreateUserResponse, error) {
	res := &CreateUserResponse{
		Body: user.User{
			ID:       "wqer",
			Username: s.req.Body.Username,
		},
	}

	return res, nil
}

func (uc useCase) AuditInfo(s *state) ucdef.AuditInfo {
	return ucdef.AuditInfo{
		Tags:          []string{"tender", "tender-main"},
		AggregateID:   &s.user.ID,
		AggregateName: &s.user.Username, // should return table name
	}
}
