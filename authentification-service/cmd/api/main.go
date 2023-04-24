package main

import (
	"authentification/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "8083"

var counter int16

type Config struct {
	DB     *sql.DB
	models data.Models
}

func main() {
	log.Println("Start authentification service ...")

	conn := connectToDB()

	if conn == nil {
		log.Panic("Can't connect to DB!")
	}

	app := Config{
		DB:     conn,
		models: data.New(conn),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
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

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not ready yet ...")
			counter++
		} else {
			log.Println("Connected to Postgres!")
			return connection
		}

		if counter > 10 {
			log.Println(err)
			return nil
		}

		log.Panicln("Backing of for 2 seconds ...")
		time.Sleep(2 * time.Second)
		continue
	}
}
