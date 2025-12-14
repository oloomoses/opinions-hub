package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

func Connect() (*sql.DB, error) {

	dns := fmt.Sprintf(
		"postgres://%s:%s@localhost:5432/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	fmt.Println(dns)

	db, err := sql.Open("pgx", dns)

	if err != nil {
		return nil, err
	}

	maxConns := 20
	maxIdle := 10

	db.SetMaxOpenConns(maxConns)
	db.SetMaxIdleConns(maxIdle)
	db.SetConnMaxLifetime(30 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	DB = db

	return db, nil
}
