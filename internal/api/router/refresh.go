package router

import (
	"log/slog"

	"github.com/fatalistix/golang-jwt-server/internal/api/handler"
	usecase "github.com/fatalistix/golang-jwt-server/internal/usecase/refresh"
	"github.com/gin-gonic/gin"
)

func Refresh(
	engine *gin.Engine,
	log *slog.Logger,
	usecase *usecase.Usecase,
) {
	engine.POST(
		"/refresh",
		handler.MakeRefreshHandlerFunc(log, usecase),
	)
}
