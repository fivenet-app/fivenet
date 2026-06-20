BEGIN;

ALTER TABLE `fivenet_user` ADD COLUMN `account_id` bigint(20) unsigned AFTER `id`;
ALTER TABLE `fivenet_user` ADD KEY `idx_account_id` (`account_id`);
ALTER TABLE `fivenet_user` ADD CONSTRAINT `fk_fivenet_user_account_id` FOREIGN KEY (`account_id`) REFERENCES `fivenet_accounts` (`id`) ON DELETE SET NULL ON UPDATE CASCADE;

UPDATE `fivenet_user` u
INNER JOIN `fivenet_user_accounts` ua ON ua.`user_id` = u.`id`
SET u.`account_id` = ua.`account_id`;

DROP TABLE IF EXISTS `fivenet_user_accounts`;

COMMIT;
