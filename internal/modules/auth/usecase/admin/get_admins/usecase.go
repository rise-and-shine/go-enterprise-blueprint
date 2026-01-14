package get_admins

import (
	"context"
	"time"

	"go-enterprise-blueprint/internal/modules/auth/domain"
	"go-enterprise-blueprint/internal/modules/auth/domain/user"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/ucdef"
)

const (
	OperationID = "get-admins"
)

type Input struct {
	Page     int   `query:"page"      validate:"omitempty,min=1"`
	PageSize int   `query:"page_size" validate:"omitempty,min=1,max=100"`
	IsActive *bool `query:"is_active"`
}

type Output struct {
	Items    []AdminOutput `json:"items"`
	Total    int           `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
}

type AdminOutput struct {
	ID           string     `json:"id"`
	Username     string     `json:"username"`
	IsSuperadmin bool       `json:"is_superadmin"`
	IsActive     bool       `json:"is_active"`
	LastActiveAt *time.Time `json:"last_active_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
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

	admins, total, err := uc.dc.AdminRepo().ListWithCount(ctx, user.AdminFilter{
		IsActive: input.IsActive,
		Limit:    input.PageSize,
		Offset:   offset,
	})
	if err != nil {
		return nil, errx.Wrap(err)
	}

	items := make([]AdminOutput, len(admins))
	for i, a := range admins {
		var lastActiveAt *time.Time
		if !a.LastActiveAt.IsZero() {
			lastActiveAt = &a.LastActiveAt
		}

		items[i] = AdminOutput{
			ID:           a.ID,
			Username:     a.Username,
			IsSuperadmin: a.IsSuperadmin,
			IsActive:     a.IsActive,
			LastActiveAt: lastActiveAt,
			CreatedAt:    a.CreatedAt,
		}
	}

	return &Output{
		Items:    items,
		Total:    total,
		Page:     input.Page,
		PageSize: input.PageSize,
	}, nil
}
