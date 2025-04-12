-- DESCRIPTION
-- This script will setup the schema for the strata database.
-- 
-- EXAMPLE USAGE
-- psql.exe -U <postgres username> -d strata -f .\database\schema.sql

DO
$do$
BEGIN
	-- Users table
	IF EXISTS (
		SELECT 1
		FROM information_schema.tables
		WHERE table_schema = 'public'
		AND table_name = 'users'
	) THEN
		DROP TABLE public.users;
	END IF;

	-- Table for storing users
	CREATE TABLE users (
		user_id SERIAL PRIMARY KEY,
		user_name TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL
	);

	-- Table for storing Roles
	CREATE TABLE IF NOT EXISTS roles (
		role_id SERIAL PRIMARY KEY,
		role_name TEXT UNIQUE NOT NULL
	);

	-- Table for storing Permissions
	CREATE TABLE IF NOT EXISTS permissions (
		permission_id SERIAL PRIMARY KEY,
		permission_name TEXT UNIQUE NOT NULL
	);

	-- Table for storing custom, or any other queries. 
	CREATE TABLE IF NOT EXISTS queries (
		query_id SERIAL PRIMARY KEY,
		query_string TEXT UNIQUE NOT NULL
	);


	---=== RELATIONS ===---
	-- Table to store relations between roles and permissions
	CREATE TABLE IF NOT EXISTS rolePermissions (
		role_id INTEGER NOT NULL references roles(role_id),
		permission_id INTEGER NOT NULL references permissions(permission_id)
	);

	-- Table to store relations between users and roles
	CREATE TABLE IF NOT EXISTS userRoles (
		user_id INTEGER NOT NULL references users(user_id),
		role_id INTEGER NOT NULL references roles(role_id)
	);

	-- Table to store relations between permissions and queries
	CREATE TABLE IF NOT EXISTS queryPermissions (
		query_id INTEGER NOT NULL references queries(query_id),
		role_id INTEGER NOT NULL references roles(role_id)
	);
END
$do$