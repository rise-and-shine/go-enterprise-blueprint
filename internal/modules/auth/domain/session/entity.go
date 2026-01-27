package session

import (
	"time"

	"github.com/rise-and-shine/pkg/pg"
)

const (
	CodeSessionNotFound = "SESSION_NOT_FOUND"
)

type Session struct {
	pg.BaseModel

	ID int64 `json:"id" bun:"id,pk,autoincrement"`

	ActorType string `json:"actor_type"`
	ActorID   string `json:"actor_id"`

	AccessToken           string    `json:"-"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshToken          string    `json:"-"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`

	IPAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`

	LastUsedAt time.Time `json:"last_used_at"`
}
