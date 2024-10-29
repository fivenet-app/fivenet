BEGIN;

ALTER TABLE `fivenet_qualifications_exam_questions` ADD `points` int(11) DEFAULT 0 NULL;

COMMIT;
