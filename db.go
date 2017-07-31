package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"io/ioutil"
)

// DbConfig db info
type DbConfig struct {
	Host     string `json:"host"`
	User     string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
	Port     string `json:"port"`
}

// InitDB initializes the database
func InitDB(config DbConfig) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", fmt.Sprint(
		"host=", config.Host,
		" user=", config.User,
		" password=", config.Password,
		" dbname=", config.Database,
		" sslmode=disable",
		" port=", config.Port,
	))
	if err != nil {
		return nil, err
	}

	schemas, err := ioutil.ReadFile("schemas.sql")
	if err != nil {
		return nil, err
	}

	// Create schema
	_, err = db.Exec(string(schemas))
	if err != nil {
		return nil, err
	}

	return db.Unsafe(), nil
}
