BEGIN;

ALTER TABLE `fivenet_job_labels` MODIFY COLUMN `color` char(7) DEFAULT '#5c7aff' NULL;

UPDATE `fivenet_job_labels` SET `color` = '#5c7aff' WHERE `color` IS NULL;

COMMIT;
