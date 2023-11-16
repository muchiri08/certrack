package mailer

import "gopkg.in/gomail.v2"

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

type Mailer struct {
	dialer *gomail.Dialer
}

func New(config *SMTPConfig) Mailer {
	dialer := gomail.NewDialer(config.Host, config.Port, config.Username, config.Password)

	return Mailer{
		dialer: dialer,
	}
}

func (m *Mailer) Send(from, to, subject string) error {
	msg := gomail.NewMessage()

	msg.SetHeader("From", from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)

	return m.dialer.DialAndSend(msg)
}
