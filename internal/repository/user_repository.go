package repository

import (
	"database/sql"

	"github.com/PranavJoshi2893/med-portal/internal/model"
	"github.com/google/uuid"
)

type UserRepository interface {
	GetAll() ([]model.GetAll, error)
	DeleteByID(id uuid.UUID) error
	GetByID(id uuid.UUID) (*model.GetByID, error)
}

type UserRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (r *UserRepo) GetAll() ([]model.GetAll, error) {

	q := `SELECT id, first_name, last_name, email FROM users`

	data, err := r.db.Query(q)
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

func (r *UserRepo) GetByID(id uuid.UUID) (*model.GetByID, error) {
	q := `SELECT id, first_name, last_name, email FROM users WHERE id = $1`

	var user model.GetByID
	if err := r.db.QueryRow(q, id).Scan(
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

func (r *UserRepo) DeleteByID(id uuid.UUID) error {

	q := `DELETE FROM users WHERE id = $1`

	res, err := r.db.Exec(q, id)

	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return model.ErrAlreadyDeleted
	}

	return nil
}
