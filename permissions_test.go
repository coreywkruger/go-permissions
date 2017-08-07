package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestCreatePermissions(t *testing.T) {
	var cases = []struct {
		AppID 			string
		Permissions []string
		Expected 		[]string
		IsErr 			bool
	}{
		{
			"697d78cb-b56d-41ad-a7a3-e2e08ebb09fb", []string{"one", "two"}, []string{"one", "two"}, false,
		}, {
			"bad app id", []string{"one", "two"}, []string{}, true,
		},
	}

	for _, tc := range cases {
		config := testConfig()
		db := testDb(config.GetString("database"))
		testCleanup(db)
		testMigrate(db)

		P := Permissionist{db}

		newPermissions, err := P.CreatePermissions(tc.Permissions, tc.AppID)
		if (err != nil) != tc.IsErr {
			t.Errorf("Unexpected error response [%v]", err)
		}
		found := 0
		for i := range newPermissions {
			for j := range tc.Permissions {
				if newPermissions[i].Name == tc.Permissions[j] {
					found += 1
				}
			}
		}
		if found != len(tc.Expected) {
			t.Errorf("Expected %d roles got %d", len(tc.Permissions), found)
		}
	}
}

func TestCreateRoles(t *testing.T) {
	var cases = []struct {
		AppID 		string
		Roles     []string
		Expected 	[]string
		IsErr 		bool
	}{
		{
			"697d78cb-b56d-41ad-a7a3-e2e08ebb09fb", []string{"one", "two"}, []string{"one", "two"}, false,
		}, {
			"bad app id", []string{"one", "two"}, []string{}, true,
		},
	}

	for _, tc := range cases {
		config := testConfig()
		db := testDb(config.GetString("database"))
		testCleanup(db)
		testMigrate(db)

		P := Permissionist{db}

		newRoles, err := P.CreateRoles(tc.Roles, tc.AppID)
		if (err != nil) != tc.IsErr {
			t.Errorf("Unexpected error response [%v]", err)
		}
		found := 0
		for i := range newRoles {
			for j := range tc.Roles {
				if newRoles[i].Name == tc.Roles[j] {
					found += 1
				}
			}
		}
		if found != len(tc.Expected) {
			t.Errorf("Expected %d roles got %d", len(tc.Roles), found)
		}
	}
}

func testCleanup(db *sqlx.DB) {
	_, err := db.Exec(`
		drop table if exists apps cascade;
		drop table if exists permissions cascade;
		drop table if exists role_permissions cascade;
		drop table if exists roles cascade;
		drop table if exists entity_roles cascade;
	`)
	if err != nil {
		log.Fatal(err)
	}
}

func testMigrate(db *sqlx.DB) {
	schemas, err := ioutil.ReadFile("migrate.sql")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(string(schemas))
	if err != nil {
		log.Fatal(err)
	}
	migrate, err := ioutil.ReadFile("seed.sql")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(string(migrate))
	if err != nil {
		log.Fatal(err)
	}
}

func testConfig() *viper.Viper {
	config := viper.New()
	config.SetConfigFile(os.Getenv("CONFIG"))
	err := config.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	return config
}

func testDb(database string) *sqlx.DB {
	db, err := sqlx.Connect("postgres", database)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
