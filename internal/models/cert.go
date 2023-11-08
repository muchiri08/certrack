package models

import "time"

type Certificate struct {
	ID         int32
	UserID     int32
	Domain     string
	Issuer     string
	ExpiryDate time.Time
	DaysLeft   int64
	Status     string
}
