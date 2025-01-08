-- +goose Up
-- +goose StatementBegin
-- Основная таблица броней
CREATE TABLE IF NOT EXISTS bookings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id UUID NOT NULL,
    user_id UUID,
    guest_name VARCHAR(255) NOT NULL,
    guest_email VARCHAR(255) NOT NULL,
    guest_phone VARCHAR(50),
    check_in TIMESTAMP WITH TIME ZONE NOT NULL,
    check_out TIMESTAMP WITH TIME ZONE NOT NULL,
    total_price DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT bookings_check_dates CHECK (check_out > check_in),
    CONSTRAINT bookings_check_price CHECK (total_price >= 0)
    );

-- История статусов брони
CREATE TABLE IF NOT EXISTS booking_status_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id UUID NOT NULL,
    status INTEGER NOT NULL,
    reason TEXT,
    changed_by VARCHAR(255) NOT NULL,
    changed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (booking_id) REFERENCES bookings(id) ON DELETE CASCADE,
    CONSTRAINT booking_status_history_status_check CHECK (status >= 0)
    );

-- Индексы для основной таблицы
CREATE INDEX idx_bookings_dates ON bookings (check_in, check_out);
CREATE INDEX idx_bookings_room_dates ON bookings (room_id, check_in, check_out);
CREATE INDEX idx_bookings_user ON bookings (user_id);

-- Индексы для истории статусов
CREATE INDEX idx_booking_status_history_booking ON booking_status_history (booking_id);
CREATE INDEX idx_booking_status_history_latest ON booking_status_history (booking_id, changed_at DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS booking_status_history;
DROP TABLE IF EXISTS bookings;
-- +goose StatementEnd