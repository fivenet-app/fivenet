BEGIN;

-- Audit-log job filters are commonly combined with the default newest-first
-- ordering. Replace the prefix-only index with one that can also serve that
-- order for a selected job.
ALTER TABLE `fivenet_audit_log`
  DROP INDEX `idx_fivenet_audit_log_user_job`,
  ADD INDEX `idx_fivenet_audit_log_user_job_created_at` (`user_job`, `created_at` DESC);

COMMIT;
