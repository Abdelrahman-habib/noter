-- Create databases
CREATE DATABASE IF NOT EXISTS `noter`;
CREATE DATABASE IF NOT EXISTS `noter_test`;

-- Create users with proper passwords
CREATE USER IF NOT EXISTS 'noter_admin'@'%' IDENTIFIED BY 'admin';
CREATE USER IF NOT EXISTS 'noter_web'@'%' IDENTIFIED BY 'pass';
CREATE USER IF NOT EXISTS 'noter_test_web'@'%' IDENTIFIED BY 'test_pass';

-- Grant permissions to noter_admin (full privileges on noter database)
GRANT ALL PRIVILEGES ON `noter`.* TO `noter_admin`@`%`;

-- Grant permissions to noter_web (limited privileges on noter database)
GRANT SELECT, INSERT, UPDATE, DELETE, REFERENCES ON `noter`.* TO `noter_web`@`%`;

-- Grant permissions to noter_test_web (full privileges on noter_test database)
GRANT ALL PRIVILEGES ON `noter_test`.* TO `noter_test_web`@`%`;

-- Flush privileges to ensure they take effect
FLUSH PRIVILEGES;
