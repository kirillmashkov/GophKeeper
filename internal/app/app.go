package app

import (
	"context"
	"log"
	"time"

	"github.com/kirillmashkov/GophKeeper.git/internal/config"
	"github.com/kirillmashkov/GophKeeper.git/internal/logger"
	"github.com/kirillmashkov/GophKeeper.git/internal/server/service"
	"github.com/kirillmashkov/GophKeeper.git/internal/storage"
	"github.com/kirillmashkov/GophKeeper.git/internal/util"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Backward-compatible globals (populated by fx on startup)
var (
	Logger        *zap.Logger
	Config        config.Config
	Database      *storage.Database
	AuthService   *service.AuthService
	SecretService *service.SecretService
	SecurityUtil  *util.SecurityUtil
)

const (
	tokenExp  = time.Hour * 3
	SecretKey = "supersecretkeysupersecretkeysupersecretkeysupersecretkeysupersecretkey"
)

var fxApp *fx.App

// Initialize builds the dependency graph using Uber FX and starts the app lifecycle.
// It preserves the original behavior by exposing the same globals and function signature.
func Initialize() error {
	fxApp = fx.New(
		fx.Provide(
			// Core
			logger.Initialize,
			provideConfig,
			storage.NewDatabase,
			provideSecurityUtil,

			// Repositories (explicit wrappers to bind to service interfaces)
			func(db *storage.Database, l *zap.Logger) service.IUserRepository {
				return storage.NewUserRepository(db, l)
			},
			func(db *storage.Database, l *zap.Logger) service.ISecretRepository {
				return storage.NewSecretRepository(db, l)
			},

			// Services
			service.NewAuthService,
			service.NewSecretService,
		),
		fx.Invoke(
			openAndMigrateDB,
			populateGlobals,
		),
	)

	if err := fxApp.Start(context.Background()); err != nil {
		// Fallback to std logger on very early initialization errors
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
		log.SetPrefix("ERROR: ")
		log.Printf("Can't start fx app: %v", err)
		return err
	}
	return nil
}

// Close gracefully stops the fx application and releases resources.
func Close() {
	if fxApp != nil {
		_ = fxApp.Stop(context.Background())
	}
}

// provideConfig adapts Config initialization for fx and returns a pointer
// because many constructors expect *config.Config.
func provideConfig(l *zap.Logger) *config.Config {
	c := config.InitConfig(l)
	return &c
}

func provideSecurityUtil() *util.SecurityUtil {
	return util.NewSecurityUtil(tokenExp, SecretKey)
}

// openAndMigrateDB wires DB lifecycle with fx and runs migrations on start.
func openAndMigrateDB(lc fx.Lifecycle, db *storage.Database, l *zap.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := db.Open(); err != nil {
				l.Error("error open database", zap.Error(err))
				return err
			}
			if err := db.Migrate(); err != nil {
				return err
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			db.Close()
			return nil
		},
	})
}

// populateGlobals keeps backward compatibility for the rest of the codebase
// that directly references app package globals.
func populateGlobals(
	l *zap.Logger,
	cfg *config.Config,
	db *storage.Database,
	auth *service.AuthService,
	secret *service.SecretService,
	su *util.SecurityUtil,
) {
	Logger = l
	Config = *cfg
	Database = db
	AuthService = auth
	SecretService = secret
	SecurityUtil = su
}
