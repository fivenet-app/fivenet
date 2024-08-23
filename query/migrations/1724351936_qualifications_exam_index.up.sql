BEGIN;

ALTER TABLE `fivenet_qualifications_exam_questions` DROP FOREIGN KEY `fk_fivenet_qualifications_exam_questions_quali_id`;
ALTER TABLE `fivenet_qualifications_exam_questions` DROP INDEX `fk_fivenet_qualifications_exam_questions_quali_id`;

ALTER TABLE `fivenet_qualifications_exam_questions` ADD KEY `fk_fivenet_qualifications_exam_questions_quali_id` (`qualification_id`);
ALTER TABLE `fivenet_qualifications_exam_questions` ADD CONSTRAINT `fk_fivenet_qualifications_exam_questions_quali_id` FOREIGN KEY (`qualification_id`) REFERENCES `fivenet_qualifications` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

COMMIT;
