BEGIN;

CREATE TABLE IF NOT EXISTS `fivenet_user_accounts` (
  `user_id` int(11) NOT NULL,
  `account_id` bigint(20) unsigned NOT NULL,
  PRIMARY KEY (`user_id`),
  KEY `idx_fivenet_user_accounts_account_id` (`account_id`),
  CONSTRAINT `fk_fivenet_user_accounts_user_id` FOREIGN KEY (`user_id`) REFERENCES `fivenet_user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_user_accounts_account_id` FOREIGN KEY (`account_id`) REFERENCES `fivenet_accounts` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

INSERT INTO `fivenet_user_accounts` (`user_id`, `account_id`)
SELECT u.`id`, a.`id`
FROM `fivenet_user` u
JOIN `fivenet_accounts` a
  ON a.`license` = COALESCE(NULLIF(u.`license`, ''), SUBSTRING_INDEX(u.`identifier`, ':', -1))
WHERE COALESCE(NULLIF(u.`license`, ''), SUBSTRING_INDEX(u.`identifier`, ':', -1)) <> ''
ON DUPLICATE KEY UPDATE
  `account_id` = VALUES(`account_id`);

ALTER TABLE `fivenet_user` DROP FOREIGN KEY `fk_fivenet_user_account_id`;
ALTER TABLE `fivenet_user` DROP INDEX `idx_account_id`;
ALTER TABLE `fivenet_user` DROP COLUMN `account_id`;

COMMIT;
