CREATE DATABASE IF NOT EXISTS `cloud_dist` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `cloud_dist`;

-- Import main system tables
SOURCE sql/cloud-dist.sql;

-- Import friend system tables
SOURCE sql/friend_tables.sql;

-- Import storage orders table
SOURCE sql/storage_orders.sql;
