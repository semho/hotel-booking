package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/semho/hotel-booking/api-gateway/internal/api/http/mapper"
	"github.com/semho/hotel-booking/api-gateway/internal/api/http/request"
	"github.com/semho/hotel-booking/api-gateway/internal/api/http/response"
	"github.com/semho/hotel-booking/pkg/errors"
	"github.com/semho/hotel-booking/pkg/logger"
	bookingpb "github.com/semho/hotel-booking/pkg/proto/booking_v1/booking"
	roompb "github.com/semho/hotel-booking/pkg/proto/room_v1/room"
	"net/http"
	"time"
)

type BookingHandler struct {
	bookingClient bookingpb.BookingServiceClient
}

func NewBookingHandler(bookingClient bookingpb.BookingServiceClient) *BookingHandler {
	return &BookingHandler{
		bookingClient: bookingClient,
	}
}

func (h *BookingHandler) RegisterRoutes(r chi.Router) {
	r.Route(
		"/api/v1", func(r chi.Router) {
			// Маршруты для бронирования
			r.Route(
				"/bookings", func(r chi.Router) {
					r.Get("/available-rooms", h.GetAvailableRooms)
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
		t := mapper.StringToProtoRoomType(roomType)
		if t == roompb.RoomType_ROOM_TYPE_UNSPECIFIED {
			return nil, errors.WithMessage(errors.ErrInvalidInput, "invalid room type")
		}
		params.Type = &t
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
