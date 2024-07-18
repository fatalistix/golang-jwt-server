package mapper

import (
	"fmt"

	"github.com/fatalistix/golang-jwt-server/internal/domain/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AccessTokenMapper struct {
	signingMethod jwt.SigningMethod
	signingKey    []byte
}

type accessTokenClaims struct {
	jwt.RegisteredClaims
	Roles []string
}

func MakeAccessTokenMapper(signingKey string) AccessTokenMapper {
	return AccessTokenMapper{
		signingMethod: jwt.SigningMethodHS256,
		signingKey:    []byte(signingKey),
	}
}

func (m AccessTokenMapper) ToSignedString(t model.AccessToken) (string, error) {
	const op = "mapper.AccessTokenMapper.ToSignedString"

	claims := accessTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    t.Issuer,
			ExpiresAt: jwt.NewNumericDate(t.ExpiresAt),
			NotBefore: jwt.NewNumericDate(t.NotBefore),
			IssuedAt:  jwt.NewNumericDate(t.IssuedAt),
			ID:        t.Id.String(),
			Subject:   fmt.Sprint(t.UserId),
		},
		Roles: t.Roles,
	}

	token := jwt.NewWithClaims(m.signingMethod, claims)
	ss, err := token.SignedString(m.signingKey)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return ss, nil
}

func (m AccessTokenMapper) FromSignedString(ss string) (model.AccessToken, error) {
	const op = "mapper.AccessTokenMapper.FromSignedString"

	claims := accessTokenClaims{}
	_, err := jwt.ParseWithClaims(ss, &claims, func(token *jwt.Token) (interface{}, error) {
		const op = "mapper.AccessTokenMapper.FromSignedString.ParseWithClaims@lambda"

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%s: unexpected signing method: %v", op, token.Method.Alg())
		}

		return m.signingKey, nil
	})
	if err != nil {
		return model.AccessToken{}, fmt.Errorf("%s: %w", op, err)
	}

	id, err := uuid.Parse(claims.ID)
	if err != nil {
		return model.AccessToken{}, fmt.Errorf("%s: %w", op, err)
	}

	var userId uint
	_, err = fmt.Sscan(claims.Subject, &userId)
	if err != nil {
		return model.AccessToken{}, fmt.Errorf("%s: %w", op, err)
	}

	return model.AccessToken{
		Issuer:    claims.Issuer,
		ExpiresAt: claims.ExpiresAt.Time,
		NotBefore: claims.NotBefore.Time,
		IssuedAt:  claims.IssuedAt.Time,
		Id:        id,
		UserId:    userId,
		Roles:     claims.Roles,
	}, nil
}
