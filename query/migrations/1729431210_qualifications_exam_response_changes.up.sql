BEGIN;

TRUNCATE `fivenet_qualifications_exam_responses`;
ALTER TABLE `fivenet_qualifications_exam_responses` DROP FOREIGN KEY `fk_fivenet_qualifications_exam_responses_question_id`;
ALTER TABLE `fivenet_qualifications_exam_responses` DROP INDEX `PRIMARY`;
ALTER TABLE `fivenet_qualifications_exam_responses` DROP `question_id`;
ALTER TABLE `fivenet_qualifications_exam_responses` ADD `qualification_id` BIGINT UNSIGNED NOT NULL;
ALTER TABLE `fivenet_qualifications_exam_responses` CHANGE `qualification_id` `qualification_id` BIGINT UNSIGNED NOT NULL FIRST;
ALTER TABLE `fivenet_qualifications_exam_responses` CHANGE `response` `responses` longtext NULL;
ALTER TABLE `fivenet_qualifications_exam_responses` ADD PRIMARY KEY (`qualification_id`, `user_id`);
ALTER TABLE `fivenet_qualifications_exam_responses` ADD CONSTRAINT `fk_fivenet_qualifications_exam_responses_quali_user_id` FOREIGN KEY (`qualification_id`) REFERENCES `fivenet_qualifications` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_qualifications_exam_responses` ADD CONSTRAINT `fk_fivenet_qualifications_exam_responses_user_id` FOREIGN KEY (`user_id`) REFERENCES `{{.UsersTableName}}` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

COMMIT;
