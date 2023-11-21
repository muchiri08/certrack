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
	dsn            string
	contextKeyUser = contextKey("user")
)

func main() {

	port := flag.Int("port", 4000, "port address")
	flag.StringVar(&dsn, "dsn", "", "PostgreSQL DSN")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	tempCache, err := newTemplateCache("./ui/html")
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

	config := mailer.SMTPConfig{
		Host:     "sandbox.smtp.mailtrap.io",
		Port:     2525,
		Username: "68d5bfa388790b",
		Password: "baefec63af537d",
		Sender:   "Mbogo <mbogokennedy08@gmail.com>",
	}

	mailer := mailer.New(&config)

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

	//test sending
	app.sendNotification()

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
