/*
 Navicat Premium Data Transfer

 Source Server         : 192.168.1.8
 Source Server Type    : MySQL
 Source Server Version : 50736
 Source Host           : 192.168.1.8:3306
 Source Schema         : cloud-disk

 Target Server Type    : MySQL
 Target Server Version : 50736
 File Encoding         : 65001

 Date: 05/05/2022 21:05:15
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for repository_pool
-- ----------------------------
DROP TABLE IF EXISTS `repository_pool`;
CREATE TABLE `repository_pool` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `identity` varchar(36) DEFAULT NULL,
  `hash` varchar(32) DEFAULT NULL COMMENT 'Unique identifier of the file',
  `name` varchar(255) DEFAULT NULL,
  `ext` varchar(30) DEFAULT NULL COMMENT 'File extension',
  `size` int(11) DEFAULT NULL COMMENT 'File size',
  `path` text DEFAULT NULL COMMENT 'File path (S3 URL or presigned URL)',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for share_basic
-- ----------------------------
DROP TABLE IF EXISTS `share_basic`;
CREATE TABLE `share_basic` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `identity` varchar(36) DEFAULT NULL,
  `user_identity` varchar(36) DEFAULT NULL,
  `repository_identity` varchar(36) DEFAULT NULL COMMENT 'Unique identifier in the public pool',
  `user_repository_identity` varchar(36) DEFAULT NULL COMMENT 'Unique identifier in the user pool',
  `expired_time` int(11) DEFAULT NULL COMMENT 'Expiration time in seconds, [0-never expires]',
  `click_num` int(11) DEFAULT '0' COMMENT 'Click count',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for user_basic
-- ----------------------------
DROP TABLE IF EXISTS `user_basic`;
CREATE TABLE `user_basic` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `identity` varchar(36) DEFAULT NULL,
  `name` varchar(60) DEFAULT NULL,
  `password` varchar(255) DEFAULT NULL COMMENT 'Bcrypt hashed password',
  `email` varchar(100) DEFAULT NULL,
  `now_volume` bigint(20) DEFAULT '0' COMMENT 'Currently used capacity (bytes)',
  `total_volume` bigint(20) DEFAULT '16106127360' COMMENT 'Total capacity (bytes), default 15GB',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for user_repository
-- ----------------------------
DROP TABLE IF EXISTS `user_repository`;
CREATE TABLE `user_repository` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `identity` varchar(36) DEFAULT NULL,
  `user_identity` varchar(36) DEFAULT NULL,
  `parent_id` int(11) DEFAULT NULL,
  `repository_identity` varchar(36) DEFAULT NULL,
  `ext` varchar(255) DEFAULT NULL COMMENT 'File or folder type',
  `name` varchar(255) DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8;

SET FOREIGN_KEY_CHECKS = 1;
