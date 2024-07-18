package register

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/fatalistix/golang-jwt-server/internal/domain/model"
)

type UserSaver interface {
	Save(ctx context.Context, user model.User) (uint, error)
}

type Usecase struct {
	userSaver   UserSaver
	encryptCost int
}

func NewUsecase(
	s UserSaver,
	encryptCost int,
) *Usecase {
	return &Usecase{
		userSaver:   s,
		encryptCost: encryptCost,
	}
}

func (u *Usecase) Handle(
	ctx context.Context,
	username, password string,
) (uint, error) {
	const op = "usecase.register.Usecase.Handle"

	encodedPassword, err := bcrypt.GenerateFromPassword([]byte(password), u.encryptCost)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	user := model.User{
		Username: username,
		Password: string(encodedPassword),
		Roles:    []model.Role{model.RoleUser},
	}

	id, err := u.userSaver.Save(ctx, user)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
