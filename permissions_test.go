package main

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	config := viper.New()
	config.SetConfigFile(os.Getenv("CONFIG"))
	err := config.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	var dbconfig DbConfig
	config.UnmarshalKey("database", &dbconfig)

	db, err := InitDB(dbconfig)
	if err != nil {
		log.Fatal(err)
	}
	_, _ = db.Exec(`
		drop table if exists apps cascade;
		drop table if exists permissions cascade;
		drop table if exists role_permissions cascade;
		drop table if exists roles cascade;
		drop table if exists entity_roles cascade;
	`)

	exitVal := m.Run()

	os.Exit(exitVal)
}
