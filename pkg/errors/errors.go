package errors

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound     = errors.New("entity not found")
	ErrConflict     = errors.New("entity already exists")
	ErrInvalidInput = errors.New("invalid input")
	ErrInternal     = errors.New("internal error")
)

// WithMessage оборачивает ошибку с дополнительным сообщением
func WithMessage(err error, message string) error {
	return fmt.Errorf("%s: %w", message, err)
}

// IsNotFound проверяет, является ли ошибка типом ErrNotFound
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// IsConflict проверяет, является ли ошибка типом ErrConflict
func IsConflict(err error) bool {
	return errors.Is(err, ErrConflict)
}

// IsInvalidInput проверяет, является ли ошибка типом ErrInvalidInput
func IsInvalidInput(err error) bool {
	return errors.Is(err, ErrInvalidInput)
}

// IsInternal проверяет, является ли ошибка типом ErrInternal
func IsInternal(err error) bool {
	return errors.Is(err, ErrInternal)
}
