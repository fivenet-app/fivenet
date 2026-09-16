BEGIN;

ALTER TABLE `fivenet_qualifications_activity`
    ADD COLUMN `attempt_id` CHAR(36) NULL AFTER `data`,
    ADD UNIQUE KEY `uq_fivenet_qualifications_activity_attempt_type` (`attempt_id`, `activity_type`);

COMMIT;
