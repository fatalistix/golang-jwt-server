package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	trmgorm "github.com/avito-tech/go-transaction-manager/drivers/gorm/v2"
	"github.com/fatalistix/golang-jwt-server/internal/api"
	"github.com/fatalistix/golang-jwt-server/internal/config"
	"github.com/fatalistix/golang-jwt-server/internal/domain/maker"
	"github.com/fatalistix/golang-jwt-server/internal/infrastructure/persistence/database"
	"github.com/fatalistix/golang-jwt-server/internal/infrastructure/persistence/migration"
	"github.com/fatalistix/golang-jwt-server/internal/infrastructure/persistence/repository"
	"github.com/fatalistix/golang-jwt-server/internal/mapper"
	"github.com/fatalistix/golang-jwt-server/internal/usecase/login"
	"github.com/fatalistix/golang-jwt-server/internal/usecase/logout"
	"github.com/fatalistix/golang-jwt-server/internal/usecase/refresh"
	"github.com/fatalistix/golang-jwt-server/internal/usecase/register"
	"github.com/fatalistix/golang-jwt-server/internal/worker"
	"gorm.io/gorm"
)

type App struct {
	db                  *gorm.DB
	server              *http.Server
	refreshTokenCleaner worker.RefreshTokenCleaner
}

func New(log *slog.Logger, cfg config.Config) (*App, error) {
	const op = "app.New"

	db, err := database.NewPostgres(cfg.Postgres)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = migration.Apply(db)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	accessTokenMapper := mapper.MakeAccessTokenMapper(cfg.Token.Secret)
	refreshTokenMapper := mapper.MakeRefreshTokenMapper(cfg.Token.Secret)
	userMapper := mapper.MakeUserMapper()

	tokenMaker := maker.MakeTokenMaker(cfg.Token)

	userRepository := repository.NewUserRepository(db, trmgorm.DefaultCtxGetter, userMapper)
	refreshTokenRepository := repository.NewRefreshTokenRepository(db, trmgorm.DefaultCtxGetter, refreshTokenMapper)

	loginUsecase := login.NewUsecase(
		userRepository,
		refreshTokenRepository,
		tokenMaker,
		accessTokenMapper,
		refreshTokenMapper,
	)

	logoutUsecase := logout.NewUsecase(
		refreshTokenMapper,
		refreshTokenRepository,
	)

	refreshUsecase := refresh.NewUsecase(
		accessTokenMapper,
		refreshTokenMapper,
		refreshTokenRepository,
		tokenMaker,
		userRepository,
	)

	registerUsecase := register.NewUsecase(
		userRepository,
		cfg.PasswordEncoder.EncryptCost,
	)

	router := api.NewRouter(cfg.Gin, log)
	api.RegisterHandlers(
		router, log,
		loginUsecase,
		logoutUsecase,
		refreshUsecase,
		registerUsecase,
	)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Http.Port),
		Handler:      router,
		ReadTimeout:  cfg.Http.ReadTimeout,
		WriteTimeout: cfg.Http.WriteTimeout,
		IdleTimeout:  cfg.Http.IdleTimeout,
	}

	refreshTokenCleaner := worker.MakeRefreshTokenCleaner(
		log,
		cfg.Token.CleanTimeout,
		refreshTokenRepository,
	)

	return &App{
		db:                  db,
		server:              server,
		refreshTokenCleaner: refreshTokenCleaner,
	}, nil
}

func (a *App) Run() error {
	const op = "app.Run"

	a.refreshTokenCleaner.Start()

	if err := a.server.ListenAndServe(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop(ctx context.Context) error {
	const op = "app.Stop"

	defer func() {
		con, _ := a.db.DB()
		_ = con.Close()
	}()

	if err := a.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	a.refreshTokenCleaner.Stop()

	return nil
}
