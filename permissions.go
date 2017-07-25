package main

import (
	"github.com/jmoiron/sqlx"
	"fmt"
	"errors"
)

// RolePermissions role_permissions schema
type RolePermissions struct {
	ID           string `json:"id" db:"id"`
	RoledID      string `json:"role_id" db:"role_id"`
	PermissionID string `json:"permission_id" db:"permission_id"`
}

// Permissions stuff
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
		return false, fmt.Errorf("Could not check permission", err.Error())
	}

	if rolePermissions == nil || len(rolePermissions) <= 0 {
		return false, nil
	}

	return true, nil
}

// getPermissions returns a list of all permissions that belong to an entity
func (permissions *Permissionist) getPermissions(entityID string, appID string) ([]string, error) {
	var perms []string
	err := permissions.DB.Select(&perms, `
	select name from permissions as p 
	inner join role_permissions as rp 
		on p.id = rp.id and p.app_id = $2 and rp.app_id = $2
	inner join entity_roles as er 
		on er.entity_id = $1 and er.app_id = $2;
	`, entityID, appID)

	if err != nil {
		return nil, fmt.Errorf("Could not get permissions", err.Error())
	}

	return perms, nil
}

// AssignRoleToEntity assigns role roleName to entity entityID
func (permissions *Permissionist) AssignRoleToEntity(entityID string, appID string, roleName string) error {
	row := permissions.DB.QueryRow(`
	insert into entity_roles (entity_id, app_id, role_id) values (
		$1, $2, (select id from roles where name = $3)
	);
	`, entityID, appID, roleName)

	if row != nil {
		return errors.New("Could not assign role to entity")
	}

	return nil
}

// AssignPermissionToRole assigns permission of permissionName to role roleName
func (permissions *Permissionist) AssignPermissionToRole(roleName string, appID string, permissionName string) error {
	row := permissions.DB.QueryRow(`
	insert into role_permissions (role_id, app_id, permission_id) values (
		(select id from roles where name = $1), $2, (select name from permissions where name = $3)
	);
	`, roleName, appID, permissionName)

	if row != nil {
		return errors.New("Could not assign permission to role")
	}

	return nil
}

// CreatePermission creates a new permission in the database
func (permissions *Permissionist) CreatePermission(permissionName string, appID string) error {
	row := permissions.DB.QueryRow(`
	insert into permissions (name, app_id) values (
		$1, $2
	);
	`, permissionName, appID)

	if row != nil {
		return errors.New("Could not create a new permission")
	}

	return nil
}

// CreateRole creates a new role in the database
func (permissions *Permissionist) CreateRole(roleName string, appID string) error {
	row := permissions.DB.QueryRow(`
	insert into roles (name) values (
		$1, $2
	);
	`, roleName, appID)

	if row != nil {
		return errors.New("Could not create a new role")
	}

	return nil
}
