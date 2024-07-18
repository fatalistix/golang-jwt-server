package logout

import (
	"context"
	"fmt"

	"github.com/fatalistix/golang-jwt-server/internal/domain/model"
	"github.com/google/uuid"
)

type RefreshTokenMapper interface {
	FromSignedString(string) (model.RefreshToken, error)
}

type RefreshTokenProvider interface {
	Delete(context.Context, uuid.UUID) error
}

type Usecase struct {
	refreshTokenMapper   RefreshTokenMapper
	refreshTokenProvider RefreshTokenProvider
}

func NewUsecase(rtm RefreshTokenMapper, rtp RefreshTokenProvider) *Usecase {
	return &Usecase{
		refreshTokenMapper:   rtm,
		refreshTokenProvider: rtp,
	}
}

func (u *Usecase) Handle(
	ctx context.Context,
	tokenStr string,
) error {
	const op = "usecase.logout.Usecase.Handle"

	token, err := u.refreshTokenMapper.FromSignedString(tokenStr)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = u.refreshTokenProvider.Delete(ctx, token.Id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
