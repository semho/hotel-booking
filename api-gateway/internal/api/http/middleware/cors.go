package middleware

import (
	"github.com/go-chi/cors"
	"net/http"
)

func CORS() func(http.Handler) http.Handler {
	return cors.Handler(
		cors.Options{
			AllowedOrigins: []string{"http://localhost:8081"}, // Разрешенные Origin
			AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowedHeaders: []string{
				"Accept",
				"Authorization",
				"Content-Type",
				"X-CSRF-Token",
				"X-Requested-With",
				"Force-Country-Code",
				"Geo-Ip-2-Country",
			}, // Разрешенные заголовки
			ExposedHeaders: []string{
				"Link",
				"Set-Cookie",
			}, // Заголовки, доступные клиенту
			AllowCredentials: true, // Разрешить передачу cookies
			MaxAge:           300,  // Кэширование preflight запросов в секундах
			Debug:            true, // Включить отладку (для разработки)
		},
	)
}
