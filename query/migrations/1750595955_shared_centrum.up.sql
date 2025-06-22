BEGIN;

ALTER TABLE `fivenet_centrum_dispatches_status` ADD COLUMN `creator_job` varchar(50) NULL;
ALTER TABLE `fivenet_centrum_units_status` ADD COLUMN `creator_job` varchar(50) NULL;

ALTER TABLE `fivenet_centrum_dispatches` CHANGE `job` `jobs` varchar(255) NOT NULL;

UPDATE `fivenet_centrum_dispatches` SET `jobs` = CONCAT('["', `jobs`, '"]') WHERE `jobs` NOT LIKE '[%]';

COMMIT;
