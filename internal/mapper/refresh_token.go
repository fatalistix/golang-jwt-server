package mapper

import (
	"fmt"

	"github.com/fatalistix/golang-jwt-server/internal/domain/model"
	"github.com/fatalistix/golang-jwt-server/internal/infrastructure/persistence/entity"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type RefreshTokenMapper struct {
	signingMethod jwt.SigningMethod
	signingKey    []byte
}

type refreshTokenClaims struct {
	jwt.RegisteredClaims
	AccessTokenId string `json:"access_token_id"`
}

func MakeRefreshTokenMapper(signingKey string) RefreshTokenMapper {
	return RefreshTokenMapper{
		signingMethod: jwt.SigningMethodHS256,
		signingKey:    []byte(signingKey),
	}
}

func (m RefreshTokenMapper) ToEntity(t model.RefreshToken) entity.RefreshToken {
	return entity.RefreshToken{
		Issuer:        t.Issuer,
		ExpiresAt:     t.ExpiresAt,
		NotBefore:     t.NotBefore,
		IssuedAt:      t.IssuedAt,
		Id:            t.Id,
		UserId:        t.UserId,
		AccessTokenId: t.AccessTokenId,
	}
}

func (m RefreshTokenMapper) FromEntity(t entity.RefreshToken) model.RefreshToken {
	return model.RefreshToken{
		Issuer:        t.Issuer,
		ExpiresAt:     t.ExpiresAt,
		NotBefore:     t.NotBefore,
		IssuedAt:      t.IssuedAt,
		Id:            t.Id,
		UserId:        t.UserId,
		AccessTokenId: t.AccessTokenId,
	}
}

func (m RefreshTokenMapper) ToSignedString(t model.RefreshToken) (string, error) {
	const op = "mapper.RefreshTokenMapper.ToSignedString"

	claims := refreshTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    t.Issuer,
			ExpiresAt: jwt.NewNumericDate(t.ExpiresAt),
			NotBefore: jwt.NewNumericDate(t.NotBefore),
			IssuedAt:  jwt.NewNumericDate(t.IssuedAt),
			ID:        t.Id.String(),
			Subject:   fmt.Sprint(t.UserId),
		},
		AccessTokenId: t.AccessTokenId.String(),
	}

	token := jwt.NewWithClaims(m.signingMethod, claims)
	ss, err := token.SignedString(m.signingKey)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return ss, nil
}

func (m RefreshTokenMapper) FromSignedString(ss string) (model.RefreshToken, error) {
	const op = "mapper.RefreshTokenMapper.FromSignedString"

	claims := refreshTokenClaims{}
	_, err := jwt.ParseWithClaims(ss, &claims, func(token *jwt.Token) (interface{}, error) {
		const op = "mapper.RefreshTokenMapper.FromSignedString.ParseWithClaims@lambda"

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%s: unexpected signing method: %v", op, token.Method.Alg())
		}

		return m.signingKey, nil
	})
	if err != nil {
		return model.RefreshToken{}, fmt.Errorf("%s: %w", op, err)
	}

	id, err := uuid.Parse(claims.ID)
	if err != nil {
		return model.RefreshToken{}, fmt.Errorf("%s: %w", op, err)
	}

	var userId uint
	_, err = fmt.Sscan(claims.Subject, &userId)
	if err != nil {
		return model.RefreshToken{}, fmt.Errorf("%s: %w", op, err)
	}

	accessTokenId, err := uuid.Parse(claims.AccessTokenId)
	if err != nil {
		return model.RefreshToken{}, fmt.Errorf("%s: %w", op, err)
	}

	return model.RefreshToken{
		Issuer:        claims.Issuer,
		ExpiresAt:     claims.ExpiresAt.Time,
		NotBefore:     claims.NotBefore.Time,
		IssuedAt:      claims.IssuedAt.Time,
		Id:            id,
		UserId:        userId,
		AccessTokenId: accessTokenId,
	}, nil
}
