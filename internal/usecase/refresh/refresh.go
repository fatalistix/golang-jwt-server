package refresh

import (
	"context"
	"fmt"
	"time"

	"github.com/fatalistix/golang-jwt-server/internal/domain/maker"
	"github.com/fatalistix/golang-jwt-server/internal/domain/model"
	"github.com/fatalistix/golang-jwt-server/internal/dto"
	"github.com/google/uuid"
)

type AccessTokenMapper interface {
	ToSignedString(model.AccessToken) (string, error)
}

type RefreshTokenMapper interface {
	FromSignedString(string) (model.RefreshToken, error)
	ToSignedString(model.RefreshToken) (string, error)
}

type RefreshTokenProvider interface {
	Delete(context.Context, uuid.UUID) error
	Save(context.Context, model.RefreshToken) (uuid.UUID, error)
	RefreshToken(context.Context, uuid.UUID) (model.RefreshToken, error)
}

type UserProvider interface {
	User(context.Context, uint) (model.User, error)
}

type Usecase struct {
	accessTokenMapper    AccessTokenMapper
	refreshTokenMapper   RefreshTokenMapper
	refreshTokenProvider RefreshTokenProvider
	tokenMaker           maker.TokenMaker
	userProvider         UserProvider
}

func NewUsecase(
	atm AccessTokenMapper,
	rtm RefreshTokenMapper,
	rtp RefreshTokenProvider,
	tm maker.TokenMaker,
	up UserProvider,
) *Usecase {
	return &Usecase{
		accessTokenMapper:    atm,
		refreshTokenMapper:   rtm,
		refreshTokenProvider: rtp,
		tokenMaker:           tm,
		userProvider:         up,
	}
}

func (u *Usecase) Handle(ctx context.Context, refreshTokenStr string) (dto.Tokens, error) {
	const op = "usecase.refresh.Usecase.Handle"

	refreshToken, err := u.refreshTokenMapper.FromSignedString(refreshTokenStr)
	if err != nil {
		return dto.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	_, err = u.refreshTokenProvider.RefreshToken(ctx, refreshToken.Id)
	if err != nil {
		return dto.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	err = u.refreshTokenProvider.Delete(ctx, refreshToken.Id)
	if err != nil {
		return dto.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		return dto.Tokens{}, fmt.Errorf("%s: refresh token expired", op)
	}

	user, err := u.userProvider.User(ctx, refreshToken.UserId)
	if err != nil {
		return dto.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	accessToken, err := u.tokenMaker.MakeAccessTokenFromUser(user)
	if err != nil {
		return dto.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	refreshToken, err = u.tokenMaker.MakeRefreshToken(user.Id, accessToken.Id)
	if err != nil {
		return dto.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	accessTokenStr, err := u.accessTokenMapper.ToSignedString(accessToken)
	if err != nil {
		return dto.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	refreshTokenStr, err = u.refreshTokenMapper.ToSignedString(refreshToken)
	if err != nil {
		return dto.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	_, err = u.refreshTokenProvider.Save(ctx, refreshToken)
	if err != nil {
		return dto.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	return dto.Tokens{
		AccessToken:  accessTokenStr,
		RefreshToken: refreshTokenStr,
	}, nil
}
