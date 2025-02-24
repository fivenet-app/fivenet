BEGIN;

ALTER TABLE `fivenet_qualifications_exam_responses` ADD COLUMN `grading` longtext NULL;
ALTER TABLE `fivenet_qualifications_results` MODIFY COLUMN score decimal(5,1) NULL;

COMMIT;
