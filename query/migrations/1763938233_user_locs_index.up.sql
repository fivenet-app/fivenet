BEGIN;

-- Table: `fivenet_centrum_user_locations`
ALTER TABLE `fivenet_centrum_user_locations`
    ADD INDEX `idx_updated_at` (`updated_at`),
    ADD INDEX `idx_job_job_grade` (`job`, `job_grade`);

COMMIT;
