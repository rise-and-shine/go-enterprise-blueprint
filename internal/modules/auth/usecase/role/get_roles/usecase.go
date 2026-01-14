package get_roles

import (
	"context"
	"time"

	"go-enterprise-blueprint/internal/modules/auth/domain"
	"go-enterprise-blueprint/internal/modules/auth/domain/rbac"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/ucdef"
)

const (
	OperationID = "get-roles"
)

type Input struct {
	Page     int `query:"page"      validate:"omitempty,min=1"`
	PageSize int `query:"page_size" validate:"omitempty,min=1,max=100"`
}

type Output struct {
	Items    []RoleOutput `json:"items"`
	Total    int          `json:"total"`
	Page     int          `json:"page"`
	PageSize int          `json:"page_size"`
}

type RoleOutput struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
	if input.Page == 0 {
		input.Page = 1
	}
	if input.PageSize == 0 {
		input.PageSize = 20
	}

	offset := (input.Page - 1) * input.PageSize

	roles, total, err := uc.dc.RoleRepo().ListWithCount(ctx, rbac.RoleFilter{
		Limit:  input.PageSize,
		Offset: offset,
	})
	if err != nil {
		return nil, errx.Wrap(err)
	}

	items := make([]RoleOutput, len(roles))
	for i, r := range roles {
		items[i] = RoleOutput{
			ID:        r.ID,
			Name:      r.Name,
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
		}
	}

	return &Output{
		Items:    items,
		Total:    total,
		Page:     input.Page,
		PageSize: input.PageSize,
	}, nil
}
