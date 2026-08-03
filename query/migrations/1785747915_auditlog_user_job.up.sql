BEGIN;

ALTER TABLE `fivenet_audit_log` MODIFY COLUMN `user_job` varchar(20) DEFAULT NULL;

UPDATE `fivenet_audit_log` SET `user_job` = NULL WHERE `user_job` = '';
UPDATE `fivenet_audit_log` SET `user_id` = NULL WHERE `user_id` = 0;

DELETE FROM `fivenet_audit_log` WHERE `service` = 'sync.SyncService' AND `user_id` IS NULL;

COMMIT;
