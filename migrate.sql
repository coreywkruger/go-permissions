CREATE TABLE IF NOT EXISTS apps (
	id UUID PRIMARY KEY,
  name VARCHAR(60) NOT NULL
);

CREATE TABLE IF NOT EXISTS roles (
	id UUID PRIMARY KEY,
	app_id UUID NOT NULL REFERENCES apps,
	name VARCHAR(60) NOT NULL
);

CREATE TABLE IF NOT EXISTS permissions (
	id UUID PRIMARY KEY,
	app_id UUID NOT NULL REFERENCES apps,
	name VARCHAR(60) NOT NULL
);

CREATE TABLE IF NOT EXISTS role_permissions (
	id UUID PRIMARY KEY,
	permission_id UUID NOT NULL REFERENCES permissions,
	role_id UUID NOT NULL REFERENCES roles
);

CREATE TABLE IF NOT EXISTS entity_roles (
	id UUID PRIMARY KEY,
	role_id UUID NOT NULL REFERENCES roles,
	entity_id UUID NOT NULL
);