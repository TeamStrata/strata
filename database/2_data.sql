-- DESCRIPTION
-- This script will populate the strata database with test data.
-- 
-- EXAMPLE USAGE
-- psql.exe -U <postgres username> -d strata -f .\database\2_data.sql

DO
$do$
DECLARE
	admin_user_id INTEGER;
	admin_role_id INTEGER;
BEGIN
	-- Users table
	IF EXISTS (
		SELECT 1
		FROM information_schema.tables
		WHERE table_schema = 'public'
		AND table_name = 'users'
	) THEN
		-- Create 'admin' user
		INSERT INTO users (user_name, password_hash)
		VALUES ('admin', '$2a$10$GEqRMhGYBay/4uXY50eyP.heui16Vs9WC//cwxt9mHijfJ.4xvi9.')
		ON CONFLICT (user_name) DO NOTHING;

		-- Get user id
		SELECT user_id INTO admin_user_id FROM users WHERE user_name = 'admin';

		-- Create 'default' role
		INSERT INTO roles (role_name)
		VALUES ('default')
		ON CONFLICT (role_name) DO NOTHING;

		-- Create 'admin' role
		INSERT INTO roles (role_name, role_color)
		VALUES ('admin', '00ADD8')
		ON CONFLICT (role_name) DO NOTHING;

		-- Get the 'admin' role id
		SELECT role_id INTO admin_role_id FROM roles WHERE role_name = 'admin';

		-- Associate role to user via userroles
		INSERT INTO userroles (user_id, role_id)
		VALUES (admin_user_id, admin_role_id)
		ON CONFLICT DO NOTHING;


		INSERT INTO settings (skey, svalue)
		VALUES ('cdb', 'postgresql://strata:atarts@strata_psql:5432/strata'),
		('ollama_host', 'strata_ollama:11434'),
		('ollama_model', 'stratabot')
		ON CONFLICT DO NOTHING;

		-- Dashboards
		INSERT INTO dashboards (dashboard_title, dashboard_content)
		VALUES
			('Main Dashboard', 'This is the main dashboard.');

		-- Insert permissions
		INSERT INTO permissions (permission_name, permission_scope)
		VALUES
			('admin_privileges', 'global'),
			('create_dashboard', 'global'),
			('delete_dashboard', 'global'),
			('edit_dashboard', 'dashboard'),
			('view_dashboard', 'dashboard')
		ON CONFLICT DO NOTHING;

		-- Assign the 'admin' role the 'admin_privileges` permission.
		INSERT INTO rolePermissions (role_id, permission_id)
		VALUES
			(2, 1),
			(2, 2),
			(2, 3),
			(2, 4),
			(2, 5);

		INSERT INTO dashboardRolePermissions (dash_id, role_id, permission_id)
		VALUES
			(1, 2, 3),
			(1, 2, 4),
			(1, 2, 5);
	END IF;
END
$do$
