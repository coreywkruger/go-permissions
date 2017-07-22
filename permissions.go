package permissionist

import (
	"github.com/jmoiron/sqlx"
	"errors"
)

// RolePermissions role_permissions schema
type RolePermissions struct {
	ID           string `json:"id" db:"id"`
	RoledID      string `json:"role_id" db:"role_id"`
	PermissionID string `json:"permission_id" db:"permission_id"`
}

// InitDB initializes the database
func InitDB(dbURI string) (*sqlx.DB, error) {
	return nil, nil
	db, err := sqlx.Connect("postgres", dbURI)

	if err != nil {
    	return nil, err
	}

	return db.Unsafe(), nil
}

// Permissions stuff
type Permissionist struct {
	DB *sqlx.DB
}

// Allowed checks if entity entityID has permission permissionName
func (permissions *Permissionist) Allowed(entityID string, permissionName string) (bool, error) {
	var rolePermissions []RolePermissions
	err := permissions.DB.Select(&rolePermissions, 
	`select * from role_permissions as rp 
		inner join permissions as p 
	on rp.permission_id = p.id and p.name = $2
		inner join entity_roles as er 
	on er.entity_id = $1;`, entityID, permissionName)

	if err != nil {
		return false, err
	}

	if rolePermissions == nil || len(rolePermissions) <= 0 {
		return false, nil
	}

	return true, nil
}

// AssignRoleToEntity assigns role roleName to entity entityID
func (permissions *Permissionist) AssignRoleToEntity(entityID string, roleName string) error {
	row := permissions.DB.QueryRow(`
	insert into entity_roles (entity_id, role_id) values (
		$1, (select id from roles where name = $2)
	);`, entityID, roleName)

	if row != nil {
		return errors.New("error")
	}

	return nil
}

// AssignPermissionToRole assigns permission of permissionName to role roleName
func (permissions *Permissionist) AssignPermissionToRole(roleName string, permissionName string) error {
	row := permissions.DB.QueryRow(`
	insert into role_permissions (role_id, permission_id) values (
		(select id from roles where name = $1), 
		(select name from permissions where name = $2)
	);`, roleName, permissionName)

	if row != nil {
		return errors.New("error")
	}

	return nil
}

// CreatePermission creates a new permission in the database
func (permissions *Permissionist) CreatePermission(permissionName string) error {
	row := permissions.DB.QueryRow(`
	insert into permissions (name) values (
		$1
	);`, permissionName)

	if row != nil {
		return errors.New("error")
	}

	return nil
}

// CreateRole creates a new role in the database
func (permissions *Permissionist) CreateRole(roleName string) error {
	row := permissions.DB.QueryRow(`
	insert into roles (name) values (
		$1
	);`, roleName)

	if row != nil {
		return errors.New("error")
	}

	return nil
}
