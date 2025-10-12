package app

import (
	"log"
	"time"

	"github.com/kirillmashkov/GophKeeper.git/internal/config"
	"github.com/kirillmashkov/GophKeeper.git/internal/logger"
	"github.com/kirillmashkov/GophKeeper.git/internal/server/service"
	"github.com/kirillmashkov/GophKeeper.git/internal/storage"
	"github.com/kirillmashkov/GophKeeper.git/internal/util"
	"go.uber.org/zap"
)

var Logger *zap.Logger
var Config config.Config
var Database *storage.Database
var AuthService *service.AuthService
var SecretService *service.SecretService
var SecurityUtil *util.SecurityUtil

const tokenExp = time.Hour * 3
const SecretKey = "supersecretkeysupersecretkeysupersecretkeysupersecretkeysupersecretkey"

func Initialize() error {
	var err error

	Logger, err = logger.Initialize()

	if err != nil {
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
		log.SetPrefix("ERROR: ")
		log.Printf("Can't init logger %v", err)
		return err
	}

	Config = config.InitConfig(Logger)

	Database = storage.NewDatabase(&Config, Logger)
	err = Database.Open()
	if err != nil {
		Logger.Error("error open database", zap.Error(err))
		return err
	}

	if err = Database.Migrate(); err != nil {
		return err
	}

	SecurityUtil = util.NewSecurityUtil(tokenExp, SecretKey)
	userStorage := storage.NewUserRepository(Database, Logger)
	secretRepository := storage.NewSecretRepository(Database, Logger)
	
	AuthService = service.NewAuthService(userStorage, Logger, SecurityUtil)
	SecretService = service.NewSecretService(secretRepository, Logger)

	return nil
}

func Close() {
	Database.Close()
}
