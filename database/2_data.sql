-- DESCRIPTION
-- This script will populate the strata database with test data.
-- 
-- EXAMPLE USAGE
-- psql.exe -U <postgres username> -d strata -f .\database\2_data.sql

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
		VALUES ('admin', '$2a$10$GEqRMhGYBay/4uXY50eyP.heui16Vs9WC//cwxt9mHijfJ.4xvi9.');
	END IF;
END
$do$
