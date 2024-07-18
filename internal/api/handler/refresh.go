package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/fatalistix/golang-jwt-server/internal/api/dto"
	modeldto "github.com/fatalistix/golang-jwt-server/internal/dto"
	"github.com/fatalistix/slogattr"
	"github.com/gin-gonic/gin"
)

type RefreshHandlerFunc = func(c *gin.Context)

type RefreshHandler interface {
	Handle(ctx context.Context, refreshToken string) (modeldto.Tokens, error)
}

func MakeRefreshHandlerFunc(log *slog.Logger, handler RefreshHandler) RefreshHandlerFunc {
	const op = "api.handler.MakeLogoutHandlerFunc"

	log = log.With(
		slog.String("op", op),
	)

	return func(c *gin.Context) {
		var request dto.RefreshRequest

		if err := c.BindJSON(&request); err != nil {
			log.Error("cannot bind request", slogattr.Err(err))

			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Message: "Invalid request body",
			})

			return
		}

		tokens, err := handler.Handle(c, request.RefreshToken)
		if err != nil {
			log.Error("cannot refresh tokens", slogattr.Err(err))

			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Message: "Cannot refresh tokens",
			})

			return
		}

		c.JSON(http.StatusCreated, dto.RefreshResponse{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		})
	}
}
