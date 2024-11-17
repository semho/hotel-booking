package handler

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/semho/hotel-booking/api-gateway/internal/api/http/request"
	"github.com/semho/hotel-booking/api-gateway/internal/api/http/response"
	"github.com/semho/hotel-booking/pkg/logger"
	pb "github.com/semho/hotel-booking/pkg/proto/auth_v1/auth"
	"net/http"
)

type AuthHandler struct {
	authClient pb.AuthServiceClient
}

func NewAuthHandler(authClient pb.AuthServiceClient) *AuthHandler {
	return &AuthHandler{
		authClient: authClient,
	}
}

func (h *AuthHandler) RegisterRoutes(r chi.Router) {
	r.Route(
		"/api/v1/auth", func(r chi.Router) {
			r.Post("/register", h.Register)
			r.Post("/login", h.Login)
			r.Post("/refresh", h.Refresh)
		},
	)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req request.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.authClient.Register(
		r.Context(), &pb.RegisterRequest{
			Email:     req.Email,
			Password:  req.Password,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Phone:     req.Phone,
		},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Log.Error("api-gateway register with authClient", err)
		return
	}

	// Устанавливаем refresh token в cookie
	if resp.RefreshToken != "" {
		http.SetCookie(
			w, &http.Cookie{
				Name:     "refresh_token",
				Value:    resp.RefreshToken,
				HttpOnly: true,
				Secure:   true,
				Path:     "/api/v1/auth/refresh",
			},
		)
	}

	// Отправляем ответ клиенту только с access token
	authResponse := response.AuthResponse{
		AccessToken: resp.AccessToken,
		User:        response.UserFromProto(resp.User),
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(authResponse)
	if err != nil {
		logger.Log.Error("api-gateway register ", err)
		return
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req request.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.authClient.Login(
		r.Context(), &pb.LoginRequest{
			Email:    req.Email,
			Password: req.Password,
		},
	)
	if err != nil {
		logger.Log.Error("api-gateway login with authClient", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Устанавливаем refresh token в cookie
	if resp.RefreshToken != "" {
		http.SetCookie(
			w, &http.Cookie{
				Name:     "refresh_token",
				Value:    resp.RefreshToken,
				HttpOnly: true,
				Secure:   true,
				Path:     "/api/v1/auth/refresh",
			},
		)
	}

	// Отправляем ответ клиенту
	authResponse := response.AuthResponse{
		AccessToken: resp.AccessToken,
		User:        response.UserFromProto(resp.User),
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(authResponse)
	if err != nil {
		logger.Log.Error("api-gateway login ", err)
		return
	}
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	// Получаем refresh token из cookie
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "refresh token not found", http.StatusUnauthorized)
		return
	}

	resp, err := h.authClient.Refresh(
		r.Context(), &pb.RefreshRequest{
			RefreshToken: cookie.Value,
		},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Устанавливаем новый refresh token в cookie
	if resp.RefreshToken != "" {
		http.SetCookie(
			w, &http.Cookie{
				Name:     "refresh_token",
				Value:    resp.RefreshToken,
				HttpOnly: true,
				Secure:   true,
				Path:     "/api/v1/auth/refresh",
			},
		)
	}

	// Отправляем ответ клиенту
	authResponse := response.AuthResponse{
		AccessToken: resp.AccessToken,
		User:        response.UserFromProto(resp.User),
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(authResponse)
	if err != nil {
		logger.Log.Error("api-gateway refresh ", err)
		return
	}
}
