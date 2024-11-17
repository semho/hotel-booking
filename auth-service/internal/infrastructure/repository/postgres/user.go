package postgres

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/semho/hotel-booking/auth-service/internal/domain/model"
	"github.com/semho/hotel-booking/auth-service/internal/domain/port"
	"github.com/semho/hotel-booking/pkg/errors"
	"time"
)

const (
	tableUsers = "users"

	idColumn        = "id"
	emailColumn     = "email"
	passwordColumn  = "password"
	firstNameColumn = "first_name"
	lastNameColumn  = "last_name"
	phoneColumn     = "phone"
	roleColumn      = "role"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

type userRepository struct {
	db      *sqlx.DB
	builder squirrel.StatementBuilderType
}

func NewUserRepository(db *sqlx.DB) port.UserRepository {
	return &userRepository{
		db:      db,
		builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	user.ID = uuid.New()
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	sql, args, err := r.builder.
		Insert(tableUsers).
		Columns(
			idColumn,
			emailColumn,
			passwordColumn,
			firstNameColumn,
			lastNameColumn,
			phoneColumn,
			roleColumn,
			createdAtColumn,
			updatedAtColumn,
		).
		Values(
			user.ID,
			user.Email,
			user.Password,
			user.FirstName,
			user.LastName,
			user.Phone,
			user.Role,
			user.CreatedAt,
			user.UpdatedAt,
		).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	if _, err := r.db.ExecContext(ctx, sql, args...); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	sql, args, err := r.builder.
		Select(
			idColumn,
			emailColumn,
			passwordColumn,
			firstNameColumn,
			lastNameColumn,
			phoneColumn,
			roleColumn,
			createdAtColumn,
			updatedAtColumn,
		).
		From(tableUsers).
		Where(squirrel.Eq{emailColumn: email}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var user model.User
	err = r.db.GetContext(ctx, &user, sql, args...)
	if err != nil {
		return nil, errors.ErrNotFound
	}

	return &user, nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	sql, args, err := r.builder.
		Select(
			idColumn,
			emailColumn,
			passwordColumn,
			firstNameColumn,
			lastNameColumn,
			phoneColumn,
			roleColumn,
			createdAtColumn,
			updatedAtColumn,
		).
		From(tableUsers).
		Where(squirrel.Eq{idColumn: id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var user model.User
	err = r.db.GetContext(ctx, &user, sql, args...)
	if err != nil {
		return nil, errors.ErrNotFound
	}

	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	user.UpdatedAt = time.Now()

	sql, args, err := r.builder.
		Update(tableUsers).
		Set(emailColumn, user.Email).
		Set(passwordColumn, user.Password).
		Set(firstNameColumn, user.FirstName).
		Set(lastNameColumn, user.LastName).
		Set(phoneColumn, user.Phone).
		Set(roleColumn, user.Role).
		Set(updatedAtColumn, user.UpdatedAt).
		Where(squirrel.Eq{idColumn: user.ID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	result, err := r.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return errors.ErrNotFound
	}

	return nil
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	sql, args, err := r.builder.
		Delete(tableUsers).
		Where(squirrel.Eq{idColumn: id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	result, err := r.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return errors.ErrNotFound
	}

	return nil
}
