BEGIN;
-- Table: fivenet_centrum_settings - Add `access` column
ALTER TABLE `fivenet_centrum_settings` ADD COLUMN `access` text AFTER `timings`;

-- Add job columns
ALTER TABLE `fivenet_centrum_units_users` ADD COLUMN `unit_job` varchar(20) NOT NULL AFTER `unit_id`;
ALTER TABLE `fivenet_centrum_units_users` ADD COLUMN `user_job` varchar(20) NOT NULL AFTER `user_id`;

ALTER TABLE `fivenet_centrum_units_status` ADD COLUMN `unit_job` varchar(20) NOT NULL AFTER `unit_id`;
ALTER TABLE `fivenet_centrum_units_status` ADD COLUMN `user_job` varchar(20) AFTER `user_id`;
ALTER TABLE `fivenet_centrum_units_status` ADD COLUMN `creator_job` varchar(20) AFTER `creator_id`;

ALTER TABLE `fivenet_centrum_dispatches_asgmts` ADD COLUMN `dispatch_job` varchar(20) NOT NULL AFTER `dispatch_id`;
ALTER TABLE `fivenet_centrum_dispatches_asgmts` ADD COLUMN `unit_job` varchar(20) NOT NULL AFTER `unit_id`;

ALTER TABLE `fivenet_centrum_dispatches_status` ADD COLUMN `dispatch_job` varchar(20) NOT NULL AFTER `dispatch_id`;
ALTER TABLE `fivenet_centrum_dispatches_status` ADD COLUMN `unit_job` varchar(20) AFTER `unit_id`;
ALTER TABLE `fivenet_centrum_dispatches_status` ADD COLUMN `user_job` varchar(20) AFTER `user_id`;

-- Table: `fivenet_centrum_users` - Drop `identifier` column
ALTER TABLE `fivenet_centrum_users` DROP FOREIGN KEY `fk_fivenet_centrum_users_identifier`;
ALTER TABLE `fivenet_centrum_users` DROP COLUMN `identifier`;

COMMIT;
