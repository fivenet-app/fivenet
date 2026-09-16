BEGIN;

ALTER TABLE `fivenet_qualifications_exam_users`
    ADD COLUMN `attempt_id` CHAR(36) NULL AFTER `user_id`;

UPDATE `fivenet_qualifications_exam_users`
SET `attempt_id` = UUID()
WHERE `attempt_id` IS NULL;

ALTER TABLE `fivenet_qualifications_exam_users`
    MODIFY COLUMN `attempt_id` CHAR(36) NOT NULL,
    ADD UNIQUE KEY `uq_fivenet_qualifications_exam_users_attempt_id` (`attempt_id`);

ALTER TABLE `fivenet_qualifications_exam_responses`
    ADD COLUMN `attempt_id` CHAR(36) NULL AFTER `user_id`;

UPDATE `fivenet_qualifications_exam_responses` r
INNER JOIN `fivenet_qualifications_exam_users` u
    ON u.`qualification_id` = r.`qualification_id`
   AND u.`user_id` = r.`user_id`
SET r.`attempt_id` = u.`attempt_id`;

DELETE r
FROM `fivenet_qualifications_exam_responses` r
LEFT JOIN `fivenet_qualifications_exam_users` u
    ON u.`qualification_id` = r.`qualification_id`
   AND u.`user_id` = r.`user_id`
WHERE u.`attempt_id` IS NULL;

ALTER TABLE `fivenet_qualifications_exam_responses`
    MODIFY COLUMN `attempt_id` CHAR(36) NOT NULL,
    DROP PRIMARY KEY,
    ADD PRIMARY KEY (`qualification_id`, `user_id`, `attempt_id`),
    ADD KEY `idx_fivenet_qualifications_exam_responses_attempt_id` (`attempt_id`),
    ADD CONSTRAINT `fk_fivenet_qualifications_exam_responses_attempt_id`
        FOREIGN KEY (`attempt_id`)
        REFERENCES `fivenet_qualifications_exam_users` (`attempt_id`)
        ON DELETE CASCADE ON UPDATE CASCADE;

COMMIT;
