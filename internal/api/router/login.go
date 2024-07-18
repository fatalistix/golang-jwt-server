package router

import (
	"log/slog"

	"github.com/fatalistix/golang-jwt-server/internal/api/handler"
	usecase "github.com/fatalistix/golang-jwt-server/internal/usecase/login"
	"github.com/gin-gonic/gin"
)

func Login(
	engine *gin.Engine,
	log *slog.Logger,
	usecase *usecase.Usecase,
) {
	engine.POST(
		"/login",
		handler.MakeLoginHandlerFunc(log, usecase),
	)
}
