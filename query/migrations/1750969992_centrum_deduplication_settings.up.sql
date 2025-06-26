BEGIN;

-- Table: fivenet_centrum_settings - Add `configuration` column
ALTER TABLE `fivenet_centrum_settings` ADD COLUMN `configuration` text AFTER `access`;

COMMIT;
