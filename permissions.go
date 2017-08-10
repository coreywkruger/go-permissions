package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"strings"
)

// App app schema
type App struct {
	ID   string `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

// RolePermission role_permissions schema
type RolePermission struct {
	ID           string `json:"id" db:"id"`
	RoledID      string `json:"role_id" db:"role_id"`
	PermissionID string `json:"permission_id" db:"permission_id"`
}

// Permission permissions schema
type Permission struct {
	ID    string `json:"id" db:"id"`
	Name  string `json:"name" db:"name"`
	AppID string `json:"app_id" db:"app_id"`
}

// Role roles schema
type Role struct {
	ID    string `json:"id" db:"id"`
	Name  string `json:"name" db:"name"`
	AppID string `json:"app_id" db:"app_id"`
}

// Permissionist owns permissions crud
type Permissionist struct {
	DB *sqlx.DB
}

// Allowed checks if entity entityID has permission permissionID
func (permissions *Permissionist) Allowed(roleID string, appID string, permissionID string) (bool, error) {
	var rolePermissions []string
	err := permissions.DB.Select(&rolePermissions, `
	SELECT r.id
	FROM roles AS r
	INNER JOIN role_permissions AS rp
		ON rp.permission_id = $3
			AND r.id = $1
			AND r.app_id = $2
			AND rp.role_id = r.id
	INNER JOIN permissions AS p
		ON p.id = $3
			AND rp.permission_id = p.id;
	`, roleID, appID, permissionID)

	if err != nil {
		return false, errors.Wrap(err, "Could not check permission")
	}

	if len(rolePermissions) <= 0 {
		return false, nil
	}

	return true, nil
}

// GetApps returns a list of all apps
func (permissions *Permissionist) GetApps() ([]App, error) {
	var apps []App
	err := permissions.DB.Select(&apps, `SELECT * FROM apps;`)
	if err != nil {
		return nil, errors.Wrap(err, "Could not get apps")
	}

	return apps, nil
}

// GetAppsByEntityID returns a list of all apps
func (permissions *Permissionist) GetAppsByEntityID(entityID string) ([]App, error) {
	var apps []App
	err := permissions.DB.Select(&apps, `
	SELECT a.id, a.name
	FROM apps AS a
	INNER JOIN roles AS r
		ON a.id = r.app_id
			AND r.id = $1;
	`, entityID)

	if err != nil {
		return nil, errors.Wrap(err, "Could not get apps")
	}

	return apps, nil
}

// GetApp returns an app by id
func (permissions *Permissionist) GetApp(appID string) (string, error) {
	var id string
	err := permissions.DB.Select(&id, `
	SELECT id
	FROM apps
	WHERE id = $1;
	`, appID)
	if err != nil {
		return "", errors.Wrap(err, "Could not get app")
	}

	return id, nil
}

// GetPermissionsByEntityID returns a list of all permissions that belong to an entity
func (permissions *Permissionist) GetPermissionsByEntityID(entityID string, appID string) ([]Permission, error) {
	var perms []Permission
	err := permissions.DB.Select(&perms, `
	SELECT *
	FROM permissions AS p
	INNER JOIN role_permissions AS rp
		ON p.id = rp.permission_id
			AND p.app_id = $2
			AND rp.app_id = $2
	INNER JOIN entity_roles AS er
		ON er.entity_id = $1
			AND er.app_id = $2;
	`, entityID, appID)

	if err != nil {
		return nil, errors.Wrap(err, "Could not get permissions")
	}

	return perms, nil
}

// GetPermissionsByRole returns a list of all permissions that belong to an entity
func (permissions *Permissionist) GetPermissionsByRole(roleID string, appID string) ([]Permission, error) {
	var perms []Permission
	err := permissions.DB.Select(&perms, `
	SELECT *
	FROM permissions AS p
	INNER JOIN role_permissions AS rp
		ON p.id = rp.permission_id
			AND rp.role_id = $1
			AND p.app_id = $2
			AND rp.app_id = $2
	`, roleID, appID)

	if err != nil {
		return nil, errors.Wrap(err, "Could not get permissions")
	}

	return perms, nil
}

// GetRoles returns a list of all roles created for an app
func (permissions *Permissionist) GetRoles(appID string) ([]Role, error) {
	roles := []Role{}
	err := permissions.DB.Select(&roles, `
	SELECT *
	FROM roles
	WHERE app_id = '`+appID+`';
	`)

	if err != nil {
		return nil, errors.Wrap(err, "Could not get roles")
	}

	return roles, nil
}

// GetRoleByID returns a role name
func (permissions *Permissionist) GetRoleByID(roleID string, appID string) (Role, error) {
	var role Role
	err := permissions.DB.Select(&role, `
	SELECT *
	FROM roles
	WHERE id = $1
		AND app_id = $2 limit 1;
	`, roleID, appID)

	if err != nil {
		return role, errors.Wrap(err, "Could not get role")
	}

	return role, nil
}

// AssignRoleToEntity assigns role roleID to entity entityID
func (permissions *Permissionist) AssignRoleToEntity(entityID string, appID string, roleID string) (string, error) {
	var id string
	err := permissions.DB.QueryRow(`
	INSERT INTO entity_roles (id, entity_id, app_id, role_id) VALUES ( 
		$1, $2, $3, $4 
	) RETURNING id;
	`, uuid.NewV4().String(), entityID, appID, roleID).Scan(&id)

	if err != nil {
		return "", errors.Wrap(err, "Could not assign role to entity")
	}

	return id, nil
}

// GrantPermissionToRole assigns permission of permissionID to role roleID
func (permissions *Permissionist) GrantPermissionToRole(roleID string, appID string, permissionID string) error {
	var id string
	err := permissions.DB.QueryRow(`
	INSERT INTO role_permissions (id, role_id, app_id, permission_id) VALUES (
		$1, $2, $3, $4
	) RETURNING id;
	`, uuid.NewV4().String(), roleID, appID, permissionID).Scan(&id)

	if err != nil {
		return errors.Wrap(err, "Could not assign permission to role")
	}

	return nil
}

// CreateApp creates a new app in the database
func (permissions *Permissionist) CreateApp(name string) (App, error) {
	var app App
	err := permissions.DB.QueryRow(`
	INSERT INTO apps (id, name) VALUES (
		$1, '$2'
	) RETURNING *;
	`, uuid.NewV4().String(), name).Scan(&app.ID, &app.Name)

	if err != nil {
		return app, errors.Wrap(err, "Could not create a new app")
	}

	return app, nil
}

// CreatePermission creates a new permission in the database
func (permissions *Permissionist) CreatePermission(permissionName string, appID string) (Permission, error) {
	var p Permission
	if len(permissionName) < 1 {
		return p, errors.New("Missing permission name")
	}
	if len(appID) < 1 {
		return p, errors.New("Missing app id")
	}
	err := permissions.DB.QueryRow(`
	INSERT INTO permissions (id, name, app_id) VALUES (
		$1, $2, $3
	) RETURNING id, name, app_id;
	`, uuid.NewV4().String(), permissionName, appID).Scan(&p.ID, &p.Name, &p.AppID)

	if err != nil {
		return p, errors.Wrap(err, "Could not create a new permission")
	}

	return p, nil
}

// CreatePermissions creates new permissions in the database
func (permissions *Permissionist) CreatePermissions(permissionNames []string, appID string) ([]Permission, error) {
	var newPermissions []Permission
	query := "INSERT INTO permissions (id, name, app_id) VALUES "
	for _, permissionName := range permissionNames {
		newPermission := Permission{
			ID:    uuid.NewV4().String(),
			Name:  permissionName,
			AppID: appID,
		}
		query += fmt.Sprintf(`('%s', '%s', '%s'),`, newPermission.ID, newPermission.Name, newPermission.AppID)
		newPermissions = append(newPermissions, newPermission)
	}
	query = strings.TrimSuffix(query, ",") + ";"
	_, err := permissions.DB.Exec(query)
	if err != nil {
		return nil, errors.Wrap(err, "Could not create new permissions")
	}

	return newPermissions, nil
}

// CreateRole creates a new role in the database
func (permissions *Permissionist) CreateRole(roleName string, appID string) (Role, error) {
	var role Role
	err := permissions.DB.QueryRow(`
	INSERT INTO roles (id, name, app_id) VALUES (
		$1, $2, $3
	) RETURNING *;
	`, uuid.NewV4().String(), roleName, appID).Scan(&role.ID, &role.Name, &role.AppID)

	if err != nil {
		return role, errors.Wrap(err, "Could not create a new role")
	}

	return role, nil
}

// CreateRoles creates a new role in the database
func (permissions *Permissionist) CreateRoles(roleNames []string, appID string) ([]Role, error) {
	var newRoles []Role
	query := "INSERT INTO roles (id, name, app_id) VALUES "
	for _, roleName := range roleNames {
		newRole := Role{
			ID:    uuid.NewV4().String(),
			Name:  roleName,
			AppID: appID,
		}
		query += fmt.Sprintf(`('%s', '%s', '%s'),`, newRole.ID, newRole.Name, newRole.AppID)
		newRoles = append(newRoles, newRole)
	}
	query = strings.TrimSuffix(query, ",") + " RETURNING name;"
	_, err := permissions.DB.Exec(query)
	if err != nil {
		return nil, errors.Wrap(err, "Could not create a new role")
	}

	return newRoles, nil
}
