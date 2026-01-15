package repository

import (
	"database/sql"

	"github.com/PranavJoshi2893/med-portal/internal/model"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (r *UserRepo) Register(user model.User) error {
	q := `INSERT INTO users(id,first_name,last_name,email,password) Values($1,$2,$3,$4,$5)`

	_, err := r.db.Exec(
		q,
		user.ID,
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

func (r *UserRepo) GetAll() ([]model.GetAll, error) {

	q := `SELECT first_name, last_name,email FROM users`

	data, err := r.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer data.Close()

	var users []model.GetAll

	for data.Next() {
		var user model.GetAll

		if err := data.Scan(
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
