package model

import (
	"time"

	"github.com/google/uuid"
)

type AccessToken struct {
	Issuer    string
	ExpiresAt time.Time
	NotBefore time.Time
	IssuedAt  time.Time
	Id        uuid.UUID
	UserId    uint
	Roles     []string
}
