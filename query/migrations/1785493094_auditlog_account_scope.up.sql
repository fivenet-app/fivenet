BEGIN;

ALTER TABLE `fivenet_audit_log`
  MODIFY COLUMN `user_id` int(11) DEFAULT NULL,
  ADD COLUMN `account_id` bigint(20) unsigned DEFAULT NULL,
  ADD KEY `idx_fivenet_audit_log_account_id` (`account_id`);

ALTER TABLE `fivenet_audit_log`
  MODIFY COLUMN `account_id` bigint(20) unsigned DEFAULT NULL AFTER `user_job`;

COMMIT;
