package main

import (
	"crypto/tls"
	"fmt"
	"math"
	"time"

	"github.com/muchiri08/certrack/internal/models"
)

func (app *application) track(userId int32, domains ...string) ([]*models.Certificate, error) {
	certs := []*models.Certificate{}
	for _, domain := range domains {
		conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", domain, 443), nil)
		if err != nil {
			return nil, err
		}
		defer conn.Close()

		cert := conn.ConnectionState().PeerCertificates[0]
		expirationDate := cert.NotAfter
		remDuration := time.Until(expirationDate)
		daysLeft := int(math.Round(remDuration.Hours() / 24))
		c := models.Certificate{
			UserID:     userId,
			Domain:     domain,
			Issuer:     cert.Issuer.CommonName,
			ExpiryDate: cert.NotAfter,
			DaysLeft:   daysLeft,
			Status:     "HEALTHY",
		}

		certs = append(certs, &c)

	}

	return certs, nil
}
