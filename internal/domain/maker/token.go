package maker

import (
	"fmt"
	"time"

	"github.com/fatalistix/golang-jwt-server/internal/config"
	"github.com/fatalistix/golang-jwt-server/internal/domain/model"
	"github.com/google/uuid"
)

type TokenMaker struct {
	accessTokenTtl  time.Duration
	refreshTokenTtl time.Duration
	issuer          string
}

func MakeTokenMaker(config config.TokenConfig) TokenMaker {
	return TokenMaker{
		accessTokenTtl:  config.AccessTokenTtl,
		refreshTokenTtl: config.RefreshTokenTtl,
		issuer:          config.Issuer,
	}
}

func (m TokenMaker) MakeAccessToken(userId uint, roles []model.Role) (model.AccessToken, error) {
	const op = "domain.maker.TokenMaker.MakeAccessToken"

	id, err := uuid.NewRandom()
	if err != nil {
		return model.AccessToken{}, fmt.Errorf("%s: %w", op, err)
	}

	rolesStr := make([]string, 0, len(roles))
	for _, v := range roles {
		rolesStr = append(rolesStr, v.String())
	}

	return model.AccessToken{
		Issuer:    m.issuer,
		ExpiresAt: time.Now().Add(m.accessTokenTtl),
		NotBefore: time.Now(),
		IssuedAt:  time.Now(),
		Id:        id,
		UserId:    userId,
		Roles:     rolesStr,
	}, nil
}

func (m TokenMaker) MakeAccessTokenFromUser(user model.User) (model.AccessToken, error) {
	const op = "domain.maker.TokenMaker.MakeAccessTokenFromUser"

	t, err := m.MakeAccessToken(user.Id, user.Roles)
	if err != nil {
		return t, fmt.Errorf("%s: %w", op, err)
	}

	return t, nil
}

func (m TokenMaker) MakeRefreshToken(userId uint, accessTokenId uuid.UUID) (model.RefreshToken, error) {
	const op = "domain.maker.TokenMaker.MakeRefreshToken"

	id, err := uuid.NewRandom()
	if err != nil {
		return model.RefreshToken{}, fmt.Errorf("%s: %w", op, err)
	}

	return model.RefreshToken{
		Issuer:        m.issuer,
		ExpiresAt:     time.Now().Add(m.refreshTokenTtl),
		NotBefore:     time.Now(),
		IssuedAt:      time.Now(),
		Id:            id,
		UserId:        userId,
		AccessTokenId: accessTokenId,
	}, nil
}
