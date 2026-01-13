package create_superadmin

import (
	"context"

	"github.com/rise-and-shine/pkg/ucdef"
)

const (
	OperationID = "create-superadmin"
)

type Request struct{}

type Response struct{}

type UseCase = ucdef.UserAction[*Request, *Response]

type usecase struct{}

func (uc *usecase) OperationID() string { return OperationID }

func (uc *usecase) Execute(ctx context.Context) (*Response, error) {
	// TODO: ...
	return nil, nil
}
