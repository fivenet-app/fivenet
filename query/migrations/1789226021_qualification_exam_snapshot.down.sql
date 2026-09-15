BEGIN;

ALTER TABLE `fivenet_qualifications_exam_users`
    DROP COLUMN `snapshot`;

COMMIT;
