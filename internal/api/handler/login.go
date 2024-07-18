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

type LoginHandlerFunc = func(c *gin.Context)

type LoginHandler interface {
	Handle(ctx context.Context, username, password string) (modeldto.Tokens, error)
}

func MakeLoginHandlerFunc(log *slog.Logger, handler LoginHandler) LoginHandlerFunc {
	const op = "api.handler.MakeLoginHandlerFunc"

	log = log.With(
		slog.String("op", op),
	)

	return func(c *gin.Context) {
		var request dto.LoginRequest

		if err := c.BindJSON(&request); err != nil {
			log.Error("cannot bind request", slogattr.Err(err))

			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Message: "Invalid request body",
			})

			return
		}

		tokens, err := handler.Handle(c, request.Username, request.Password)
		if err != nil {
			log.Error("cannot login user", slogattr.Err(err))

			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Message: "Cannot login user",
			})

			return
		}

		c.JSON(http.StatusOK, dto.LoginResponse{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		})
	}
}
