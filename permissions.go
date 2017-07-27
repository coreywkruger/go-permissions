package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/satori/go.uuid"
)

// RolePermissions role_permissions schema
type RolePermissions struct {
	ID           string `json:"id" db:"id"`
	RoledID      string `json:"role_id" db:"role_id"`
	PermissionID string `json:"permission_id" db:"permission_id"`
}

// Permissionist owns permissions crud
type Permissionist struct {
	DB *sqlx.DB
}

// Allowed checks if entity entityID has permission permissionName
func (permissions *Permissionist) Allowed(entityID string, appID string, permissionName string) (bool, error) {
	var rolePermissions []RolePermissions
	err := permissions.DB.Select(&rolePermissions,
		`select * from role_permissions as rp 
		inner join permissions as p 
	on rp.permission_id = p.id and rp.app_id = $2 and p.name = $3
		inner join entity_roles as er 
	on er.app_id = $2 and er.entity_id = $1;
	`, entityID, appID, permissionName)

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

// GetApp returns an app by id
func (permissions *Permissionist) GetApp(appID string) (string, error) {
	var id string
	err := permissions.DB.Select(&id, `select id from apps where id = $1;`, appID)
	if err != nil {
		return "", fmt.Errorf("Could not get app: %s", err.Error())
	}

	return id, nil
}

// GetPermissions returns a list of all permissions that belong to an entity
func (permissions *Permissionist) GetPermissions(entityID string, appID string) ([]string, error) {
	var perms []string
	err := permissions.DB.Select(&perms, `
	select name from permissions as p 
	inner join role_permissions as rp 
		on p.id = rp.id and p.app_id = $2 and rp.app_id = $2
	inner join entity_roles as er 
		on er.entity_id = $1 and er.app_id = $2;
	`, entityID, appID)

	if err != nil {
		return nil, fmt.Errorf("Could not get permissions: %s", err.Error())
	}

	return perms, nil
}

// GetRoles returns a list of all roles created for an app
func (permissions *Permissionist) GetRoles(appID string) ([]string, error) {
	var roleIDs []string
	err := permissions.DB.Select(&roleIDs, `
	select id from roles where app_id = $1;
	`, appID)

	if err != nil {
		return nil, fmt.Errorf("Could not get roles: %s", err.Error())
	}

	return roleIDs, nil
}

// AssignRoleToEntity assigns role roleName to entity entityID
func (permissions *Permissionist) AssignRoleToEntity(entityID string, appID string, roleName string) (string, error) {
	var id string
	err := permissions.DB.QueryRow(`
	insert into entity_roles (id, entity_id, app_id, role_id) values (
		$1, $2, $3, (select id from roles where name = $3)
	) returning id;
	`, uuid.NewV4().String(), entityID, appID, roleName).Scan(&id)

	if err != nil {
		return "", fmt.Errorf("Could not assign role to entity: %s", err.Error())
	}

	return id, nil
}

// AssignPermissionToRole assigns permission of permissionName to role roleName
func (permissions *Permissionist) AssignPermissionToRole(roleName string, appID string, permissionName string) (string, error) {
	var id string
	err := permissions.DB.QueryRow(`
	insert into role_permissions (id, role_id, app_id, permission_id) values (
		$1, (select id from roles where name = $2), $3, (select name from permissions where name = $4)
	) returning id;
	`, uuid.NewV4().String(), roleName, appID, permissionName).Scan(&id)

	if err != nil {
		return "", fmt.Errorf("Could not assign permission to role: %s", err.Error())
	}

	return id, nil
}

// CreateApp creates a new app in the database
func (permissions *Permissionist) CreateApp() (string, error) {
	var id string
	err := permissions.DB.QueryRow(`
	insert into apps (id) values (
		$1
	) returning id;
	`, uuid.NewV4().String()).Scan(&id)

	if err != nil {
		return "", fmt.Errorf("Could not create a new app: %s", err.Error())
	}

	return id, nil
}

// CreatePermission creates a new permission in the database
func (permissions *Permissionist) CreatePermission(permissionName string, appID string) (string, error) {
	var id string
	err := permissions.DB.QueryRow(`
	insert into permissions (id, name, app_id) values (
		$1, $2, $3
	) returning id;
	`, uuid.NewV4().String(), permissionName, appID).Scan(&id)

	if err != nil {
		return "", fmt.Errorf("Could not create a new permission: %s", err.Error())
	}

	return id, nil
}

// CreateRole creates a new role in the database
func (permissions *Permissionist) CreateRole(roleName string, appID string) (string, error) {
	var id string
	err := permissions.DB.QueryRow(`
	insert into roles (id, name, app_id) values (
		$1, $2, $3
	) returning id;
	`, uuid.NewV4().String(), roleName, appID).Scan(&id)

	if err != nil {
		return "", fmt.Errorf("Could not create a new role: %s", err.Error())
	}

	return id, nil
}
