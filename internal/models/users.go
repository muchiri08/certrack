package models

import (
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

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
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return ErrDuplicateUsername
		default:
			return err
		}
	}

	return nil
}

func (m *UserModel) GetUserByUsername(username string) (*User, error) {
	var user User
	stmt := `SELECT id, username, password_hash FROM users WHERE username = $1`

	err := m.DB.QueryRow(stmt, &username).Scan(&user.Id, &user.Username, &user.Password.Hash)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoRecordFound
		default:
			return nil, err
		}
	}

	return &user, nil

}

func (p *Password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.Plaintext = &plaintextPassword
	p.Hash = hash

	return nil
}

func (p *Password) Matches(plainPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.Hash, []byte(plainPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
