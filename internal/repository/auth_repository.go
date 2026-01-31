package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/PranavJoshi2893/med-portal/internal/model"
	"github.com/lib/pq"
)

type AuthRepository interface {
	Register(ctx context.Context, user model.User) error
	Login(ctx context.Context, email string) (*model.GetByEmail, error)
}

type AuthRepo struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepo {
	return &AuthRepo{
		db: db,
	}
}

func (r *AuthRepo) Register(ctx context.Context, user model.User) error {
	q := `INSERT INTO users(id,first_name,last_name,email,password) Values($1,$2,$3,$4,$5)`

	_, err := r.db.ExecContext(
		ctx,
		q,
		user.ID,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password,
	)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return fmt.Errorf("%w", model.ErrAlreadyExists)
		}
		return err
	}

	return nil
}

func (r *AuthRepo) Login(ctx context.Context, email string) (*model.GetByEmail, error) {
	q := `SELECT id, password FROM users WHERE email=$1`

	var user model.GetByEmail
	if err := r.db.QueryRowContext(ctx, q, email).Scan(
		&user.ID,
		&user.Password,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, model.ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}
