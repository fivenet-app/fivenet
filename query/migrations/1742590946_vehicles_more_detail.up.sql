BEGIN;

ALTER TABLE `fivenet_owned_vehicles` ADD COLUMN `job` varchar(40) NULL AFTER `owner`;
ALTER TABLE `fivenet_owned_vehicles` ADD COLUMN `data` text DEFAULT NULL;

COMMIT;
