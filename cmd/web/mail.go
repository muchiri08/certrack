package main

import (
	"errors"

	"github.com/muchiri08/certrack/internal/models"
)

func (app *application) sendNotification() {
	app.infoLog.Println("sending notification...")
	data, err := app.models.Certs.GetAlmostExpired()
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNoRecordFound):
			app.infoLog.Println("No certs yet to expire!")
			return
		default:
			app.errorLog.Println(err)
			return
		}
	}

	for _, d := range data {
		if err := app.mailer.Send(d.Email, d); err != nil {
			app.errorLog.Println(err)
		}
	}
}
