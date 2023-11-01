package models

import (
	"database/sql"
	"errors"
	"time"
)

var ErrDuplicateEmail = errors.New("duplicate email")

type User struct {
	Id         int32
	Username   string
	Email      string
	Password   Password
	CreatedAt  time.Time
	Subscribed bool
}

type Password struct {
	Plaintext *string
	Hash      []byte
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) NewUser(user *User) error {
	stmt := `INSERT INTO users (username, email, password_hash, created_at, subscribed) VALUES ($1, $2, $3, CURRENT_TIMESTAMP, False)`

	args := []interface{}{user.Username, user.Email, user.Password.Hash}

	_, err := m.DB.Exec(stmt, args...)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}
