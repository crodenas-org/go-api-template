-- One-time database and user setup.
-- Run as the postgres superuser:
--   psql -h localhost -U postgres -f db/setup.sql
--
-- Passwords should be sourced from your secrets manager (Ansible Vault).
-- Replace :hw_admin_password and :hw_app_password with actual values,
-- or use \set in psql:
--   psql -h localhost -U postgres \
--     -v hw_admin_password='...' \
--     -v hw_app_password='...' \
--     -f db/setup.sql

-- Database
CREATE DATABASE hello_world_go;

-- Migration user (DDL privileges)
CREATE USER hw_admin WITH PASSWORD :'hw_admin_password';
GRANT ALL PRIVILEGES ON DATABASE hello_world_go TO hw_admin;
ALTER DATABASE hello_world_go OWNER TO hw_admin;

-- Application user (DML only)
CREATE USER hw_app WITH PASSWORD :'hw_app_password';

-- Connect to the app DB to grant schema-level privileges
\c hello_world_go

GRANT USAGE ON SCHEMA public TO hw_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO hw_app;
ALTER DEFAULT PRIVILEGES IN SCHEMA public
    GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO hw_app;
