BEGIN;

UPDATE `fivenet_audit_log` SET `user_job` = '' WHERE `user_job` IS NULL;

ALTER TABLE `fivenet_audit_log` MODIFY COLUMN `user_job` varchar(20) NOT NULL;

COMMIT;
