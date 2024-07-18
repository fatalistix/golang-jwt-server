package repository

import (
	"context"
	"fmt"
	"time"

	trmgorm "github.com/avito-tech/go-transaction-manager/drivers/gorm/v2"
	"github.com/fatalistix/golang-jwt-server/internal/domain/model"
	"github.com/fatalistix/golang-jwt-server/internal/infrastructure/persistence/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshTokenMapper interface {
	ToEntity(model.RefreshToken) entity.RefreshToken
	FromEntity(entity.RefreshToken) model.RefreshToken
}

type RefreshTokenRepository struct {
	db     *gorm.DB
	getter *trmgorm.CtxGetter
	mapper RefreshTokenMapper
}

func NewRefreshTokenRepository(db *gorm.DB, c *trmgorm.CtxGetter, m RefreshTokenMapper) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		db:     db,
		getter: c,
		mapper: m,
	}
}

func (r *RefreshTokenRepository) Save(ctx context.Context, refreshToken model.RefreshToken) (uuid.UUID, error) {
	const op = "infrastructure.persistence.repository.RefreshTokenRepository.Save"

	entityRefreshToken := r.mapper.ToEntity(refreshToken)
	result := r.getter.DefaultTrOrDB(ctx, r.db).Create(&entityRefreshToken)
	if result.Error != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, result.Error)
	}

	return entityRefreshToken.Id, nil
}

func (r *RefreshTokenRepository) Update(ctx context.Context, refreshToken model.RefreshToken) error {
	const op = "infrastructure.persistence.repository.RefreshTokenRepository.Update"

	entityRefreshToken := r.mapper.ToEntity(refreshToken)
	result := r.getter.DefaultTrOrDB(ctx, r.db).Updates(&entityRefreshToken)
	if result.Error != nil {
		return fmt.Errorf("%s: %w", op, result.Error)
	}

	return nil
}

func (r *RefreshTokenRepository) Delete(ctx context.Context, refreshTokenId uuid.UUID) error {
	const op = "infrastructure.persistence.repository.RefreshTokenRepository.Delete"

	result := r.getter.DefaultTrOrDB(ctx, r.db).Where("id = ?", refreshTokenId).Delete(&[]entity.RefreshToken{})
	if result.Error != nil {
		return fmt.Errorf("%s: %w", op, result.Error)
	}

	return nil
}

func (r *RefreshTokenRepository) RefreshToken(ctx context.Context, id uuid.UUID) (model.RefreshToken, error) {
	const op = "infrastructure.persistence.repository.RefreshTokenRepository.Delete"

	entityRefreshToken := entity.RefreshToken{}
	result := r.getter.DefaultTrOrDB(ctx, r.db).First(&entityRefreshToken, "id = ?", id)
	if result.Error != nil {
		return model.RefreshToken{}, fmt.Errorf("%s: %w", op, result.Error)
	}

	modelRefreshToken := r.mapper.FromEntity(entityRefreshToken)
	return modelRefreshToken, nil
}

func (r *RefreshTokenRepository) DeleteExpired() error {
	const op = "infrastructure.persistence.repository.RefreshTokenRepository.DeleteExpired"

	result := r.db.Where("expires_at < ?", time.Now()).Delete(&[]entity.RefreshToken{})
	if result.Error != nil {
		return fmt.Errorf("%s: %w", op, result.Error)
	}

	return nil
}
