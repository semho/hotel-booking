package middleware

import (
	"context"
	"github.com/semho/hotel-booking/api-gateway/internal/constants"
	"github.com/semho/hotel-booking/pkg/logger"
	pb "github.com/semho/hotel-booking/pkg/proto/auth_v1/auth"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	authClient pb.AuthServiceClient
}

func NewAuthMiddleware(authClient pb.AuthServiceClient) *AuthMiddleware {
	return &AuthMiddleware{
		authClient: authClient,
	}
}

func (m *AuthMiddleware) ValidateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			logger.Log.Info(
				"starting token validation",
				"path", r.URL.Path,
				"method", r.Method,
			)
			// Извлекаем токен из заголовка
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				logger.Log.Info("missing authorization header")
				http.Error(w, "missing authorization header", http.StatusUnauthorized)
				return
			}

			// Проверяем формат токена
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				logger.Log.Info(
					"invalid authorization header format",
					"parts_length", len(parts),
					"first_part", parts[0],
				)
				http.Error(w, "invalid authorization header format", http.StatusUnauthorized)
				return
			}

			// Валидируем токен через auth service
			resp, err := m.authClient.Validate(
				r.Context(), &pb.ValidateRequest{
					AccessToken: parts[1],
				},
			)
			if err != nil {
				logger.Log.Error(
					"failed to validate token",
					"error", err,
				)
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			// Добавляем информацию о пользователе в контекст
			ctx := context.WithValue(r.Context(), constants.USER, resp.User)
			logger.Log.Info(
				"token validated successfully",
				"user_id", resp.User.Id,
				"user_email", resp.User.Email,
			)
			next.ServeHTTP(w, r.WithContext(ctx))
		},
	)
}
