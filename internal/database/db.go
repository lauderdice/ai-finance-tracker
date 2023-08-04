package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 6432
	user     = "lauderdice"
	password = "lauderdice"
	dbname   = "finance"
)

var DB *sql.DB

func InitDatabase() error {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}

	err = DB.Ping()
	if err != nil {
		return err
	}

	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL
		);

		CREATE TABLE IF NOT EXISTS payments (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users (id),
			value DECIMAL(10,2) NOT NULL,
			item VARCHAR(255) NOT NULL,
			category VARCHAR(255) NOT NULL,
			date DATE NOT NULL,
			notes TEXT
		);
	`)
	if err != nil {
		return err
	}

	return nil
}
