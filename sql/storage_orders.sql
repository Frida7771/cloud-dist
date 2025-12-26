-- ----------------------------
-- Table structure for storage_orders
-- ----------------------------
DROP TABLE IF EXISTS `storage_orders`;
CREATE TABLE `storage_orders` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `identity` varchar(36) DEFAULT NULL COMMENT 'Order unique identifier',
  `user_identity` varchar(36) DEFAULT NULL COMMENT 'User who made the purchase',
  `stripe_session_id` varchar(255) DEFAULT NULL COMMENT 'Stripe Checkout Session ID',
  `stripe_payment_intent_id` varchar(255) DEFAULT NULL COMMENT 'Stripe Payment Intent ID',
  `storage_amount` bigint(20) DEFAULT NULL COMMENT 'Storage capacity purchased (bytes)',
  `price_amount` int(11) DEFAULT NULL COMMENT 'Price in cents (e.g., 999 = $9.99)',
  `currency` varchar(10) DEFAULT 'usd' COMMENT 'Currency code',
  `status` varchar(20) DEFAULT 'pending' COMMENT 'Order status: pending, paid, failed, refunded',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_identity` (`user_identity`),
  KEY `idx_stripe_session_id` (`stripe_session_id`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Storage purchase orders';

