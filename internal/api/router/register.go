package router

import (
	"log/slog"

	"github.com/fatalistix/golang-jwt-server/internal/api/handler"
	usecase "github.com/fatalistix/golang-jwt-server/internal/usecase/register"
	"github.com/gin-gonic/gin"
)

func Register(
	engine *gin.Engine,
	log *slog.Logger,
	usecase *usecase.Usecase,
) {
	engine.POST(
		"/register",
		handler.MakeRegisterHandlerFunc(log, usecase),
	)
}
