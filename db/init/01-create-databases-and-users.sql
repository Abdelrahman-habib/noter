-- Create roles and databases (development defaults).
-- Statements are guarded with \gexec so psql re-runs are safe.

-- Create roles
SELECT format('CREATE ROLE %I LOGIN PASSWORD %L', 'noter_admin', 'admin')
WHERE NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'noter_admin') \gexec
SELECT format('CREATE ROLE %I LOGIN PASSWORD %L', 'noter_web', 'pass')
WHERE NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'noter_web') \gexec
SELECT format('CREATE ROLE %I LOGIN PASSWORD %L', 'noter_test_web', 'test_pass')
WHERE NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'noter_test_web') \gexec

-- Create databases
SELECT format('CREATE DATABASE %I OWNER %I', 'noter', 'noter_admin')
WHERE NOT EXISTS (SELECT 1 FROM pg_database WHERE datname = 'noter') \gexec
SELECT format('CREATE DATABASE %I OWNER %I', 'noter_test', 'noter_test_web')
WHERE NOT EXISTS (SELECT 1 FROM pg_database WHERE datname = 'noter_test') \gexec

-- Grant noter_web access to tables and sequences created by noter_admin
\connect noter
GRANT USAGE ON SCHEMA public TO noter_web;
ALTER DEFAULT PRIVILEGES FOR ROLE noter_admin IN SCHEMA public GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO noter_web;
ALTER DEFAULT PRIVILEGES FOR ROLE noter_admin IN SCHEMA public GRANT USAGE, SELECT ON SEQUENCES TO noter_web;
