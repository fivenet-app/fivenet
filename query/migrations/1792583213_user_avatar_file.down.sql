BEGIN;

ALTER TABLE `fivenet_user_props`
  ADD COLUMN `avatar_file_id` bigint unsigned DEFAULT NULL AFTER `avatar`;

ALTER TABLE `fivenet_user_props`
  ADD CONSTRAINT `fk_fivenet_user_props_avatar_file_id` FOREIGN KEY (`avatar_file_id`)
    REFERENCES `fivenet_files` (`id`);

UPDATE `fivenet_user_props` p
INNER JOIN `fivenet_user` u ON u.`id` = p.`user_id`
SET p.`avatar_file_id` = u.`avatar_file_id`
WHERE u.`avatar_file_id` IS NOT NULL;

INSERT INTO `fivenet_user_files` (`user_id`, `file_id`, `kind`)
SELECT `id`, `avatar_file_id`, 'avatar'
FROM `fivenet_user`
WHERE `avatar_file_id` IS NOT NULL
ON DUPLICATE KEY UPDATE `file_id` = VALUES(`file_id`);

ALTER TABLE `fivenet_user`
  DROP FOREIGN KEY `fk_fivenet_user_avatar_file_id`,
  DROP INDEX `idx_avatar_file_id`,
  DROP COLUMN `avatar_file_id`;

COMMIT;
