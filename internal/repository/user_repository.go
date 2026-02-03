package repository

import (
	"context"
	"database/sql"

	"github.com/PranavJoshi2893/med-portal/internal/model"
	"github.com/google/uuid"
)

type UserRepository interface {
	GetAll(ctx context.Context) ([]model.GetAll, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.GetByID, error)
	UpdateByID(ctx context.Context, id uuid.UUID, data *model.UpdateUser) error
}

type UserRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (r *UserRepo) GetAll(ctx context.Context) ([]model.GetAll, error) {

	q := `SELECT id, first_name, last_name, email FROM users WHERE is_deleted = false`

	data, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer data.Close()

	var users []model.GetAll

	for data.Next() {
		var user model.GetAll

		if err := data.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Email,
		); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := data.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.GetByID, error) {
	q := `SELECT id, first_name, last_name, email FROM users WHERE id = $1 AND is_deleted = false`

	var user model.GetByID
	if err := r.db.QueryRowContext(ctx, q, id).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, model.ErrNotFound
		}
		return nil, err
	}

	return &user, nil

}

func (r *UserRepo) DeleteByID(ctx context.Context, id uuid.UUID) error {

	q := `SELECT is_deleted FROM users WHERE id = $1`

	var isDeleted bool
	if err := r.db.QueryRowContext(ctx, q, id).Scan(&isDeleted); err != nil {
		if err == sql.ErrNoRows {
			return model.ErrNotFound
		}
		return err
	}

	if isDeleted {
		return model.ErrAlreadyDeleted
	}

	q = `UPDATE users SET is_deleted = true WHERE id = $1`

	if _, err := r.db.ExecContext(ctx, q, id); err != nil {
		return err
	}

	return nil
}

func (r *UserRepo) UpdateByID(ctx context.Context, id uuid.UUID, data *model.UpdateUser) error {
	if data == nil || data.FirstName == nil || data.LastName == nil {
		return model.ErrBadRequest
	}

	q := `SELECT is_deleted FROM users WHERE id=$1`
	var isDeleted bool
	if err := r.db.QueryRowContext(ctx, q, id).Scan(&isDeleted); err != nil {
		if err == sql.ErrNoRows {
			return model.ErrNotFound
		}
		return err
	}

	if isDeleted {
		return model.ErrAlreadyDeleted
	}

	q = `UPDATE users SET first_name=$1, last_name=$2 WHERE id = $3 AND is_deleted=false`
	if _, err := r.db.ExecContext(ctx, q, *data.FirstName, *data.LastName, id); err != nil {
		return err
	}

	return nil
}
