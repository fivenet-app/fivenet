BEGIN;

-- Table: `fivenet_calendar_subs` - Remove entry_id column and fix indexes, if it exists
ALTER TABLE `fivenet_calendar_subs` DROP FOREIGN KEY `fk_fivenet_calendar_subs_entry_id`;
ALTER TABLE `fivenet_calendar_subs` DROP KEY `idx_fivenet_calendar_subs_unique`;
ALTER TABLE `fivenet_calendar_subs` DROP INDEX `idx_fivenet_calendar_subs_confirmed`;
ALTER TABLE `fivenet_calendar_subs` DROP COLUMN `entry_id`;

ALTER TABLE `fivenet_calendar_subs` ADD UNIQUE KEY `idx_fivenet_calendar_subs_unique` (`calendar_id`,`user_id`);

-- Table: `fivenet_accounts_oauth2` - Add columns to store tokens
ALTER TABLE `fivenet_accounts_oauth2`
  ADD COLUMN `access_token` VARCHAR(512),
  ADD COLUMN `refresh_token` VARCHAR(512),
  ADD COLUMN `token_type` VARCHAR(16),
  ADD COLUMN `scope` VARCHAR(128),
  ADD COLUMN `expires_in` INT,
  ADD COLUMN `obtained_at` DATETIME;

COMMIT;
