package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestGrantPermissionToRole(t *testing.T) {
	var cases = []struct {
		RoleID       string
		PermissionID string
		IsErr        bool
	}{
		{
			"c1688c91-b818-4917-a20e-b95a2006c07f", "73017965-b16c-4c6e-9ec1-1e1272594648", false, // Role 'Admin' has a permission
		},
	}

	for _, tc := range cases {
		config := testConfig()
		db := testDb(config.GetString("database"))
		testCleanup(db)
		testMigrate(db)

		P := Permissionist{db}

		err := P.GrantPermissionToRole(tc.RoleID, tc.PermissionID)
		if (err != nil) != tc.IsErr {
			t.Errorf("Unexpected error response [%v]", err)
		}
	}
}

func TestEntityIsAllowed(t *testing.T) {
	var cases = []struct {
		EntityID     string
		PermissionID string
		Expected     bool
		IsErr        bool
	}{
		{
			"809e5e2f-0555-4d81-8f91-d6d8f0d4ea79", "5bee1c60-43e4-460e-80ae-b7c3b8774033", true, false, // Role 'Admin' has a permission
		}, {
			"07df4a77-6243-41cd-a421-90c524ef2203", "28a212cc-51eb-4e17-95e1-2baa65e55b16", false, false, // Role 'Customer' doens't have permission
		}, {
			"07df4a77-6243-41cd-a421-90c524ef2203", "5bee1c60-43e4-460e-80ae-b7c3b8774033", true, false, // Role 'Customer' has permission
		}, {
			"809e5e2f-0555-4d81-8f91-d6d8f0d4ea79", "bad permission id", false, true, // Error caused by bad permission id
		}, {
			"809e5e2f-0555-4d81-8f91-d6d8f0d4ea79", "bad permission id", false, true, // Error caused by bad role id
		},
	}

	for _, tc := range cases {
		config := testConfig()
		db := testDb(config.GetString("database"))
		testCleanup(db)
		testMigrate(db)

		P := Permissionist{db}

		allowed, err := P.EntityIsAllowed(tc.EntityID, tc.PermissionID)
		if (err != nil) != tc.IsErr {
			t.Errorf("Unexpected error response [%v]", err)
		}
		if allowed != tc.Expected {
			t.Errorf("Expected permission to be '%t' got '%t'", tc.Expected, allowed)
		}
	}
}

func TestRoleIsAllowed(t *testing.T) {
	var cases = []struct {
		RoleID       string
		PermissionID string
		Expected     bool
		IsErr        bool
		Description  string
	}{
		{
			"c51003fc-2ae4-4296-9d5e-325c76a40316", "5bee1c60-43e4-460e-80ae-b7c3b8774033", true, false, "Role 'Admin' has a permission",
		}, {
			"c1688c91-b818-4917-a20e-b95a2006c07f", "28a212cc-51eb-4e17-95e1-2baa65e55b16", false, false, "Role 'Customer' doens't have permission",
		}, {
			"c1688c91-b818-4917-a20e-b95a2006c07f", "5bee1c60-43e4-460e-80ae-b7c3b8774033", true, false, "Role 'Customer' has permission",
		}, {
			"c51003fc-2ae4-4296-9d5e-325c76a40316", "bad permission id", false, true, "Error caused by malformed permission id",
		}, {
			"bad role id", "5bee1c60-43e4-460e-80ae-b7c3b8774033", false, true, "Error caused by malformed role id",
		},
	}

	for _, tc := range cases {
		config := testConfig()
		db := testDb(config.GetString("database"))
		testCleanup(db)
		testMigrate(db)

		P := Permissionist{db}

		allowed, err := P.RoleIsAllowed(tc.RoleID, tc.PermissionID)
		if (err != nil) != tc.IsErr {
			log.Println(tc.Expected)
			t.Errorf("Unexpected error response [%v]", err)
		}
		if allowed != tc.Expected {
			t.Errorf("Expected permission to be '%t' got '%t'", tc.Expected, allowed)
		}
	}
}

func TestGetApps(t *testing.T) {
	var cases = []struct {
		Expected []string
		IsErr    bool
	}{
		{
			[]string{"697d78cb-b56d-41ad-a7a3-e2e08ebb09fb"}, false,
		},
	}

	for _, tc := range cases {
		config := testConfig()
		db := testDb(config.GetString("database"))
		testCleanup(db)
		testMigrate(db)

		P := Permissionist{db}

		apps, err := P.GetApps()
		if (err != nil) != tc.IsErr {
			t.Errorf("Unexpected error response [%v]", err)
		}
		found := 0
		for i := range apps {
			for j := range tc.Expected {
				if apps[i].ID == tc.Expected[j] {
					found++
				}
			}
		}
		if found != len(tc.Expected) {
			t.Errorf("Expected %d roles got %d", len(tc.Expected), found)
		}
	}
}

func TestGetPermissionsByEntityID(t *testing.T) {
	var cases = []struct {
		EntityID string
		AppID    string
		Expected []string
		IsErr    bool
	}{
		{
			"809e5e2f-0555-4d81-8f91-d6d8f0d4ea79", "697d78cb-b56d-41ad-a7a3-e2e08ebb09fb", []string{
				"5bee1c60-43e4-460e-80ae-b7c3b8774033",
				"73017965-b16c-4c6e-9ec1-1e1272594648",
				"28a212cc-51eb-4e17-95e1-2baa65e55b16",
			}, false,
		},
	}

	for _, tc := range cases {
		config := testConfig()
		db := testDb(config.GetString("database"))
		testCleanup(db)
		testMigrate(db)

		P := Permissionist{db}

		permissions, err := P.GetPermissionsByEntityID(tc.EntityID, tc.AppID)
		if (err != nil) != tc.IsErr {
			t.Errorf("Unexpected error response [%v]", err)
		}
		found := 0
		for i := range permissions {
			for j := range tc.Expected {
				if permissions[i].ID == tc.Expected[j] {
					found++
				}
			}
		}
		if found != len(tc.Expected) {
			t.Errorf("Expected %d permissions got %d", len(tc.Expected), found)
		}
	}
}

func TestGetRolesByEntityID(t *testing.T) {
	var cases = []struct {
		EntityID string
		Expected []string
		IsErr    bool
	}{
		{
			"809e5e2f-0555-4d81-8f91-d6d8f0d4ea79", []string{"c51003fc-2ae4-4296-9d5e-325c76a40316"}, false,
		},
	}

	for _, tc := range cases {
		config := testConfig()
		db := testDb(config.GetString("database"))
		testCleanup(db)
		testMigrate(db)

		P := Permissionist{db}

		apps, err := P.GetRolesByEntityID(tc.EntityID)
		if (err != nil) != tc.IsErr {
			t.Errorf("Unexpected error response [%v]", err)
		}
		found := 0
		for i := range apps {
			for j := range tc.Expected {
				if apps[i].ID == tc.Expected[j] {
					found++
				}
			}
		}
		if found != len(tc.Expected) {
			t.Errorf("Expected %d roles got %d", len(tc.Expected), found)
		}
	}
}

func TestGetAppsByEntityID(t *testing.T) {
	var cases = []struct {
		EntityID string
		Expected []string
		IsErr    bool
	}{
		{
			"809e5e2f-0555-4d81-8f91-d6d8f0d4ea79", []string{"697d78cb-b56d-41ad-a7a3-e2e08ebb09fb"}, false,
		},
	}

	for _, tc := range cases {
		config := testConfig()
		db := testDb(config.GetString("database"))
		testCleanup(db)
		testMigrate(db)

		P := Permissionist{db}

		apps, err := P.GetAppsByEntityID(tc.EntityID)
		if (err != nil) != tc.IsErr {
			t.Errorf("Unexpected error response [%v]", err)
		}
		found := 0
		for i := range apps {
			for j := range tc.Expected {
				if apps[i].ID == tc.Expected[j] {
					found++
				}
			}
		}
		if found != len(tc.Expected) {
			t.Errorf("Expected %d roles got %d", len(tc.Expected), found)
		}
	}
}

func TestGetPermissionsByRoleID(t *testing.T) {
	var cases = []struct {
		RoleID   string
		Expected []string
		IsErr    bool
	}{
		{
			"c51003fc-2ae4-4296-9d5e-325c76a40316", []string{"5bee1c60-43e4-460e-80ae-b7c3b8774033", "73017965-b16c-4c6e-9ec1-1e1272594648", "28a212cc-51eb-4e17-95e1-2baa65e55b16"}, false,
		},
		{
			"c1688c91-b818-4917-a20e-b95a2006c07f", []string{"5bee1c60-43e4-460e-80ae-b7c3b8774033"}, false,
		},
	}

	for _, tc := range cases {
		config := testConfig()
		db := testDb(config.GetString("database"))
		testCleanup(db)
		testMigrate(db)

		P := Permissionist{db}

		permissions, err := P.GetPermissionsByRoleID(tc.RoleID)
		if (err != nil) != tc.IsErr {
			t.Errorf("Unexpected error response [%v]", err)
		}
		found := 0
		for i := range permissions {
			for j := range tc.Expected {
				if permissions[i].ID == tc.Expected[j] {
					found++
				}
			}
		}
		if found != len(tc.Expected) {
			t.Errorf("Expected %d roles got %d", len(tc.Expected), found)
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
			"697d78cb-b56d-41ad-a7a3-e2e08ebb09fb", "taco-eating", "taco-eating", false,
		}, {
			"bad-app-id", "taco-eating", "", true,
		}, {
			"697d78cb-b56d-41ad-a7a3-e2e08ebb09fb", "", "", true,
		}, {
			"", "taco-eating", "", true,
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
		}, {
			"697d78cb-b56d-41ad-a7a3-e2e08ebb09fb", []string{"read", "write"}, []string{}, true,
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

func TestAssignRoleToEntity(t *testing.T) {
	var cases = []struct {
		EntityID string
		RoleID   string
		IsErr    bool
	}{
		{
			"some entity", "c51003fc-2ae4-4296-9d5e-325c76a40316", false,
		}, {
			"some entity", "bad role id", true,
		},
	}

	for _, tc := range cases {
		config := testConfig()
		db := testDb(config.GetString("database"))
		testCleanup(db)
		testMigrate(db)

		P := Permissionist{db}

		err := P.AssignRoleToEntity(tc.EntityID, tc.RoleID)
		if (err != nil) != tc.IsErr {
			t.Errorf("Unexpected error response [%v]", err)
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
		}, {
			"697d78cb-b56d-41ad-a7a3-e2e08ebb09fb", []string{"admin", "customer"}, []string{}, true,
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
