BEGIN;

ALTER TABLE `fivenet_user`
  ADD COLUMN `avatar_file_id` bigint unsigned DEFAULT NULL AFTER `phone_number`,
  ADD KEY `idx_avatar_file_id` (`avatar_file_id`),
  ADD CONSTRAINT `fk_fivenet_user_avatar_file_id` FOREIGN KEY (`avatar_file_id`)
    REFERENCES `fivenet_files` (`id`) ON DELETE SET NULL;

-- Avatar files are identity/profile data, not investigation attachments.
UPDATE `fivenet_user` u
INNER JOIN `fivenet_user_files` uf
  ON uf.`user_id` = u.`id` AND uf.`kind` = 'avatar'
SET u.`avatar_file_id` = uf.`file_id`
WHERE u.`avatar_file_id` IS NULL;

DELETE FROM `fivenet_user_files` WHERE `kind` = 'avatar';

ALTER TABLE `fivenet_user_props`
  DROP FOREIGN KEY `fk_fivenet_user_props_avatar_file_id`,
  DROP COLUMN `avatar_file_id`;

COMMIT;
