package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/carlescere/scheduler"
	_ "github.com/lib/pq"

	"github.com/muchiri08/certrack/internal/mailer"
	"github.com/muchiri08/certrack/internal/models"
)

type application struct {
	infoLog        *log.Logger
	errorLog       *log.Logger
	templateCache  map[string]*template.Template
	models         models.Models
	sessionManager *scs.SessionManager
	mailer         mailer.Mailer
}

type contextKey string

var (
	smtpConfig     mailer.SMTPConfig
	dsn            string
	contextKeyUser = contextKey("user")
)

func main() {

	port := flag.Int("port", 4000, "port address")
	flag.StringVar(&dsn, "dsn", "", "PostgreSQL DSN")
	flag.IntVar(&smtpConfig.Port, "smtp-port", 2525, "SMTP port")
	flag.StringVar(&smtpConfig.Host, "smtp-host", "sandbox.smtp.mailtrap.io", "SMTP host")
	flag.StringVar(&smtpConfig.Username, "smtp-user", "68d5bfa388790b", "SMTP username")
	flag.StringVar(&smtpConfig.Password, "smtp-pwd", "baefec63af537d", "SMTP password")
	flag.StringVar(&smtpConfig.Sender, "smtp-sender", "mbogokennedy08@gmail.com", "SMTP sender email")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	tempCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	db, err := openDb()
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	sessionManager := scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.SameSite = http.SameSiteStrictMode

	mailer := mailer.New(&smtpConfig)

	app := &application{
		infoLog:        infoLog,
		errorLog:       errorLog,
		templateCache:  tempCache,
		models:         models.NewModel(db),
		sessionManager: sessionManager,
		mailer:         mailer,
	}

	srv := http.Server{
		Addr:         fmt.Sprintf(":%d", *port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	//sending notifications
	scheduler.Every().Sunday().At("00:00").Run(app.sendNotification)

	infoLog.Println("Starting the server on port", *port)
	err = srv.ListenAndServe()
	if err != nil {
		errorLog.Fatal(err)
	}

}

func openDb() (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
