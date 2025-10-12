package handler

import (
	"github.com/go-chi/chi"
	"github.com/kirillmashkov/GophKeeper.git/internal/server/httsserver/middleware/security"
	"go.uber.org/zap"
)

type Handler struct {
	Mux         *chi.Mux
	UserService IAuthService
	SecretService ISecretService
	Logger *zap.Logger
}

func NewRouter(authService IAuthService, secretService ISecretService, logger *zap.Logger) chi.Router {
	r := chi.NewRouter()

	h := NewHandler(authService, secretService, logger)
	r.Post("/api/signup", h.SignUp)
	r.Post("/api/signin", h.SignIn)
	r.Group(func (r chi.Router) {
		r.Use(security.Auth)
		r.Post("/api/createsecret", h.CreateSecret)
		r.Put("/api/updatesecret", h.UpdateSecret)
		r.Get("/api/{name}", h.GetSecret)
	})

	return r
}

func NewHandler(authService IAuthService, secretService ISecretService, logger *zap.Logger) *Handler {
	return &Handler{
		Mux:         chi.NewMux(),
		UserService: authService,
		SecretService: secretService,
		Logger: logger,
	}
}
