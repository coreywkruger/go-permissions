package main

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestCreatePermission(t *testing.T) {
	var testCases = []struct {
		n        int
		expected int
	}{
		{1, 1},
	}

	for _, test := range testCases {
		log.Println(test)

		config := testConfig()
		db := testDb(config.GetString("database"))
		testCleanup(db)

		P := Permissionist{db}

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
}

func TestCreateRoles(t *testing.T) {
	var testCases = []struct {
		args      []string
		expected  []string
		errorMessage interface{}
	}{
		{
			[]string{"one", "two", "three"},
			[]string{"one", "two", "three"},
			nil,
		}, {
			[]string{"one", "two", "three"},
			[]string{"one", "two", "three"},
			interface{}{},
		},
	}

	for _, testCase := range testCases {
		config := testConfig()
		db := testDb(config.GetString("database"))
		testCleanup(db)

		P := Permissionist{db}

		app, err := P.CreateApp("Taco-App")
		if err != nil {
			t.Errorf("Error: %s", err.Error())
		}

		err = P.CreateRoles(testCase.args, app.ID)
		if err != nil {
			t.Errorf("Error: %s", err.Error())
		}

		roles, err := P.GetRoles(app.ID)
		if err != nil {
			t.Errorf("Error: %s", err.Error())
		}

		for i := 0; i < len(roles); i++ {
			if roles[i].Name != testCase.expected[i] {
				t.Errorf("Error: Expected %s, Got %s", testCase.expected[i], roles[i])
			}
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
	schemas, err := ioutil.ReadFile("schemas.sql")
	if err != nil {
		log.Fatal(err)
	}
	// Create schema
	_, err = db.Exec(string(schemas))
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
