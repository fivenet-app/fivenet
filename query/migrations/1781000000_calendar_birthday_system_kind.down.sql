BEGIN;

ALTER TABLE `fivenet_calendar`
    DROP INDEX `idx_fivenet_calendar_job_system_kind`;

ALTER TABLE `fivenet_calendar`
    DROP COLUMN `system_kind`;

COMMIT;
