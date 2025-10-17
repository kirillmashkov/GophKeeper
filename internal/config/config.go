package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"
)

type Config struct {
	ServerAddress 	string "env:\"SERVER_ADDRESS\""
	DB 				string	"env:\"DB\""
}

var serverArg Config
var serverEnv Config

func init() {
	flag.StringVar(&serverArg.ServerAddress, "a", "localhost:8080", "server host")
	flag.StringVar(&serverArg.DB, "d", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable", "db connection string")
}

func InitConfig(l *zap.Logger) Config {
	err := env.Parse(&serverEnv)
	if err != nil {
		l.Error("Can't read env variables")
	}

	var serverConfig Config

	serverConfig.DB = getConfigString(serverEnv.DB, serverArg.DB)
	serverConfig.ServerAddress = getConfigString(serverEnv.ServerAddress, serverArg.ServerAddress)

	return serverConfig
}

func getConfigString(env string, arg string) string {
	if env == "" {
		return arg
	} else {
		return env
	}
}