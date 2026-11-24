BEGIN;

ALTER TABLE `fivenet_audit_log`
  DROP INDEX `idx_fivenet_audit_log_user_job_created_at`,
  ADD INDEX `idx_fivenet_audit_log_user_job` (`user_job`);

COMMIT;
