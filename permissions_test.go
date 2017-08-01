package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

var DB *sqlx.DB
var P Permissionist

func TestMain(m *testing.M) {
	config := viper.New()
	config.SetConfigFile(os.Getenv("CONFIG"))
	err := config.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	var dbConfig DbConfig
	config.UnmarshalKey("database", &dbConfig)

	DB, err = InitDB(dbConfig)
	if err != nil {
		log.Fatal(err)
	}

	P = Permissionist{
		DB: DB,
	}

	exitVal := m.Run()

	os.Exit(exitVal)
}

func TestCreatePermission(t *testing.T) {
	err := cleanup(DB)
	if err != nil {
		log.Fatal(err)
	}

	app, err := P.CreateApp("Taco App")
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
	role, err := P.CreateRole("Taco Role", app.ID)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
	permission, err := P.CreatePermission("Taco Permission", app.ID)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
	err = P.GrantPermissionToRole(role.ID, app.ID, permission.ID)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
}

func cleanup(db *sqlx.DB) error {
	_, err := db.Exec(`
		drop table if exists apps cascade;
		drop table if exists permissions cascade;
		drop table if exists role_permissions cascade;
		drop table if exists roles cascade;
		drop table if exists entity_roles cascade;
	`)

	schemas, err := ioutil.ReadFile("schemas.sql")
	if err != nil {
		return err
	}

	// Create schema
	_, err = db.Exec(string(schemas))
	if err != nil {
		return err
	}

	return err
}
