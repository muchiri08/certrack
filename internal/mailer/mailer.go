package mailer

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"

	"gopkg.in/gomail.v2"
)

//go:embed "template"
var templateFS embed.FS

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Sender   string
}

type Mailer struct {
	dialer *gomail.Dialer
	sender string
}

func New(config *SMTPConfig) Mailer {
	fmt.Printf("smtp configs: %v\n", *config)
	dialer := gomail.NewDialer(config.Host, config.Port, config.Username, config.Password)

	return Mailer{
		dialer: dialer,
		sender: config.Sender,
	}
}

func (m *Mailer) Send(to string, data interface{}) error {
	tmpl, err := template.New("email").ParseFS(templateFS, "template/template.email.html")
	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}

	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	msg := gomail.NewMessage()

	msg.SetHeader("From", m.sender)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject.String())
	msg.SetBody("text/plain", plainBody.String())
	msg.AddAlternative("text/html", htmlBody.String())

	return m.dialer.DialAndSend(msg)
}
