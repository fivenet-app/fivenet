BEGIN;

ALTER TABLE `fivenet_qualifications_activity`
    DROP KEY `uq_fivenet_qualifications_activity_attempt_type`,
    DROP COLUMN `attempt_id`;

COMMIT;
