BEGIN;

-- Table: fivenet_job_conduct_files
CREATE TABLE IF NOT EXISTS `fivenet_job_conduct_files` (
    `conduct_id` bigint unsigned NOT NULL,
    `file_id` bigint unsigned NOT NULL,
    PRIMARY KEY (`conduct_id`, `file_id`),
    KEY `idx_file_id` (`file_id`),
    CONSTRAINT `fk_fivenet_job_conduct_files_conduct_id` FOREIGN KEY (`conduct_id`) REFERENCES `fivenet_job_conduct` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_fivenet_job_conduct_files_file_id` FOREIGN KEY (`file_id`) REFERENCES `fivenet_files` (`id`) ON DELETE RESTRICT
);

ALTER TABLE `fivenet_job_conduct` ADD COLUMN `draft` tinyint(1) DEFAULT '0';
ALTER TABLE `fivenet_job_conduct` CHANGE `draft` `draft` tinyint(1) DEFAULT '0' AFTER `type`;

ALTER TABLE `fivenet_job_conduct` ADD INDEX `idx_draft` (`draft`);

COMMIT;
