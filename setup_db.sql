CREATE DATABASE IF NOT EXISTS `cloud-dist` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `cloud-dist`;

-- Import main system tables
SOURCE sql/cloud-dist.sql;

-- Import friend system tables
SOURCE sql/friend_tables.sql;
