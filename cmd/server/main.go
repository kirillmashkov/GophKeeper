package main

import (
	"flag"

	"github.com/kirillmashkov/GophKeeper.git/internal/app"
	"github.com/kirillmashkov/GophKeeper.git/internal/server"
	httpsserver "github.com/kirillmashkov/GophKeeper.git/internal/server/httsserver"
	"go.uber.org/zap"
)

func main() {
	flag.Parse()
	err := app.Initialize()
	if err != nil {
		panic(err)
	}

	var restServer server.Server
	app.Logger.Info("%s", zap.String("addr", app.Config.ServerAddress))
	restServer = httpsserver.NewHTTPS(app.Config.ServerAddress, app.AuthService, app.SecretService, app.Logger)
	restServer.Run()

	defer app.Close()
}