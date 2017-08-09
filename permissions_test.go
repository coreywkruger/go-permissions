package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestAllowed(t *testing.T) {
	var cases = []struct {
		EntityID     string
		AppID        string
		PermissionID string
		Expected     bool
		IsErr        bool
	}{
		{
			"c51003fc-2ae4-4296-9d5e-325c76a40316", "697d78cb-b56d-41ad-a7a3-e2e08ebb09fb", "5bee1c60-43e4-460e-80ae-b7c3b8774033", true, false,
		}, {
			"c51003fc-2ae4-4296-9d5e-325c76a40316", "697d78cb-b56d-41ad-a7a3-e2e08ebb09fb", "bad permission id", false, true,
		}, {
			"c1688c91-b818-4917-a20e-b95a2006c07f", "697d78cb-b56d-41ad-a7a3-e2e08ebb09fb", "28a212cc-51eb-4e17-95e1-2baa65e55b16", false, false,
		},
	}

	for _, tc := range cases {
		config := testConfig()
		db := testDb(config.GetString("database"))
		testCleanup(db)
		testMigrate(db)

		P := Permissionist{db}

		allowed, err := P.Allowed(tc.EntityID, tc.AppID, tc.PermissionID)
		if (err != nil) != tc.IsErr {
			log.Println(tc.Expected)
			t.Errorf("Unexpected error response [%v]", err)
		}
		if allowed != tc.Expected {
			t.Errorf("Expected permission to be '%t' got '%t'", tc.Expected, allowed)
		}
	}
}

func TestCreatePermission(t *testing.T) {
	var cases = []struct {
		AppID      string
		Permission string
		Expected   string
		IsErr      bool
	}{
		{
			"697d78cb-b56d-41ad-a7a3-e2e08ebb09fb", "read", "read", false,
		}, {
			"bad app id", "read", "", true,
		}, {
			"697d78cb-b56d-41ad-a7a3-e2e08ebb09fb", "", "", true,
		}, {
			"", "read", "", true,
		},
	}

	for _, tc := range cases {
		config := testConfig()
		db := testDb(config.GetString("database"))
		testCleanup(db)
		testMigrate(db)

		P := Permissionist{db}

		newPermission, err := P.CreatePermission(tc.Permission, tc.AppID)
		if (err != nil) != tc.IsErr {
			log.Println(tc.Expected)
			t.Errorf("Unexpected error response [%v]", err)
		}
		if newPermission.Name != tc.Expected {
			t.Errorf("Expected permission name of '%s' got '%s'", tc.Expected, newPermission.Name)
		}
	}
}

func TestCreatePermissions(t *testing.T) {
	var cases = []struct {
		AppID       string
		Permissions []string
		Expected    []string
		IsErr       bool
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
					found++
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
		AppID    string
		Roles    []string
		Expected []string
		IsErr    bool
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
					found++
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
		DROP TABLE IF EXISTS apps CASCADE;
		DROP TABLE IF EXISTS permissions CASCADE;
		DROP TABLE IF EXISTS role_permissions CASCADE;
		DROP TABLE IF EXISTS roles CASCADE;
		DROP TABLE IF EXISTS entity_roles CASCADE;
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
