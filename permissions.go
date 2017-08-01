package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/satori/go.uuid"
)

// RolePermission role_permissions schema
type RolePermission struct {
	ID           string `json:"id" db:"id"`
	RoledID      string `json:"role_id" db:"role_id"`
	PermissionID string `json:"permission_id" db:"permission_id"`
}

// Permission permissions schema
type Permission struct {
	ID   string `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

// Role roles schema
type Role struct {
	ID   string `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

// Permissionist owns permissions crud
type Permissionist struct {
	DB *sqlx.DB
}

// Allowed checks if entity entityID has permission permissionID
func (permissions *Permissionist) Allowed(entityID string, appID string, permissionID string) (bool, error) {
	var rolePermissions []RolePermission
	err := permissions.DB.Select(&rolePermissions, `
		select * from role_permissions as rp 
			inner join permissions as p 
		on rp.permission_id = p.id and rp.app_id = $2 and p.id = $3
			inner join entity_roles as er 
		on er.app_id = $2 and er.entity_id = $1;
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
	select id from apps as a 
	inner join entity_roles as er 
		on a.id = er.app_id and er.id = $1;
	`, entityID)

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

// GetPermissionsByEntityID returns a list of all permissions that belong to an entity
func (permissions *Permissionist) GetPermissionsByEntityID(entityID string, appID string) ([]Permission, error) {
	var perms []Permission
	err := permissions.DB.Select(&perms, `
	select * from permissions as p 
	inner join role_permissions as rp 
		on p.id = rp.permission_id and p.app_id = $2 and rp.app_id = $2
	inner join entity_roles as er 
		on er.entity_id = $1 and er.app_id = $2;
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
	select * from permissions as p 
	inner join role_permissions as rp 
		on p.id = rp.permission_id and rp.role_id = $1 and p.app_id = $2 and rp.app_id = $2
	`, roleID, appID)

	if err != nil {
		return nil, fmt.Errorf("Could not get permissions: %s", err.Error())
	}

	return perms, nil
}

// GetRoles returns a list of all roles created for an app
func (permissions *Permissionist) GetRoles(appID string) ([]Role, error) {
	var roleIDs []Role
	err := permissions.DB.Select(&roleIDs, `
	select * from roles where app_id = $1;
	`, appID)

	if err != nil {
		return nil, fmt.Errorf("Could not get roles: %s", err.Error())
	}

	return roleIDs, nil
}

// GetRoleByID returns a role name
func (permissions *Permissionist) GetRoleByID(roleID string, appID string) (Role, error) {
	var role Role
	err := permissions.DB.Select(&role, `
	select * from roles where id = $1 and app_id = $2 limit 1;
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
	insert into entity_roles (id, entity_id, app_id, role_id) values (
		$1, $2, $3, (select id from roles where id = $4)
	) returning id;
	`, uuid.NewV4().String(), entityID, appID, roleID).Scan(&id)

	if err != nil {
		return "", fmt.Errorf("Could not assign role to entity: %s", err.Error())
	}

	return id, nil
}

// GrantPermissionToRole assigns permission of permissionID to role roleID
func (permissions *Permissionist) GrantPermissionToRole(roleID string, appID string, permissionID string) (string, error) {
	var id string
	err := permissions.DB.QueryRow(`
	insert into role_permissions (id, role_id, app_id, permission_id) values (
		$1, (select id from roles where id = $2), $3, (select id from permissions where id = $4)
	) returning id;
	`, uuid.NewV4().String(), roleID, appID, permissionID).Scan(&id)

	if err != nil {
		return "", fmt.Errorf("Could not assign permission to role: %s", err.Error())
	}

	return id, nil
}

// CreateApp creates a new app in the database
func (permissions *Permissionist) CreateApp(name string) (string, error) {
	var id string
	err := permissions.DB.QueryRow(`
	insert into apps (id, name) values (
		$1, $2
	) returning id;
	`, uuid.NewV4().String(), name).Scan(&id)

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
