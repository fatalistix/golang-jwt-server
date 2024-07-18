package repository

import (
	"context"
	"fmt"

	trmgorm "github.com/avito-tech/go-transaction-manager/drivers/gorm/v2"

	"github.com/fatalistix/golang-jwt-server/internal/domain/model"
	"github.com/fatalistix/golang-jwt-server/internal/infrastructure/persistence/entity"
	"gorm.io/gorm"
)

type UserMapper interface {
	ToEntity(model.User) entity.User
	FromEntityToDomain(entity.User) (model.User, error)
}

type UserRepository struct {
	db     *gorm.DB
	getter *trmgorm.CtxGetter
	mapper UserMapper
}

func NewUserRepository(db *gorm.DB, c *trmgorm.CtxGetter, m UserMapper) *UserRepository {
	return &UserRepository{
		db:     db,
		getter: c,
		mapper: m,
	}
}

func (r *UserRepository) Save(ctx context.Context, user model.User) (uint, error) {
	const op = "infrastructure.persistence.repository.UserRepository.Save"

	entityUser := r.mapper.ToEntity(user)
	result := r.getter.DefaultTrOrDB(ctx, r.db).Create(&entityUser)
	if result.Error != nil {
		return 0, fmt.Errorf("%s: %w", op, result.Error)
	}

	return entityUser.ID, nil
}

func (r *UserRepository) UserByUsername(ctx context.Context, username string) (model.User, error) {
	const op = "infrastructure.persistence.repository.UserRepository.UserByUsername"

	entityUser := entity.User{}
	result := r.getter.DefaultTrOrDB(ctx, r.db).First(&entityUser, "username = ?", username)
	if result.Error != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, result.Error)
	}

	modelUser, err := r.mapper.FromEntityToDomain(entityUser)
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return modelUser, nil
}

func (r *UserRepository) User(ctx context.Context, id uint) (model.User, error) {
	const op = "infrastructure.persistence.repository.UserRepository.User"

	entityUser := entity.User{}
	result := r.getter.DefaultTrOrDB(ctx, r.db).First(&entityUser, "id = ?", id)
	if result.Error != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, result.Error)
	}

	modelUser, err := r.mapper.FromEntityToDomain(entityUser)
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return modelUser, nil
}
