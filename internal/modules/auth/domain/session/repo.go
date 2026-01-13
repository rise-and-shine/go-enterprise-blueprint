package session

import "github.com/rise-and-shine/pkg/repogen"

type SessionFilter struct{}

type SessionRepo struct {
	repogen.Repo[Session, SessionFilter]
}
