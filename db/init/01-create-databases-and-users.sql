-- Create roles (development passwords)
CREATE ROLE noter_admin LOGIN PASSWORD 'admin';
CREATE ROLE noter_web   LOGIN PASSWORD 'pass';
CREATE ROLE noter_test_web LOGIN PASSWORD 'test_pass';

-- Create databases
CREATE DATABASE noter OWNER noter_admin;
CREATE DATABASE noter_test OWNER noter_test_web;

-- Grant noter_web access to tables and sequences created by noter_admin
\connect noter
GRANT USAGE ON SCHEMA public TO noter_web;
ALTER DEFAULT PRIVILEGES FOR ROLE noter_admin IN SCHEMA public GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO noter_web;
ALTER DEFAULT PRIVILEGES FOR ROLE noter_admin IN SCHEMA public GRANT USAGE, SELECT ON SEQUENCES TO noter_web;
