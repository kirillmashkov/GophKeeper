package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/kirillmashkov/GophKeeper.git/internal/server/httsserver/middleware/security"
	"github.com/kirillmashkov/GophKeeper.git/internal/server/model"
	"go.uber.org/zap"
)

// IAuthService определяет интерфейс для сервиса аутентификации
type IAuthService interface {
	SignUp(ctx context.Context, email string, password string) (string, error)
	SignIn(ctx context.Context, email string, password string) (bool, string, error)
}

type ISecretService interface {
	CreateSecret(ctx context.Context, userID string, createSecret model.CreateSecretRequest) error
	UpdateSecret(ctx context.Context, userID string, updatesecret model.UpdateSecretRequest) error
	GetSecret(ctx context.Context, userID string, name string) (model.GetSecretResponse, error)
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusBadRequest)
		return
	}

	var request model.SignUpRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		h.Logger.Error("cannot parse request JSON body", zap.Error(err))
		http.Error(w, "cannot parse request JSON body", http.StatusBadRequest)
		return
	}

	token, err := h.UserService.SignUp(r.Context(), request.Email, request.Password)
	if err != nil {
		if errors.Is(err, model.ErrUserConflict) {
			w.WriteHeader(http.StatusConflict)
			return
		}

		errorString := fmt.Sprintf("Something went wrong when create user %s", request.Email)
		http.Error(w, errorString, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Authorization", "Bearer "+token)
	w.WriteHeader(http.StatusCreated)

}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusBadRequest)
		return
	}

	var request model.SignUpRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		h.Logger.Error("cannot parse request JSON body", zap.Error(err))
		http.Error(w, "cannot parse request JSON body", http.StatusBadRequest)
		return
	}

	checkUser, token, err := h.UserService.SignIn(r.Context(), request.Email, request.Password)
	if checkUser && err == nil {
		w.Header().Set("Authorization", "Bearer "+token)
		w.WriteHeader(http.StatusOK)
		return
	}

	if checkUser && err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !checkUser {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
}

func (h *Handler) CreateSecret(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusBadRequest)
		return
	}

	var request model.CreateSecretRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		h.Logger.Error("cannot parse request JSON body", zap.Error(err))
		http.Error(w, "cannot parse request JSON body", http.StatusBadRequest)
		return
	}
	
	u := security.UserIDType("userID")
	userID := fmt.Sprintf("%v", r.Context().Value(u))
	if userID == "" {
		h.Logger.Error("No userID. Unauthorized")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	err := h.SecretService.CreateSecret(r.Context(), userID, request)

	if err != nil {
		if errors.Is(err, model.ErrSecretConflict) {
			w.WriteHeader(http.StatusConflict)
			return
		}
		errorString := fmt.Sprintf("Something went wrong when create secret %s", request.Name)
		http.Error(w, errorString, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) UpdateSecret(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Only PUT requests are allowed!", http.StatusBadRequest)
		return
	}

	var request model.UpdateSecretRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		h.Logger.Error("cannot parse request JSON body", zap.Error(err))
		http.Error(w, "cannot parse request JSON body", http.StatusBadRequest)
		return
	}
	
	u := security.UserIDType("userID")
	userID := fmt.Sprintf("%v", r.Context().Value(u))
	if userID == "" {
		h.Logger.Error("No userID. Unauthorized")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	err := h.SecretService.UpdateSecret(r.Context(), userID, request)
	if err != nil {
		if errors.Is(err, model.ErrSecretNoFound) {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		errorString := fmt.Sprintf("Something went wrong when update secret %s", request.Name)
		http.Error(w, errorString, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetSecret(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only Get requests are allowed!", http.StatusBadRequest)
		return
	}

	u := security.UserIDType("userID")
	userID := fmt.Sprintf("%v", r.Context().Value(u))
	if userID == "" {
		h.Logger.Error("No userID. Unauthorized")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	name := r.URL.Path[len("/api/"):]

	secret, err := h.SecretService.GetSecret(r.Context(), userID, name)
	if err != nil {
		if errors.Is(err, model.ErrSecretNoFound) {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		errorString := fmt.Sprintf("Something went wrong when update secret %s", name)
		http.Error(w, errorString, http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(secret); err != nil {
		h.Logger.Debug("error encoding result", zap.Error(err))
		return
	}
}