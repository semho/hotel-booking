package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/semho/hotel-booking/api-gateway/internal/api/http/mapper"
	"github.com/semho/hotel-booking/api-gateway/internal/api/http/middleware"
	"github.com/semho/hotel-booking/api-gateway/internal/api/http/request"
	"github.com/semho/hotel-booking/api-gateway/internal/api/http/response"
	"github.com/semho/hotel-booking/pkg/errors"
	"github.com/semho/hotel-booking/pkg/logger"
	roompb "github.com/semho/hotel-booking/pkg/proto/room_v1/room"
)

type RoomHandler struct {
	roomClient     roompb.RoomServiceClient
	authMiddleware *middleware.AuthMiddleware
}

func NewRoomHandler(
	roomClient roompb.RoomServiceClient,
	authMiddleware *middleware.AuthMiddleware,
) *RoomHandler {
	return &RoomHandler{
		roomClient:     roomClient,
		authMiddleware: authMiddleware,
	}
}

func (h *RoomHandler) RegisterRoutes(r chi.Router) {
	r.Route(
		"/api/v1/rooms", func(r chi.Router) {

			// Публичные маршруты
			//r.Group(
			//	func(r chi.Router) {
			//
			//	},
			//)
			// Защищенные маршруты
			r.Group(
				func(r chi.Router) {
					r.Use(h.authMiddleware.ValidateToken)
					r.Post("/", h.CreateRoom)
				},
			)

		},
	)

	// Добавим также корневой маршрут для проверки работоспособности API
	r.Get(
		"/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		},
	)
}

func (h *RoomHandler) respondWithJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Log.Error("failed to encode response", "error", err)
	}
}

func (h *RoomHandler) respondWithError(w http.ResponseWriter, code int, err error) {
	var errorCode string
	switch {
	case errors.IsNotFound(err):
		errorCode = "NOT_FOUND"
	case errors.IsConflict(err):
		errorCode = "CONFLICT"
	case errors.IsInvalidInput(err):
		errorCode = "INVALID_INPUT"
	default:
		errorCode = "INTERNAL_ERROR"
	}

	h.respondWithJSON(
		w, code, response.Error{
			Code:    errorCode,
			Message: err.Error(),
		},
	)
}

func (h *RoomHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	//TODO: сделать проверку на админа. Лучше не тут, а сделать дополнительный middleware isAdmin
	//userInfo, ok := r.Context().Value(app.USER).(*authpb.UserInfo)
	//if !ok {
	//	http.Error(w, "unauthorized", http.StatusUnauthorized)
	//	return
	//}

	var req request.CreateRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.Error("failed to decode request body", "error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	protoReq := mapper.HttpToProto(req)
	resp, err := h.roomClient.CreateRoom(ctx, protoReq)
	if err != nil {
		logger.Log.Error("failed to create rooms", "error", err)
		h.respondWithError(
			w, http.StatusInternalServerError,
			errors.WithMessage(errors.ErrInternal, "failed to create rooms"),
		)
		return
	}

	// Конвертируем ответ в HTTP формат
	rooms := mapper.ToHTTPRoom(resp.GetRoom())

	// Отправляем ответ
	h.respondWithJSON(w, http.StatusOK, rooms)
}
