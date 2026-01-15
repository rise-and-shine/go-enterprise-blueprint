package portal

import "go-enterprise-blueprint/internal/portal/auth"

type portal struct{}

func New() auth.Portal {
	return &portal{}
}
