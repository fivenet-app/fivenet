BEGIN;

-- Table: fivenet_qualifications_exam_questions_files
CREATE TABLE IF NOT EXISTS `fivenet_qualifications_exam_questions_files` (
    `question_id` bigint unsigned NOT NULL,
    `file_id` bigint unsigned NOT NULL,
    PRIMARY KEY (`question_id`, `file_id`),
    KEY `idx_file_id` (`file_id`),
    CONSTRAINT `fk_fivenet_qualifications_exam_questions_files_question_id` FOREIGN KEY (`question_id`) REFERENCES `fivenet_qualifications_exam_questions` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_fivenet_qualifications_exam_questions_files_file_id` FOREIGN KEY (`file_id`) REFERENCES `fivenet_files` (`id`) ON DELETE RESTRICT
);

-- Table: `fivenet_qualifications` - Add draft field
ALTER TABLE `fivenet_qualifications` ADD COLUMN `draft` tinyint(1) DEFAULT '0';
ALTER TABLE `fivenet_qualifications` CHANGE `draft` `draft` tinyint(1) DEFAULT '0' AFTER `closed`;

ALTER TABLE `fivenet_qualifications` ADD INDEX `idx_draft` (`draft`);

COMMIT;
