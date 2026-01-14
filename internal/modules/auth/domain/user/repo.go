package user

import "github.com/rise-and-shine/pkg/repogen"

type AdminFilter struct {
	ID           *string
	Username     *string
	IsSuperadmin *bool
	IsActive     *bool

	Limit  int
	Offset int
}

type AdminRepo interface {
	repogen.Repo[Admin, AdminFilter]
}
