package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var schema = `
CREATE TABLE IF NOT EXISTS apps (
	id UUID PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS roles (
	id UUID PRIMARY KEY,
	app_id UUID NOT NULL REFERENCES apps,
	name VARCHAR(60) NOT NULL
);

CREATE TABLE IF NOT EXISTS permissions (
	id UUID PRIMARY KEY,
	app_id UUID NOT NULL REFERENCES apps,
	name VARCHAR(60) NOT NULL,
	entity_id UUID NOT NULL
);

CREATE TABLE IF NOT EXISTS role_permissions (
	id UUID PRIMARY KEY,
	app_id UUID NOT NULL REFERENCES apps,
	permission_id UUID NOT NULL,
	entity_id UUID NOT NULL
);

CREATE TABLE IF NOT EXISTS entity_roles (
	id UUID PRIMARY KEY,
	app_id UUID NOT NULL REFERENCES apps,
	entity_id UUID NOT NULL,
	role_id UUID NOT NULL
);
`

type DbConfig struct {
	Host 		string `json:"host"`
	User 		string `json:"username"`
	Password 	string `json:"password"`
	Database 	string `json:"database"`
	Port 		string `json:"port"`

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

	// Create schema
	_, err = db.Exec(schema)
	if err != nil {
		return nil, err
	}

	return db.Unsafe(), nil
}