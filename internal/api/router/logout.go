package router

import (
	"log/slog"

	"github.com/fatalistix/golang-jwt-server/internal/api/handler"
	usecase "github.com/fatalistix/golang-jwt-server/internal/usecase/logout"
	"github.com/gin-gonic/gin"
)

func Logout(
	engine *gin.Engine,
	log *slog.Logger,
	usecase *usecase.Usecase,
) {
	engine.POST(
		"/logout",
		handler.MakeLogoutHandlerFunc(log, usecase),
	)
}
