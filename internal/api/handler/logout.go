package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/fatalistix/golang-jwt-server/internal/api/dto"
	"github.com/fatalistix/slogattr"
	"github.com/gin-gonic/gin"
)

type LogoutHandlerFunc = func(c *gin.Context)

type LogoutHandler interface {
	Handle(ctx context.Context, refreshToken string) error
}

func MakeLogoutHandlerFunc(log *slog.Logger, handler LogoutHandler) LogoutHandlerFunc {
	const op = "api.handler.MakeLogoutHandlerFunc"

	log = log.With(
		slog.String("op", op),
	)

	return func(c *gin.Context) {
		var request dto.LogoutRequest

		if err := c.BindJSON(&request); err != nil {
			log.Error("cannot bind request", slogattr.Err(err))

			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Message: "Invalid request body",
			})

			return
		}

		err := handler.Handle(c, request.RefreshToken)
		if err != nil {
			log.Error("cannot logout user", slogattr.Err(err))

			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Message: "Cannot logout user",
			})

			return
		}

		c.JSON(http.StatusCreated, dto.LogoutResponse{})
	}
}
