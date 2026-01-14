package repository

import (
	"database/sql"

	"github.com/PranavJoshi2893/med-portal/internal/model"
)

type UserRepo struct {
	db *sql.DB
}

func (r *UserRepo) Register(user *model.CreateUser) error {
	q := `INSERT INTO users(id,first_name,last_name,email,password) Values($1,$2,$3,$4)`

	_, err := r.db.Exec(
		q,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepo) Login(user *model.LoginUser) (string, error) {

	q := `SELECT password FROM users WHERE email = $1`

	var hashedPassword string
	err := r.db.QueryRow(q, user.Email).Scan(&hashedPassword)
	if err != nil {
		return "", err
	}

	return hashedPassword, nil
}
