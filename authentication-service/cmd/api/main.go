package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dineshkuncham/authentication-service/data"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type Config struct {
	Repo   data.Repository
	Client *http.Client
}

const webport = "80"

var retries int64

func main() {
	log.Println("Starting authentication service")

	conn := connectToDB()
	if conn == nil {
		log.Panic("Can't connect to Postgres!")
	}

	api := Config{
		Client: &http.Client{},
	}

	api.setupRepo(conn)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webport),
		Handler: api.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")
	if dsn == "" {
		log.Panic("DSN env variable is not found")
	}
	for {
		conn, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not ready yet.....")
			retries++
		} else {
			log.Println("Connected to Postges!")
			return conn
		}
		if retries > 10 {
			log.Println(err)
			return nil
		}
		log.Println("Backing off for 2 seconds.....")
		time.Sleep(2 * time.Second)
		continue
	}

}

func (app *Config) setupRepo(conn *sql.DB) {
	db := data.NewPostgresRepository(conn)
	app.Repo = db
}
