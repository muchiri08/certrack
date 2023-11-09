package models

import (
	"database/sql"
	"log"
	"time"
)

type Certificate struct {
	ID         int32
	UserID     int32
	Domain     string
	Issuer     string
	ExpiryDate time.Time
	DaysLeft   int64
	Status     string
}

type CertModel struct {
	DB *sql.DB
}

func NewCertModel(db *sql.DB) *CertModel {
	return &CertModel{
		DB: db,
	}
}

func (m *CertModel) Inser(cert *Certificate) error {
	stmt := `INSERT INTO certs(user_id, domain, issuer, expiry_date, days_left, status VALUES($1, $2, $3, $4, $5))`

	args := []interface{}{cert.UserID, cert.Domain, cert.Issuer, cert.ExpiryDate, cert.DaysLeft, cert.Status}

	_, err := m.DB.Exec(stmt, args...)
	if err != nil {
		log.Println("Error:", err.Error())
		return err
	}

	return nil
}
