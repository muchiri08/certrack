package models

import (
	"database/sql"
	"errors"
)

var (
	ErrDuplicateEmail    = errors.New("duplicate email")
	ErrDuplicateUsername = errors.New("duplicate username")
	ErrNoRecordFound     = errors.New("no record found")
	ErrInvalidPassword   = errors.New("invalid credentials")
)

type Models struct {
	Users UserModel
	Certs CertModel
}

func NewModel(db *sql.DB) Models {
	return Models{
		Users: UserModel{db},
		Certs: CertModel{db},
	}
}
