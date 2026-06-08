BEGIN;

ALTER TABLE `fivenet_calendar`
    ADD COLUMN `system_kind` smallint(2) DEFAULT NULL AFTER `job`;

ALTER TABLE `fivenet_calendar`
    ADD UNIQUE KEY `idx_fivenet_calendar_job_system_kind` (`job`, `system_kind`);

COMMIT;
