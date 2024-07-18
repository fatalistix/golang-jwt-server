package api

import (
	"log/slog"

	ginlogger "github.com/FabienMht/ginslog/logger"
	ginrecovery "github.com/FabienMht/ginslog/recovery"
	"github.com/fatalistix/golang-jwt-server/internal/api/router"
	"github.com/fatalistix/golang-jwt-server/internal/config"
	"github.com/fatalistix/golang-jwt-server/internal/usecase/login"
	"github.com/fatalistix/golang-jwt-server/internal/usecase/logout"
	"github.com/fatalistix/golang-jwt-server/internal/usecase/refresh"
	"github.com/fatalistix/golang-jwt-server/internal/usecase/register"
	"github.com/gin-gonic/gin"
)

func NewRouter(cfg config.GinConfig, log *slog.Logger) *gin.Engine {
	r := gin.Default()
	r.Use(ginlogger.New(log))
	r.Use(ginrecovery.New(log))
	gin.SetMode(cfg.RunMode)

	return r
}

func RegisterHandlers(
	engine *gin.Engine,
	log *slog.Logger,
	login *login.Usecase,
	logout *logout.Usecase,
	refresh *refresh.Usecase,
	register *register.Usecase,
) {
	router.Login(engine, log, login)
	router.Logout(engine, log, logout)
	router.Refresh(engine, log, refresh)
	router.Register(engine, log, register)
}
