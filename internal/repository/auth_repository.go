package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/PranavJoshi2893/med-portal/internal/model"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type AuthRepository interface {
	Register(ctx context.Context, user model.User) error
	Login(ctx context.Context, email string) (*model.GetByEmail, error)
	StoreRefreshToken(ctx context.Context, id uuid.UUID, userID uuid.UUID, token string, expiresAt time.Time) error
	RevokeRefreshToken(ctx context.Context, token string) error
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
	q := `SELECT id, password, role FROM users WHERE email=$1 AND is_deleted = false`

	var user model.GetByEmail
	if err := r.db.QueryRowContext(ctx, q, email).Scan(
		&user.ID,
		&user.Password,
		&user.Role,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, model.ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *AuthRepo) StoreRefreshToken(ctx context.Context, id uuid.UUID, userID uuid.UUID, token string, expiresAt time.Time) error {
	q := `INSERT INTO refresh_tokens(id, user_id, token_hash, expires_at) VALUES($1, $2, $3, $4)`

	_, err := r.db.ExecContext(ctx, q, id, userID, token, expiresAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *AuthRepo) RevokeRefreshToken(ctx context.Context, token string) error {
	q := `UPDATE refresh_tokens SET revoked = true WHERE token_hash = $1 AND revoked = false`

	res, err := r.db.ExecContext(ctx, q, token)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return model.ErrNotFound
	}

	return nil
}
