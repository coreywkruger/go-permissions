package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var schema = `
CREATE TABLE IF NOT EXISTS permissions (
	id uuid PRIMARY KEY,
	app_id uuid,
	name varchar(60),
	owner_id uuid
);

CREATE TABLE IF NOT EXISTS role_permissions (
	id uuid PRIMARY KEY,
	app_id uuid,
	permission_id uuid,
	owner_id uuid
);

CREATE TABLE IF NOT EXISTS entity_roles (
	id uuid PRIMARY KEY,
	app_id uuid,
	entity_id uuid,
	role_id uuid
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