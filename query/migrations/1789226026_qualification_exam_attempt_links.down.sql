BEGIN;

ALTER TABLE `fivenet_qualifications_results`
    DROP FOREIGN KEY `fk_fivenet_qualifications_results_exam_attempt_id`,
    DROP KEY `idx_fivenet_qualifications_results_exam_attempt_id`,
    DROP COLUMN `exam_attempt_id`;

ALTER TABLE `fivenet_qualifications_requests`
    DROP FOREIGN KEY `fk_fivenet_qualifications_requests_exam_attempt_id`,
    DROP KEY `idx_fivenet_qualifications_requests_exam_attempt_id`,
    DROP COLUMN `exam_attempt_id`;

COMMIT;
