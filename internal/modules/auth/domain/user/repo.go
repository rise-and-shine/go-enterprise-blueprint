package user

import "github.com/rise-and-shine/pkg/repogen"

type AdminFilter struct{}

type AdminRepo interface {
	repogen.Repo[Admin, AdminFilter]
}
