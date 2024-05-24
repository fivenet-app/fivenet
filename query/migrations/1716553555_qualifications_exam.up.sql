BEGIN;

-- Table: fivenet_qualifications_exam
CREATE TABLE
    IF NOT EXISTS `fivenet_qualifications_exam` (
        `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
        `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
        `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3),
        `deleted_at` datetime(3) DEFAULT NULL,
        `qualification_id` bigint(20) unsigned NOT NULL,
        `settings` longtext,
        `questions` longtext,
        PRIMARY KEY (`id`),
        UNIQUE KEY `idx_fivenet_qualifications_exam_quali_id_unique` (`qualification_id`),
        KEY `idx_fivenet_qualifications_exam_deleted_at` (`deleted_at`),
        CONSTRAINT `fk_fivenet_qualifications_exam_quali_id` FOREIGN KEY (`qualification_id`) REFERENCES `fivenet_qualifications` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
    ) ENGINE = InnoDB;

-- Table: fivenet_qualifications_exam_responses
CREATE TABLE
    IF NOT EXISTS `fivenet_qualifications_exam_responses` (
        `user_id` int(11) NOT NULL,
        `qualification_id` bigint(20) unsigned NOT NULL,
        `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
        `started_at` datetime(3) DEFAULT NULL,
        `ended_at` datetime(3) DEFAULT NULL,
        `responses` longtext,
        `closed` tinyint(1) DEFAULT 0,
        PRIMARY KEY (`user_id`, `qualification_id`),
        KEY `idx_fivenet_qualifications_exam_responses_closed` (`closed`),
        CONSTRAINT `fk_fivenet_qualifications_exam_responses_quali_id` FOREIGN KEY (`qualification_id`) REFERENCES `fivenet_qualifications` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
    ) ENGINE = InnoDB;

COMMIT;
