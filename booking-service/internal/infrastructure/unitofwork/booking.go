package unitofwork

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/semho/hotel-booking/booking-service/internal/domain/port"
)

type bookingUnitOfWork struct {
	db *sqlx.DB
}

func NewBookingUnitOfWork(db *sqlx.DB) port.BookingUnitOfWork {
	return &bookingUnitOfWork{db: db}
}

func (uow *bookingUnitOfWork) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := uow.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	txCtx := context.WithValue(ctx, "tx", tx)

	if err = fn(txCtx); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
