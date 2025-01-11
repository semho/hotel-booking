package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/semho/hotel-booking/api-gateway/internal/constants"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/semho/hotel-booking/api-gateway/internal/api/http/mapper"
	"github.com/semho/hotel-booking/api-gateway/internal/api/http/middleware"
	"github.com/semho/hotel-booking/api-gateway/internal/api/http/request"
	"github.com/semho/hotel-booking/api-gateway/internal/api/http/response"
	"github.com/semho/hotel-booking/pkg/errors"
	"github.com/semho/hotel-booking/pkg/logger"
	authpb "github.com/semho/hotel-booking/pkg/proto/auth_v1/auth"
	bookingpb "github.com/semho/hotel-booking/pkg/proto/booking_v1/booking"
	roompb "github.com/semho/hotel-booking/pkg/proto/room_v1/room"
)

type BookingHandler struct {
	bookingClient  bookingpb.BookingServiceClient
	authMiddleware *middleware.AuthMiddleware
}

func NewBookingHandler(
	bookingClient bookingpb.BookingServiceClient,
	authMiddleware *middleware.AuthMiddleware,
) *BookingHandler {
	return &BookingHandler{
		bookingClient:  bookingClient,
		authMiddleware: authMiddleware,
	}
}

func (h *BookingHandler) RegisterRoutes(r chi.Router) {
	r.Route(
		"/api/v1", func(r chi.Router) {
			// Маршруты для бронирования
			r.Route(
				"/bookings", func(r chi.Router) {
					// Публичные маршруты
					r.Group(
						func(r chi.Router) {
							r.Get("/available-rooms", h.GetAvailableRooms)
						},
					)
					// Защищенные маршруты
					r.Group(
						func(r chi.Router) {
							r.Use(h.authMiddleware.ValidateToken)
							r.Post("/", h.CreateBooking)
						},
					)
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

// @Summary Get available rooms for booking
// @Description Returns list of available rooms for specified dates and criteria
// @Tags bookings
// @Accept json
// @Produce json
// @Param checkIn query string true "Check-in date (YYYY-MM-DD)"
// @Param checkOut query string true "Check-out date (YYYY-MM-DD)"
// @Param capacity query integer false "Minimum room capacity"
// @Param type query string false "Room type (STANDARD, DELUXE, SUITE)"
// @Success 200 {array} response.AvailableRoom
// @Failure 400 {object} response.Error
// @Failure 500 {object} response.Error
// @Router /api/v1/bookings/available-rooms [get]
func (h *BookingHandler) GetAvailableRooms(w http.ResponseWriter, r *http.Request) {
	logger.Log.Info(
		"received request for available rooms",
		"path", r.URL.Path,
		"query", r.URL.RawQuery,
		"method", r.Method,
	)

	ctx := r.Context()

	// Парсим и валидируем параметры запроса
	params, err := h.parseSearchParams(r)
	if err != nil {
		if errors.IsInvalidInput(err) {
			logger.Log.Error("failed to parse search params", "error", err)
			h.respondWithError(w, http.StatusBadRequest, err)
		} else {
			h.respondWithError(
				w,
				http.StatusInternalServerError,
				errors.WithMessage(errors.ErrInternal, "failed to parse parameters"),
			)
		}
		return
	}

	// Формируем gRPC запрос
	req := &bookingpb.GetAvailableRoomsRequest{
		CheckIn:  mapper.TimeToProtoTimestamp(params.CheckIn),
		CheckOut: mapper.TimeToProtoTimestamp(params.CheckOut),
		Capacity: params.Capacity,
		Type:     params.Type,
	}

	// Устанавливаем timeout для запроса
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Делаем запрос к booking service
	resp, err := h.bookingClient.GetAvailableRooms(ctx, req)
	if err != nil {
		logger.Log.Error("failed to get available rooms", "error", err)
		h.respondWithError(
			w, http.StatusInternalServerError,
			errors.WithMessage(errors.ErrInternal, "failed to get available rooms"),
		)
		return
	}

	// Конвертируем ответ в HTTP формат
	rooms := mapper.ProtoToAvailableRooms(resp.Rooms)

	// Отправляем ответ
	h.respondWithJSON(w, http.StatusOK, rooms)
}

func (h *BookingHandler) parseSearchParams(r *http.Request) (*request.SearchParams, error) {
	checkIn, err := time.Parse("2006-01-02", r.URL.Query().Get("checkIn"))
	if err != nil {
		return nil, errors.WithMessage(errors.ErrInvalidInput, "invalid check-in date format")
	}

	checkOut, err := time.Parse("2006-01-02", r.URL.Query().Get("checkOut"))
	if err != nil {
		return nil, errors.WithMessage(errors.ErrInvalidInput, "invalid check-out date format")
	}

	params := &request.SearchParams{
		CheckIn:  checkIn,
		CheckOut: checkOut,
	}

	// Опциональные параметры
	if capacity := r.URL.Query().Get("capacity"); capacity != "" {
		var capVal int32
		if _, err = fmt.Sscan(capacity, &capVal); err != nil {
			return nil, errors.WithMessage(errors.ErrInvalidInput, "invalid capacity value")
		}
		params.Capacity = &capVal
	}

	if roomType := r.URL.Query().Get("type"); roomType != "" {
		// Проверяем, что значение есть в enum
		if val, ok := roompb.RoomType_value[roomType]; !ok {
			return nil, errors.WithMessage(errors.ErrInvalidInput, "invalid room type")
		} else {
			t := roompb.RoomType(val)
			params.Type = &t
		}
	}

	// Валидируем параметры
	if err = params.Validate(); err != nil {
		return nil, err
	}

	return params, nil
}

func (h *BookingHandler) respondWithJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Log.Error("failed to encode response", "error", err)
	}
}

func (h *BookingHandler) respondWithError(w http.ResponseWriter, code int, err error) {
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

// TODO: заглушка до реализации
func (h *BookingHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	// Получаем данные пользователя из контекста (установленные в AuthMiddleware)
	userInfo, ok := r.Context().Value(constants.USER).(*authpb.UserInfo)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Читаем тело запроса
	var req request.CreateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Создаем заглушку ответа
	res := response.CreateBookingResponse{
		ID:     uuid.New().String(), // Генерируем фейковый ID
		RoomID: req.RoomID,
		UserInfo: &response.UserInfo{
			ID:        userInfo.Id,
			Email:     userInfo.Email,
			FirstName: userInfo.FirstName,
			LastName:  userInfo.LastName,
			Role:      userInfo.Role.String(),
		},
		Status:  "PENDING",
		Message: "Booking created successfully (stub response)",
	}

	// Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		logger.Log.Error("failed to encode response", "error", err)
	}

	logger.Log.Info(
		"stub booking created",
		"user_id", userInfo.Id,
		"room_id", req.RoomID,
		"booking_id", res.ID,
	)
}
