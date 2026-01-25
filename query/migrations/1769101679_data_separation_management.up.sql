BEGIN;

ALTER TABLE `fivenet_user` ADD KEY `idx_phone_number` (`phone_number`); -- Add index on phone_number

-- Table: `fivenet_user_licenses`
-- Drop dependent constraints and indexes
ALTER TABLE `fivenet_user_licenses` DROP CONSTRAINT `fk_fivenet_user_licenses_owner`;
ALTER TABLE `fivenet_user_licenses` DROP INDEX `fk_fivenet_user_licenses_owner`;

ALTER TABLE `fivenet_user_licenses` DROP CONSTRAINT `fk_fivenet_user_licenses_type`;
ALTER TABLE `fivenet_user_licenses` DROP INDEX `PRIMARY`;

-- Table: `fivenet_user`
-- Modify primary key on fivenet_user to be `id` instead of `identifier`
-- `identifier` column was primary key
ALTER TABLE `fivenet_user` DROP INDEX `PRIMARY`;

-- Set `id` as primary key
ALTER TABLE `fivenet_user` ADD PRIMARY KEY (`id`);
-- `id` column was unique key
ALTER TABLE `fivenet_user` DROP INDEX `idx_id`;

-- Set `identifier` as unique key
ALTER TABLE `fivenet_user` ADD UNIQUE KEY `idx_identifier` (`identifier`);

-- Add `account_id` column
ALTER TABLE `fivenet_user` ADD COLUMN `account_id` bigint(20) unsigned AFTER `id`;
-- Add index on `account_id` (can't be unique as multiple users can belong to same account)
ALTER TABLE `fivenet_user` ADD KEY `idx_account_id` (`account_id`);
-- Add foreign key constraint to `fivenet_account`
ALTER TABLE `fivenet_user` ADD CONSTRAINT `fk_fivenet_user_account_id` FOREIGN KEY (`account_id`) REFERENCES `fivenet_accounts` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- Add `deleted_at` column for soft deletes
ALTER TABLE `fivenet_user` ADD COLUMN `deleted_at` datetime(3) DEFAULT NULL AFTER `last_seen`;
ALTER TABLE `fivenet_user` ADD COLUMN `deleted_reason` VARCHAR(255) DEFAULT NULL AFTER `deleted_at`;

-- Fill in `account_id` for existing users with "matching" accounts
UPDATE `fivenet_user` u
JOIN `fivenet_accounts` a ON SUBSTRING_INDEX(u.`identifier`, ':', -1) = a.`license`
SET u.`account_id` = a.`id`;

-- Table: `fivenet_user_licenses`
-- Create new `user_id` column to replace `owner` in fivenet_user_licenses
ALTER TABLE `fivenet_user_licenses` ADD COLUMN `user_id` int(11) NOT NULL FIRST;
ALTER TABLE `fivenet_user_licenses` ADD PRIMARY KEY (`type`, `user_id`);

-- Migrate `owner` (identifier) to `user_id`
UPDATE `fivenet_user_licenses` ul
JOIN `fivenet_user` u ON ul.`owner` = u.`identifier`
SET ul.`user_id` = u.`id`, ul.`owner` = NULL;

-- Remove "broken" records where no matching user was found..
DELETE FROM `fivenet_user_licenses` WHERE `user_id` = 0;
ALTER TABLE `fivenet_user_licenses` ADD CONSTRAINT `fk_fivenet_user_licenses_user_id` FOREIGN KEY (`user_id`) REFERENCES `fivenet_user`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- Recreate foreign key for `type`
ALTER TABLE `fivenet_user_licenses` ADD CONSTRAINT `fk_fivenet_user_licenses_type` FOREIGN KEY (`type`) REFERENCES `fivenet_licenses` (`type`) ON DELETE CASCADE ON UPDATE CASCADE;

-- Table: `fivenet_owned_vehicles`
-- Drop dependent constraints and indexes
ALTER TABLE `fivenet_owned_vehicles` DROP INDEX `idx_fivenet_owned_vehicles_ownerplate`;
ALTER TABLE `fivenet_owned_vehicles` DROP INDEX `idx_fivenet_owned_vehicles_owner`;
ALTER TABLE `fivenet_owned_vehicles` DROP INDEX `idx_fivenet_owned_vehicles_owner_type`;
ALTER TABLE `fivenet_owned_vehicles` DROP INDEX `idx_fivenet_owned_vehicles_owner_model_type`;

ALTER TABLE `fivenet_owned_vehicles` ADD COLUMN `user_id` int(11) NULL AFTER `owner`;

-- Migrate `owner` (identifier) to `user_id`
UPDATE `fivenet_owned_vehicles` owv
JOIN `fivenet_user` u ON owv.`owner` = u.`identifier`
SET owv.`user_id` = u.`id`, owv.`owner` = NULL;

-- Remove "broken" records where no matching user was found..
DELETE FROM `fivenet_owned_vehicles` WHERE `user_id` = 0;
ALTER TABLE `fivenet_owned_vehicles` ADD CONSTRAINT `fk_fivenet_owned_vehicles_user_id` FOREIGN KEY (`user_id`) REFERENCES `fivenet_user`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- Recreate indexes with `user_id` instead of `owner`
ALTER TABLE `fivenet_owned_vehicles` ADD UNIQUE KEY `idx_fivenet_owned_vehicles_userplate` (`user_id`, `plate`);
ALTER TABLE `fivenet_owned_vehicles` ADD UNIQUE KEY `idx_fivenet_owned_vehicles_jobplate` (`job`, `plate`);
ALTER TABLE `fivenet_owned_vehicles` ADD KEY `idx_fivenet_owned_vehicles_user` (`user_id`);
ALTER TABLE `fivenet_owned_vehicles` ADD KEY `idx_fivenet_owned_vehicles_user_type` (`user_id`, `type`);
ALTER TABLE `fivenet_owned_vehicles` ADD KEY `idx_fivenet_owned_vehicles_user_model_type` (`user_id`, `model`, `type`);

-- Table: `fivenet_user_phone_numbers`
CREATE TABLE IF NOT EXISTS `fivenet_user_phone_numbers` (
  `user_id` int(11) NOT NULL,
  `phone_number` varchar(15) NOT NULL,
  `is_primary` tinyint(1) NOT NULL DEFAULT '0',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`user_id`, `phone_number`),
  UNIQUE KEY `idx_phone_number` (`phone_number`),
  KEY `idx_is_primary` (`is_primary`),
  CONSTRAINT `fk_fivenet_user_phone_numbers_user_id` FOREIGN KEY (`user_id`) REFERENCES `fivenet_user`(`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Add all current `phone_number`s as primary numbers in `fivenet_user_phone_numbers`
INSERT INTO `fivenet_user_phone_numbers` (`user_id`, `phone_number`, `is_primary`)
SELECT u.`id` AS `user_id`, u.`phone_number` AS `phone_number`, 1 AS `is_primary`
FROM `fivenet_user` u
WHERE u.`phone_number` IS NOT NULL AND u.`phone_number` != '';

-- Table: `fivenet_user_jobs`
CREATE TABLE IF NOT EXISTS `fivenet_user_jobs` (
  `user_id` int(11) NOT NULL,
  `job` varchar(50) NOT NULL,
  `grade` int(11) NOT NULL DEFAULT '0',
  `is_primary` tinyint(1) NOT NULL DEFAULT '0',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`user_id`, `job`),
  KEY `idx_job` (`job`),
  KEY `idx_is_primary` (`is_primary`),
  CONSTRAINT `fk_fivenet_user_jobs_user_id` FOREIGN KEY (`user_id`) REFERENCES `fivenet_user`(`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Add all current `job`s as primary job in `fivenet_user_jobs`
INSERT INTO `fivenet_user_jobs` (`user_id`, `job`, `grade`, `is_primary`)
SELECT u.`id` AS `user_id`, u.`job` AS `job`,  u.`job_grade` AS `grade`, 1 AS `is_primary`
FROM `fivenet_user` u
WHERE u.`job` IS NOT NULL AND u.`job` != '';

-- Drop old `owner` column from both tables
ALTER TABLE `fivenet_user_licenses` DROP COLUMN `owner`;
ALTER TABLE `fivenet_owned_vehicles` DROP COLUMN `owner`;

COMMIT;
