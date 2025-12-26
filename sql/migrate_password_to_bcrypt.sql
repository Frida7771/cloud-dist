-- Migration script to update password column for bcrypt support
-- Bcrypt hashes are 60 characters long, so we need to increase the column size

ALTER TABLE `user_basic` 
MODIFY COLUMN `password` varchar(255) DEFAULT NULL COMMENT 'Bcrypt hashed password';

