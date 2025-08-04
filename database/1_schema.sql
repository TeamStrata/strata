-- DESCRIPTION
-- This script will setup the schema for the strata database.
-- 
-- EXAMPLE USAGE
-- psql.exe -U <postgres username> -d strata -f .\database\1_schema.sql

DO
$do$
BEGIN
	-- Users table
	IF EXISTS (
		SELECT 1
		FROM information_schema.tables
		WHERE table_schema = 'public'
	) THEN
		DROP SCHEMA public CASCADE;
		CREATE SCHEMA public;
	END IF;

	DO $$ BEGIN
		-- Enums
		CREATE TYPE CHART_TYPE AS ENUM (
			'line',
			'area',
			'column',
			'bar',
			'treemap',
			'heatmap',
			'pie',
			'radar',
			'polar',
			'scatter'
		);
	EXCEPTION
		WHEN duplicate_object THEN null;
	END $$;

	-- Table for storing users
	CREATE TABLE users (
		user_id SERIAL PRIMARY KEY,
		user_name TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL
	);

	-- Table for storing Roles
	CREATE TABLE IF NOT EXISTS roles (
		role_id SERIAL PRIMARY KEY,
		role_name TEXT UNIQUE NOT NULL,
		role_color TEXT
	);

	-- Scope for permissions
	CREATE TYPE scope AS ENUM ('global', 'dashboard'); 

	-- Table for storing Permissions
	CREATE TABLE IF NOT EXISTS permissions (
		permission_id SERIAL PRIMARY KEY,
		permission_name TEXT UNIQUE NOT NULL,
		permission_scope scope NOT NULL
	);

	-- Table for storing custom, or any other queries. 
	CREATE TABLE IF NOT EXISTS queries (
		query_id SERIAL PRIMARY KEY,
		query_name TEXT UNIQUE NOT NULL,
		query_string TEXT NOT NULL
	);

	-- Table for storing charts
	CREATE TABLE IF NOT EXISTS chart (
		chart_id SERIAL PRIMARY KEY,
		chart_title TEXT NOT NULL,
		chart_type CHART_TYPE NOT NULL
	);

	-- Chart series
	CREATE TABLE IF NOT EXISTS chartSeries (
		series_id SERIAL PRIMARY KEY,
		chart_id INTEGER NOT NULL references chart(chart_id) ON DELETE CASCADE,
		query_id INTEGER NOT NULL references queries(query_id) ON DELETE CASCADE,
		x_column TEXT NOT NULL,
		y_column TEXT NOT NULL
	);

	-- Dashboard 
	CREATE TABLE IF NOT EXISTS dashboards (
		dashboard_id SERIAL PRIMARY KEY,
		dashboard_title TEXT,
		dashboard_content TEXT
	);

	-- Dashboard Graphs
	CREATE TABLE IF NOT EXISTS dashboardGraphs (
		dashboard_id INTEGER NOT NULL references dashboards(dashboard_id),
		chart_id INTEGER NOT NULL references chart(chart_id),
		size_x INTEGER NOT NULL,
		size_y INTEGER NOT NULL,
		chart_order INTEGER NOT NULL
	);

	---=== RELATIONS ===---
	-- Table to store relations between roles and permissions
	CREATE TABLE IF NOT EXISTS rolePermissions (
		role_id INTEGER NOT NULL references roles(role_id) ON DELETE CASCADE,
		permission_id INTEGER NOT NULL references permissions(permission_id) ON DELETE CASCADE,
		PRIMARY KEY (role_id, permission_id)
	);

	CREATE TABLE IF NOT EXISTS dashboardRolePermissions (
		dash_id INTEGER NOT NULL references dashboards(dashboard_id) ON DELETE CASCADE,
		role_id INTEGER NOT NULL references roles(role_id) ON DELETE CASCADE,
		permission_id INTEGER NOT NULL references permissions(permission_id) ON DELETE CASCADE,
		PRIMARY KEY (role_id, dash_id, permission_id),
		FOREIGN KEY (role_id, permission_id) REFERENCES rolePermissions(role_id, permission_id)
	);

	-- Table to store relations between users and roles
	CREATE TABLE IF NOT EXISTS userRoles (
		user_id INTEGER NOT NULL references users(user_id) ON DELETE CASCADE,
		role_id INTEGER NOT NULL references roles(role_id) ON DELETE CASCADE
	);

	-- Table to store relations between permissions and queries
	CREATE TABLE IF NOT EXISTS queryPermissions (
		query_id INTEGER NOT NULL references queries(query_id) ON DELETE CASCADE,
		role_id INTEGER NOT NULL references roles(role_id) ON DELETE CASCADE
	);

	-- Table to store settings such as the client database connection string
	CREATE TABLE IF NOT EXISTS settings (
		skey TEXT UNIQUE NOT NULL,
		svalue TEXT NOT NULL
	);
END
$do$
