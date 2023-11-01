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

	_ "github.com/lib/pq"

	"github.com/muchiri08/certrack/internal/models"
)

type application struct {
	infoLog       *log.Logger
	errorLog      *log.Logger
	templateCache map[string]*template.Template
	models        models.Models
}

var dsn string

func main() {

	port := flag.Int("port", 4000, "port address")
	flag.StringVar(&dsn, "dsn", "postgres://root:KiNuThiaPro$2@localhost/certrack", "PostgreSQL DSN")

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

	app := &application{
		infoLog:       infoLog,
		errorLog:      errorLog,
		templateCache: tempCache,
		models:        models.NewModel(db),
	}

	srv := http.Server{
		Addr:         fmt.Sprintf(":%d", *port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

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
