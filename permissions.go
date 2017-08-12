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

// EntityIsAllowed checks if entity entityID has permission permissionID
func (permissions *Permissionist) EntityIsAllowed(entityID string, permissionID string) (bool, error) {
	var entityRoleIDs []string
	err := permissions.DB.Select(&entityRoleIDs, `
	SELECT p.id
	FROM permissions AS p
	INNER JOIN entity_roles AS er
		ON er.entity_id = $1
	INNER JOIN role_permissions AS rp
		ON rp.permission_id = $2
			AND rp.permission_id = p.id
			AND rp.role_id = er.role_id;
	`, entityID, permissionID)

	if err != nil {
		return false, errors.Wrap(err, "Could not check permission")
	}

	if len(entityRoleIDs) == 0 {
		return false, nil
	}

	return true, nil
}

// RoleIsAllowed checks if entity roleID has permission permissionID
func (permissions *Permissionist) RoleIsAllowed(roleID string, permissionID string) (bool, error) {
	var rolePermissionIDs []string
	err := permissions.DB.Select(&rolePermissionIDs, `
	SELECT rp.id
	FROM permissions AS p
	INNER JOIN role_permissions AS rp
		ON p.id = $2
			AND rp.permission_id = p.id
			AND rp.role_id = $1;
	`, roleID, permissionID)

	if err != nil {
		return false, errors.Wrap(err, "Could not check permission")
	}

	if len(rolePermissionIDs) <= 0 {
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
func (permissions *Permissionist) GetPermissionsByRoleID(roleID string) ([]Permission, error) {
	var perms []Permission
	err := permissions.DB.Select(&perms, `
	SELECT p.id, p.name, p.app_id
	FROM permissions AS p
	INNER JOIN role_permissions AS rp
		ON p.id = rp.permission_id
			AND rp.role_id = $1;
	`, roleID)

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
	WHERE app_id = $1;
	`, appID)

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
func (permissions *Permissionist) AssignRoleToEntity(entityID string, roleID string) error {
	_, err := permissions.DB.Exec(`
	INSERT INTO entity_roles AS er (id, entity_id, role_id) VALUES (
		$1, $2, $3 
	);
	`, uuid.NewV4().String(), entityID, roleID)

	if err != nil {
		return errors.Wrap(err, "Could not assign role to entity")
	}

	return nil
}

// GrantPermissionToRole assigns permission of permissionID to role roleID
func (permissions *Permissionist) GrantPermissionToRole(roleID string, permissionID string) error {
	_, err := permissions.DB.Exec(`
	INSERT INTO role_permissions (id, role_id, permission_id) VALUES (
		$1, $2, $3
	);
	`, uuid.NewV4().String(), roleID, permissionID)

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
