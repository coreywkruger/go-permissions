package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
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
func (permissions *Permissionist) Allowed(entityID string, appID string, permissionID string) (bool, error) {
	var rolePermissions []RolePermission
	err := permissions.DB.Select(&rolePermissions, `
	SELECT *
	FROM role_permissions AS rp
	INNER JOIN permissions AS p
		ON rp.permission_id = p.id
			AND rp.app_id = $2
			AND p.id = $3
	INNER JOIN entity_roles AS er
		ON er.app_id = $2
			AND er.entity_id = $1;
	`, entityID, appID, permissionID)

	if err != nil {
		return false, fmt.Errorf("Could not check permission: %s", err.Error())
	}

	if rolePermissions == nil || len(rolePermissions) <= 0 {
		return false, nil
	}

	return true, nil
}

// GetApps returns a list of all apps
func (permissions *Permissionist) GetApps() ([]string, error) {
	var apps []string
	err := permissions.DB.Select(&apps, `select id from apps;`)
	if err != nil {
		return nil, fmt.Errorf("Could not get apps: %s", err.Error())
	}

	return apps, nil
}

// GetAppsByEntityID returns a list of all apps
func (permissions *Permissionist) GetAppsByEntityID(entityID string) ([]string, error) {
	var apps []string
	err := permissions.DB.Select(&apps, `
	SELECT id
	FROM apps AS a
	INNER JOIN entity_roles AS er
		ON a.id = er.app_id
			AND er.id = $1;
	`, entityID)

	if err != nil {
		return nil, fmt.Errorf("Could not get apps: %s", err.Error())
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
		return "", fmt.Errorf("Could not get app: %s", err.Error())
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
		return nil, fmt.Errorf("Could not get permissions: %s", err.Error())
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
		return nil, fmt.Errorf("Could not get permissions: %s", err.Error())
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
		return nil, fmt.Errorf("Could not get roles: %s", err.Error())
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
		return role, fmt.Errorf("Could not get role: %s", err.Error())
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
		return "", fmt.Errorf("Could not assign role to entity: %s", err.Error())
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
		return fmt.Errorf("Could not assign permission to role: %s", err.Error())
	}

	return nil
}

// CreateApp creates a new app in the database
func (permissions *Permissionist) CreateApp(name string) (*App, error) {
	var app App
	err := permissions.DB.QueryRow(`
	INSERT INTO apps (id, name) VALUES (
		$1, $2
	) RETURNING *;
	`, uuid.NewV4().String(), name).Scan(&app.ID, &app.Name)

	if err != nil {
		return nil, fmt.Errorf("Could not create a new app: %s", err.Error())
	}

	return &app, nil
}

// CreatePermission creates a new permission in the database
func (permissions *Permissionist) CreatePermission(permissionName string, appID string) (*Permission, error) {
	var p Permission
	err := permissions.DB.QueryRow(`
	INSERT INTO permissions (id, name, app_id) VALUES (
		$1, $2, $3
	) RETURNING *;
	`, uuid.NewV4().String(), permissionName, appID).Scan(&p.ID, &p.Name, &p.AppID)

	if err != nil {
		return nil, fmt.Errorf("Could not create a new permission: %s", err.Error())
	}

	return &p, nil
}

// CreateRole creates a new role in the database
func (permissions *Permissionist) CreateRole(roleName string, appID string) (*Role, error) {
	var role Role
	err := permissions.DB.QueryRow(`
	INSERT INTO roles (id, name, app_id) VALUES (
		$1, $2, $3
	) RETURNING *;
	`, uuid.NewV4().String(), roleName, appID).Scan(&role.ID, &role.Name, &role.AppID)

	if err != nil {
		return nil, fmt.Errorf("Could not create a new role: %s", err.Error())
	}

	return &role, nil
}

// CreateRoles creates a new role in the database
func (permissions *Permissionist) CreateRoles(roleNames []string, appID string) ([]Role, error) {
	var newRoles []Role
	query := "INSERT INTO roles (id, name, app_id) VALUES "
	for _, roleName := range roleNames {
		newRole := Role{
			ID: uuid.NewV4().String(), 
			Name: roleName, 
			AppID: appID,
		}
		query += `('` + newRole.ID + `', '` + newRole.Name + `', '` + newRole.AppID + `'), `
		newRoles = append(newRoles, newRole)
	}
	query = strings.TrimSuffix(query, ", ")
	_, err := permissions.DB.Exec(query + ";")
	if err != nil {
		return nil, fmt.Errorf("Could not create a new role: %s", err.Error())
	}

	return newRoles, nil
}
