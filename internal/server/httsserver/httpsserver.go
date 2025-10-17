package httpsserver

import (
	"context"
	"net/http"

	"github.com/kirillmashkov/GophKeeper.git/internal/server"
	"github.com/kirillmashkov/GophKeeper.git/internal/server/httsserver/handler"
	"go.uber.org/zap"

	"golang.org/x/crypto/acme/autocert"
)

type HTTPS struct {
	server *http.Server
}

func (s *HTTPS) Run() error {
	return s.server.ListenAndServe()
}

// Shutdown - остановка https сервера
func (s *HTTPS) Shutdown() error {
	return s.server.Shutdown(context.Background())
}

func NewHTTPS(addr string, authService handler.IAuthService, secretService handler.ISecretService, logger *zap.Logger) server.Server {
	httpServer := &http.Server{
		Addr:    addr,
		Handler: handler.NewRouter(authService, secretService, logger),
	}

	manager := &autocert.Manager{
		Cache:      autocert.DirCache("cache-dir"),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("mysite.ru", "www.mysite.ru"),
	}

	httpServer.TLSConfig = manager.TLSConfig()

	server := &HTTPS{
		server: httpServer,
	}

	return server
}
