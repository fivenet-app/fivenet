BEGIN;

ALTER TABLE `fivenet_qualifications_requests`
    ADD COLUMN `exam_attempt_id` CHAR(36) NULL AFTER `approver_job`,
    ADD KEY `idx_fivenet_qualifications_requests_exam_attempt_id` (`exam_attempt_id`),
    ADD CONSTRAINT `fk_fivenet_qualifications_requests_exam_attempt_id`
        FOREIGN KEY (`exam_attempt_id`)
        REFERENCES `fivenet_qualifications_exam_users` (`attempt_id`)
        ON DELETE SET NULL ON UPDATE CASCADE;

ALTER TABLE `fivenet_qualifications_results`
    ADD COLUMN `exam_attempt_id` CHAR(36) NULL AFTER `auto_graded`,
    ADD KEY `idx_fivenet_qualifications_results_exam_attempt_id` (`exam_attempt_id`),
    ADD CONSTRAINT `fk_fivenet_qualifications_results_exam_attempt_id`
        FOREIGN KEY (`exam_attempt_id`)
        REFERENCES `fivenet_qualifications_exam_users` (`attempt_id`)
        ON DELETE SET NULL ON UPDATE CASCADE;

UPDATE `fivenet_qualifications_requests` r
INNER JOIN `fivenet_qualifications_exam_users` e
    ON e.`qualification_id` = r.`qualification_id`
   AND e.`user_id` = r.`user_id`
SET r.`exam_attempt_id` = e.`attempt_id`
WHERE r.`deleted_at` IS NULL;

COMMIT;
