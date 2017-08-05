package main

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"log"
	"os"
)

func InitConfig() *viper.Viper {
	config := viper.New()
	config.SetConfigFile(os.Getenv("CONFIG"))
	err := config.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	return config
}

func InitDb(database string) *sqlx.DB {
	db, err := sqlx.Connect("postgres", database)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
