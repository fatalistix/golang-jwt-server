package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshToken struct {
	Id            uuid.UUID `gorm:"primaryKey"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	Issuer        string
	ExpiresAt     time.Time
	NotBefore     time.Time
	IssuedAt      time.Time
	User          User `gorm:"foreignKey:UserId;constraint:OnUpdate:NO ACTION;OnDelete:NO ACTION"`
	UserId        uint
	AccessTokenId uuid.UUID
}
