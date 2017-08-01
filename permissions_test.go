package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"log"
	"os"
	"testing"
)

var db *sqlx.DB

func TestMain(m *testing.M) {
	config := viper.New()
	config.SetConfigFile(os.Getenv("CONFIG"))
	err := config.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	var dbconfig DbConfig
	config.UnmarshalKey("database", &dbconfig)

	db, err = InitDB(dbconfig)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`
		drop table if exists apps cascade;
		drop table if exists permissions cascade;
		drop table if exists role_permissions cascade;
		drop table if exists roles cascade;
		drop table if exists entity_roles cascade;
	`)
	if err != nil {
		log.Fatal(err)
	}

	exitVal := m.Run()

	os.Exit(exitVal)
}

func TestCreatePermission(t *testing.T) {
	err := cleanup(db)
	if err != nil {
		log.Fatal(err)
	}
}

func cleanup(database *sqlx.DB) error {
	_, err := database.Exec(`
		drop table if exists apps cascade;
		drop table if exists permissions cascade;
		drop table if exists role_permissions cascade;
		drop table if exists roles cascade;
		drop table if exists entity_roles cascade;
	`)
	return err
}
