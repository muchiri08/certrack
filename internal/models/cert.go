package models

import (
	"database/sql"
	"errors"
	"log"
	"time"
)

type Certificate struct {
	ID         int32
	UserID     int32
	Domain     string
	Issuer     string
	ExpiryDate time.Time
	DaysLeft   int
	Status     string
	CreatedAt  time.Time
}

type CertModel struct {
	DB *sql.DB
}

func NewCertModel(db *sql.DB) *CertModel {
	return &CertModel{
		DB: db,
	}
}

func (m *CertModel) Insert(cert *Certificate) error {
	stmt := `INSERT INTO certs(user_id, domain, issuer, expiry_date, days_left, status) VALUES($1, $2, $3, $4, $5, $6)`

	args := []interface{}{cert.UserID, cert.Domain, cert.Issuer, cert.ExpiryDate, cert.DaysLeft, cert.Status}

	_, err := m.DB.Exec(stmt, args...)
	if err != nil {
		log.Println("Error:", err.Error())
		return err
	}

	return nil
}

func (m *CertModel) GetCerts(userId int32) ([]*Certificate, error) {
	stmt := `SELECT * FROM certs WHERE user_id = $1`

	rows, err := m.DB.Query(stmt, userId)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoRecordFound
		default:
			return nil, err
		}

	}

	var certs []*Certificate
	for rows.Next() {
		var c Certificate

		err := rows.Scan(&c.ID, &c.Domain, &c.Issuer, &c.ExpiryDate, &c.DaysLeft, &c.Status, &c.CreatedAt, &c.UserID)
		if err != nil {
			return nil, err
		}

		certs = append(certs, &c)
	}

	return certs, nil

}

type Data struct {
	Username string
	Domain   string
	DaysLeft int
	Email    string
}

func (m *CertModel) GetAlmostExpired() ([]*Data, error) {
	stmt := `SELECT u.username, c.domain, c.days_left, u.email FROM certs c JOIN users u ON c.user_id = u.id WHERE c.days_left < 50`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoRecordFound
		default:
			return nil, err
		}

	}

	var dataSlice []*Data
	for rows.Next() {
		var data Data

		err := rows.Scan(&data.Username, &data.Domain, &data.DaysLeft, &data.Email)
		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				return nil, ErrNoRecordFound
			default:
				return nil, err
			}
		}

		dataSlice = append(dataSlice, &data)
	}

	return dataSlice, nil
}
