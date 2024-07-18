package database

import (
	"fmt"

	"github.com/fatalistix/golang-jwt-server/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgres(cfg config.PostgresConfig) (*gorm.DB, error) {
	const op = "infrastructure.persistence.database.NewPostgres"

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port,
		cfg.User, cfg.Password,
		cfg.Db, cfg.SslMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	sqlDb, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = sqlDb.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return db, nil
}
