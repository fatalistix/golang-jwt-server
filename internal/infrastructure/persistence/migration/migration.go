package migration

import (
	"fmt"

	"github.com/fatalistix/golang-jwt-server/internal/infrastructure/persistence/entity"
	"gorm.io/gorm"
)

func Apply(db *gorm.DB) error {
	const op = "infrastructure.persistence.migration.Apply"

	err := db.AutoMigrate(&entity.User{}, &entity.Role{}, &entity.RefreshToken{})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
