BEGIN;

-- Table: fivenet_oauth2_accounts
CREATE TABLE IF NOT EXISTS `fivenet_oauth2_accounts` (
  `account_id` bigint(20) unsigned NOT NULL,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `provider` varchar(255) NOT NULL,
  `external_id` bigint(20) unsigned NOT NULL,
  `username` varchar(255) NOT NULL,
  `avatar` varchar(255) NOT NULL,
  UNIQUE KEY `idx_fivenet_oauth2_accounts_unique` (`account_id`,`provider`),
  UNIQUE KEY `idx_fivenet_oauth2_accounts_provider_external_id` (`provider`,`external_id`),
  CONSTRAINT `fk_fivenet_oauth2_accounts_account_id` FOREIGN KEY (`account_id`) REFERENCES `fivenet_accounts` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

COMMIT;
