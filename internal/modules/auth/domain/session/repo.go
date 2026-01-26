package session

import "github.com/rise-and-shine/pkg/repogen"

type Filter struct {
	ID           *int64
	ActorType    *string
	ActorID      *string
	AccessToken  *string
	RefreshToken *string

	Limit  int
	Offset int
}

type Repo interface {
	repogen.Repo[Session, Filter]
}
