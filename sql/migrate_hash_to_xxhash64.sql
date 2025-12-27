-- Migrate hash field from varchar(64) to varchar(16) to support xxHash64
-- xxHash64 produces 64-bit hash, which is 16 hexadecimal characters

ALTER TABLE `repository_pool` 
MODIFY COLUMN `hash` varchar(16) DEFAULT NULL COMMENT 'Unique identifier of the file (xxHash64 hash)';


