BEGIN;

-- Table: `fivenet_qualifications_exam_questions` - Add `order` column
ALTER TABLE `fivenet_qualifications_exam_questions` ADD `order` int(11) DEFAULT 0 AFTER `points`;

COMMIT;
