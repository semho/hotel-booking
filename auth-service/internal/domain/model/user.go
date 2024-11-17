package model

import (
	"github.com/google/uuid"
	"time"
)

type UserRole string

const (
	UserRoleUser  UserRole = "USER"
	UserRoleAdmin UserRole = "ADMIN"
)

type User struct {
	ID        uuid.UUID `db:"id"`
	Email     string    `db:"email"`
	Password  string    `db:"password"` // Хэшированный пароль
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	Phone     *string   `db:"phone"`
	Role      UserRole  `db:"role"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
