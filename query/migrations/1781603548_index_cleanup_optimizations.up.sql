BEGIN;

-- INDEX CLEANUP AND OPTIMIZATION CHANGES

-- fivenet_documents: High-volume table with frequent state/draft filtering

-- Remove redundant single-column index (covered by composite index)
ALTER TABLE `fivenet_documents`
DROP INDEX `idx_draft`;

-- Add composite indexes for common query patterns
-- Pattern 1: WHERE state=? AND deleted_at IS NULL AND closed=?
ALTER TABLE `fivenet_documents`
ADD INDEX `idx_state_deleted_closed` (`state`, `deleted_at`, `closed`);

-- Pattern 2: WHERE creator_id=? AND state=? AND deleted_at IS NULL
ALTER TABLE `fivenet_documents`
ADD INDEX `idx_creator_state_deleted` (`creator_id`, `state`, `deleted_at`);

-- Pattern 3: WHERE category_id=? AND deleted_at IS NULL AND draft=?
ALTER TABLE `fivenet_documents`
ADD INDEX `idx_category_deleted_draft` (`category_id`, `deleted_at`, `draft`);

ALTER TABLE `fivenet_documents`
ADD INDEX `idx_fivenet_documents_public_deleted_updated_id` (`public`, `deleted_at`, `updated_at` DESC, `id`);

-- fivenet_audit_log: High-volume append-only log table

-- Remove low-selectivity index (rarely used in filters)
ALTER TABLE `fivenet_audit_log`
DROP INDEX `idx_fivenet_audit_log_target_user_job`;

-- Remove single-column index (covered by composite below)
ALTER TABLE `fivenet_audit_log`
DROP INDEX `idx_action`;

-- Add composite indexes for common audit query patterns
-- Pattern 1: WHERE user_id IN (...) AND created_at BETWEEN ? AND ? + pagination
ALTER TABLE `fivenet_audit_log`
ADD INDEX `idx_user_created_action` (`user_id`, `created_at` DESC, `action`);

-- Pattern 2: WHERE service IN (...) AND method IN (...) AND created_at
ALTER TABLE `fivenet_audit_log`
ADD INDEX `idx_service_method_created` (`service`, `method`, `created_at` DESC);

-- Pattern 3: WHERE created_at DESC + result filtering for reports
ALTER TABLE `fivenet_audit_log`
ADD INDEX `idx_created_result` (`created_at` DESC, `result`);

-- fivenet_calendar_entries: Time-series table with time-range queries

-- Remove redundant index (covered by calendar_start_end composite)
ALTER TABLE `fivenet_calendar_entries`
DROP INDEX `idx_fivenet_calendar_entries_times`;

-- Remove low-selectivity single-column index
ALTER TABLE `fivenet_calendar_entries`
DROP INDEX `idx_fivenet_calendar_entries_job`;

-- Add composite indexes for calendar query patterns
-- Pattern 1: WHERE calendar_id=? AND deleted_at IS NULL AND closed=?
ALTER TABLE `fivenet_calendar_entries`
ADD INDEX `idx_calendar_deleted_closed` (`calendar_id`, `deleted_at`, `closed`);

-- Pattern 2: WHERE start_time >= ? AND start_time <= ? AND calendar_id=?
ALTER TABLE `fivenet_calendar_entries`
ADD INDEX `idx_start_calendar_deleted` (`start_time`, `calendar_id`, `deleted_at`);

-- fivenet_centrum_dispatches: Dispatch table

-- Add composite indexes for dispatch query patterns
-- Pattern 1: WHERE creator_id=? ORDER BY created_at DESC
ALTER TABLE `fivenet_centrum_dispatches`
ADD INDEX `idx_creator_created` (`creator_id`, `created_at` DESC);

-- Pattern 2: WHERE postal=? ORDER BY created_at DESC
ALTER TABLE `fivenet_centrum_dispatches`
ADD INDEX `idx_postal_created` (`postal`, `created_at` DESC);

-- Pattern 3: WHERE anon=? ORDER BY created_at DESC
ALTER TABLE `fivenet_centrum_dispatches`
ADD INDEX `idx_anon_created` (`anon`, `created_at` DESC);

-- fivenet_documents_activity: Activity log table

-- Add composite index for activity type + document lookups
-- Pattern: WHERE document_id=? AND activity_type=? ORDER BY created_at DESC
ALTER TABLE `fivenet_documents_activity`
ADD INDEX `idx_document_type_created` (`document_id`, `activity_type`, `created_at` DESC);

-- fivenet_documents_approval_tasks: Approval task queries

-- Add index for approval lookup queries
-- Pattern: WHERE document_id=? AND snapshot_date=? AND status=?
ALTER TABLE `fivenet_documents_approval_tasks`
ADD INDEX `idx_doc_snapshot_status` (`document_id`, `snapshot_date`, `status`);

-- fivenet_calendar_entries: Additional pattern for recurring queries

-- Improve recurring/recurring_until queries
-- Pattern: WHERE calendar_id=? AND recurring_until >= ? AND start_time <= ?
ALTER TABLE `fivenet_calendar_entries`
ADD INDEX `idx_calendar_recurring_times` (`calendar_id`, `recurring_until`, `start_time`);

COMMIT;
