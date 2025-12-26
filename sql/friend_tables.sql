-- Friend system tables
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for friend
-- ----------------------------
DROP TABLE IF EXISTS `friend`;
CREATE TABLE `friend` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `identity` varchar(36) DEFAULT NULL,
  `user_identity` varchar(36) DEFAULT NULL COMMENT 'The user who has this friend',
  `friend_identity` varchar(36) DEFAULT NULL COMMENT 'The friend user identity',
  `status` varchar(20) DEFAULT 'active' COMMENT 'active, blocked',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_friend` (`user_identity`, `friend_identity`, `deleted_at`),
  KEY `idx_user_identity` (`user_identity`),
  KEY `idx_friend_identity` (`friend_identity`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for friend_request
-- ----------------------------
DROP TABLE IF EXISTS `friend_request`;
CREATE TABLE `friend_request` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `identity` varchar(36) DEFAULT NULL,
  `from_user_identity` varchar(36) DEFAULT NULL COMMENT 'User who sent the request',
  `to_user_identity` varchar(36) DEFAULT NULL COMMENT 'User who received the request',
  `status` varchar(20) DEFAULT 'pending' COMMENT 'pending, accepted, rejected',
  `message` varchar(500) DEFAULT NULL COMMENT 'Optional message',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_from_user` (`from_user_identity`),
  KEY `idx_to_user` (`to_user_identity`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for friend_share
-- ----------------------------
DROP TABLE IF EXISTS `friend_share`;
CREATE TABLE `friend_share` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `identity` varchar(36) DEFAULT NULL,
  `from_user_identity` varchar(36) DEFAULT NULL COMMENT 'User who shared the file',
  `to_user_identity` varchar(36) DEFAULT NULL COMMENT 'Friend who received the share',
  `repository_identity` varchar(36) DEFAULT NULL COMMENT 'The shared file',
  `user_repository_identity` varchar(36) DEFAULT NULL COMMENT 'User file reference',
  `message` varchar(500) DEFAULT NULL COMMENT 'Optional message',
  `is_read` tinyint(1) DEFAULT 0 COMMENT 'Whether the friend has read it',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_from_user` (`from_user_identity`),
  KEY `idx_to_user` (`to_user_identity`),
  KEY `idx_repository` (`repository_identity`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

SET FOREIGN_KEY_CHECKS = 1;

