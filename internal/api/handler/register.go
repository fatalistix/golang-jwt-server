package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/fatalistix/golang-jwt-server/internal/api/dto"
	"github.com/fatalistix/slogattr"
	"github.com/gin-gonic/gin"
)

type RegisterHandlerFunc = func(c *gin.Context)

type RegisterHandler interface {
	Handle(ctx context.Context, username, password string) (uint, error)
}

func MakeRegisterHandlerFunc(log *slog.Logger, handler RegisterHandler) RegisterHandlerFunc {
	const op = "api.handler.MakeRegisterHandlerFunc"

	log = log.With(
		slog.String("op", op),
	)

	return func(c *gin.Context) {
		var request dto.RegisterRequest

		if err := c.BindJSON(&request); err != nil {
			log.Error("cannot bind request", slogattr.Err(err))

			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Message: "Invalid request body",
			})

			return
		}

		id, err := handler.Handle(c, request.Username, request.Password)
		if err != nil {
			log.Error("cannot register user", slogattr.Err(err))

			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Message: "Cannot register user",
			})

			return
		}

		c.JSON(http.StatusCreated, dto.RegisterResponse{
			Id: id,
		})
	}
}
