package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/semho/hotel-booking/api-gateway/internal/api/http/request"
	"github.com/semho/hotel-booking/api-gateway/internal/api/http/response"
	"github.com/semho/hotel-booking/pkg/logger"
	pb "github.com/semho/hotel-booking/pkg/proto/auth_v1/auth"
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
			r.Post("/validate", h.Validate)
		},
	)
}

func (h *AuthHandler) Validate(w http.ResponseWriter, r *http.Request) {
	var req request.ValidateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.authClient.Validate(
		r.Context(), &pb.ValidateRequest{
			AccessToken: req.AccessToken,
		},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Log.Error("api-gateway validate with authClient", err)
		return
	}

	logger.Log.Info(
		"received response from auth service",
		"valid", resp.Valid,
		"user", resp.User,
	)

	// Отправляем ответ клиенту
	authResponse := response.ValidateResponse{
		Valid:    resp.Valid,
		UserInfo: response.UserFromProto(resp.User),
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(authResponse)
	if err != nil {
		logger.Log.Error("api-gateway validate ", err)
		return
	}
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

	logger.Log.Info(
		"received response from auth service",
		"access_token", resp.AccessToken,
		"refresh_token", resp.RefreshToken,
		"user", resp.User,
	)

	// Устанавливаем refresh token в cookie
	if resp.RefreshToken != "" {
		logger.Log.Info(
			"api-gateway: setting refresh token cookie",
			"refresh_token_length", len(resp.RefreshToken),
		)
		cookie := &http.Cookie{
			Name:     "refresh_token",
			Value:    resp.RefreshToken,
			HttpOnly: true,
			Secure:   false,
			Path:     "/",
			MaxAge:   60 * 60 * 24 * 7, // 7 дней
			SameSite: http.SameSiteLaxMode,
		}
		http.SetCookie(w, cookie)

		// Проверяем, что cookie был установлен
		cookies := w.Header().Values("Set-Cookie")
		logger.Log.Info(
			"api-gateway: cookies after setting",
			"all_cookies", cookies,
			"cookie_count", len(cookies),
			"cookie_details", map[string]interface{}{
				"name":      cookie.Name,
				"path":      cookie.Path,
				"secure":    cookie.Secure,
				"httpOnly":  cookie.HttpOnly,
				"sameSite":  cookie.SameSite,
				"maxAge":    cookie.MaxAge,
				"raw_value": cookie.String(),
			},
		)
	} else {
		logger.Log.Warn(
			"api-gateway: refresh token is empty",
			"grpc_response", resp,
		)
	}

	// Отправляем ответ клиенту
	authResponse := response.AuthResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		User:         response.UserFromProto(resp.User),
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

	logger.Log.Info(
		"api-gateway: received grpc response",
		"access_token_length", len(resp.AccessToken),
		"refresh_token_length", len(resp.RefreshToken),
		"access_token_empty", resp.AccessToken == "",
		"refresh_token_empty", resp.RefreshToken == "",
	)

	// Отправляем ответ клиенту
	authResponse := response.AuthResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		User:         response.UserFromProto(resp.User),
	}

	logger.Log.Info(
		"api-gateway: sending response",
		"access_token_length", len(authResponse.AccessToken),
		"refresh_token_length", len(authResponse.RefreshToken),
		"access_token_empty", authResponse.AccessToken == "",
		"refresh_token_empty", authResponse.RefreshToken == "",
	)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(authResponse)
	if err != nil {
		logger.Log.Error("api-gateway login encoding error", err)
		return
	}
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	// Читаем refresh token из тела запроса
	var req request.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.Error(
			"api-gateway: failed to decode refresh request",
			"error", err,
		)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.RefreshToken == "" {
		logger.Log.Error("api-gateway: refresh token is empty in request")
		http.Error(w, "refresh token is required", http.StatusBadRequest)
		return
	}

	resp, err := h.authClient.Refresh(
		r.Context(), &pb.RefreshRequest{
			RefreshToken: req.RefreshToken,
		},
	)
	if err != nil {
		logger.Log.Error(
			"api-gateway: refresh token validation failed",
			"error", err,
		)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	logger.Log.Info(
		"api-gateway: received new tokens from auth service",
		"access_token_length", len(resp.AccessToken),
		"refresh_token_length", len(resp.RefreshToken),
	)

	// Отправляем ответ клиенту
	authResponse := response.AuthResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		User:         response.UserFromProto(resp.User),
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(authResponse)
	if err != nil {
		logger.Log.Error("api-gateway refresh encoding error", err)
		return
	}
}
