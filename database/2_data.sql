-- DESCRIPTION
-- This script will populate the strata database with test data.
-- 
-- EXAMPLE USAGE
-- psql.exe -U <postgres username> -d strata -f .\database\data.sql

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
		INSERT INTO users (user_name, password_hash)
		VALUES ('gopher', '123');
	END IF;
END
$do$