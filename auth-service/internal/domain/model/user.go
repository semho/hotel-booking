package model

import (
	"github.com/google/uuid"
	pb "github.com/semho/hotel-booking/pkg/proto/auth_v1/auth"
	"time"
)

// type UserRole string
type UserRole = pb.UserRole

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
