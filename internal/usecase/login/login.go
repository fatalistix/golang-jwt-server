package login

import (
	"context"
	"fmt"

	"github.com/fatalistix/golang-jwt-server/internal/domain/maker"
	"github.com/fatalistix/golang-jwt-server/internal/domain/model"
	"github.com/fatalistix/golang-jwt-server/internal/dto"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserProvider interface {
	UserByUsername(ctx context.Context, username string) (model.User, error)
}

type RefreshTokenSaver interface {
	Save(ctx context.Context, refreshToken model.RefreshToken) (uuid.UUID, error)
}

type AccessTokenMapper interface {
	ToSignedString(t model.AccessToken) (string, error)
}

type RefreshTokenMapper interface {
	ToSignedString(t model.RefreshToken) (string, error)
}

type Usecase struct {
	userProvider       UserProvider
	refreshTokenSaver  RefreshTokenSaver
	tokenMaker         maker.TokenMaker
	accessTokenMapper  AccessTokenMapper
	refreshTokenMapper RefreshTokenMapper
}

func NewUsecase(
	p UserProvider,
	s RefreshTokenSaver,
	m maker.TokenMaker,
	atm AccessTokenMapper,
	rtm RefreshTokenMapper,
) *Usecase {
	return &Usecase{
		userProvider:       p,
		refreshTokenSaver:  s,
		tokenMaker:         m,
		accessTokenMapper:  atm,
		refreshTokenMapper: rtm,
	}
}

func (u *Usecase) Handle(ctx context.Context, username, password string) (dto.Tokens, error) {
	const op = "usecase.login.Usecase.Handle"

	user, err := u.userProvider.UserByUsername(ctx, username)
	if err != nil {
		return dto.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return dto.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	accessToken, err := u.tokenMaker.MakeAccessTokenFromUser(user)
	if err != nil {
		return dto.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	refreshToken, err := u.tokenMaker.MakeRefreshToken(user.Id, accessToken.Id)
	if err != nil {
		return dto.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	accessTokenStr, err := u.accessTokenMapper.ToSignedString(accessToken)
	if err != nil {
		return dto.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	refreshTokenStr, err := u.refreshTokenMapper.ToSignedString(refreshToken)
	if err != nil {
		return dto.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	_, err = u.refreshTokenSaver.Save(ctx, refreshToken)
	if err != nil {
		return dto.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	return dto.Tokens{
		AccessToken:  accessTokenStr,
		RefreshToken: refreshTokenStr,
	}, nil
}
