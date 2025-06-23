BEGIN;

-- Table: fivenet_centrum_settings - Add `access` column
ALTER TABLE `fivenet_centrum_settings` ADD COLUMN `access` text AFTER `timings`;

COMMIT;
