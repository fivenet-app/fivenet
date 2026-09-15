BEGIN;

ALTER TABLE `fivenet_qualifications_exam_users`
    ADD COLUMN `snapshot` LONGTEXT NULL AFTER `ended_at`;

COMMIT;
