package security

import (
	"context"
	"net/http"
	"strings"

	"github.com/kirillmashkov/GophKeeper.git/internal/app"
	"go.uber.org/zap"
)

type UserIDType string

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		u := UserIDType("userID")
		if auth == "" {
			app.Logger.Warn("No token")
			c := context.WithValue(r.Context(), u, "")
			next.ServeHTTP(w, r.WithContext(c))
			return
		}

		token := (strings.Split(auth, "Bearer "))[1]
		userID, err := getUserIDFromToken(token)

		app.Logger.Info("userID", zap.String("userID", userID))
		if err != nil {
			c := context.WithValue(r.Context(), u, "")
			next.ServeHTTP(w, r.WithContext(c))
			return
		}

		c := context.WithValue(r.Context(), u, userID)
		next.ServeHTTP(w, r.WithContext(c))

	})
}

func getUserIDFromToken(tokenString string) (string, error) {
	token, claims, err := app.SecurityUtil.ParseJWT(tokenString)

	if err != nil {
		app.Logger.Error("Can't parse token", zap.Error(err))
		return "", err
	}

	if !token.Valid {
		app.Logger.Warn("Token is not valid")
		return "", err
	}

	if claims.UserID == "" {
		app.Logger.Warn("Token doesn't contain UserID")
		return "", err
	}

	return claims.UserID, nil
}

