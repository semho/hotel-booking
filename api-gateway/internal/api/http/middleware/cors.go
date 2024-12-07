package middleware

import (
	"github.com/go-chi/cors"
	"github.com/semho/hotel-booking/api-gateway/internal/config"
	"net/http"
)

func CORS(cfg config.CORSConfig) func(http.Handler) http.Handler {
	return cors.Handler(
		cors.Options{
			AllowedOrigins:   cfg.Origins, // Разрешенные Origin
			AllowedMethods:   cfg.Methods,
			AllowedHeaders:   cfg.Headers,        // Разрешенные заголовки
			ExposedHeaders:   cfg.ExposedHeaders, // Заголовки, доступные клиенту
			AllowCredentials: cfg.Credentials,    // Разрешить передачу cookies
			MaxAge:           cfg.MaxAge,         // Кэширование preflight запросов в секундах
			Debug:            cfg.Debug,          // Включить отладку (для разработки)
		},
	)
}
