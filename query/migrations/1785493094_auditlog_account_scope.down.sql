BEGIN;

UPDATE `fivenet_audit_log` SET `user_id` = 0 WHERE `user_id` IS NULL;

ALTER TABLE `fivenet_audit_log`
  DROP INDEX `idx_fivenet_audit_log_account_id`,
  DROP COLUMN `account_id`,
  MODIFY COLUMN `user_id` int(11) NOT NULL;

COMMIT;
