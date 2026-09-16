BEGIN;

ALTER TABLE `fivenet_qualifications_exam_responses`
    DROP FOREIGN KEY `fk_fivenet_qualifications_exam_responses_attempt_id`,
    DROP KEY `idx_fivenet_qualifications_exam_responses_attempt_id`,
    DROP PRIMARY KEY,
    ADD PRIMARY KEY (`qualification_id`, `user_id`),
    DROP COLUMN `attempt_id`;

ALTER TABLE `fivenet_qualifications_exam_users`
    DROP KEY `uq_fivenet_qualifications_exam_users_attempt_id`,
    DROP COLUMN `attempt_id`;

COMMIT;
