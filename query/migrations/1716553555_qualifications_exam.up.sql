BEGIN;

-- Table: fivenet_qualifications_exam_questions
CREATE TABLE
    IF NOT EXISTS `fivenet_qualifications_exam_questions` (
        `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
        `qualification_id` bigint(20) unsigned NOT NULL,
        `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
        `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3),
        `title` varchar(512) NOT NULL,
        `description` varchar(1024) DEFAULT NULL,
        `data` longtext,
        `answer` longtext,
        PRIMARY KEY (`id`),
        CONSTRAINT `fk_fivenet_qualifications_exam_questions_quali_id` FOREIGN KEY (`qualification_id`) REFERENCES `fivenet_qualifications` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
    ) ENGINE = InnoDB;

-- Table: fivenet_qualifications_exam_users
CREATE TABLE
    IF NOT EXISTS `fivenet_qualifications_exam_users` (
        `qualification_id` bigint(20) unsigned NOT NULL,
        `user_id` int(11) NOT NULL,
        `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
        `started_at` datetime(3) DEFAULT NULL,
        `ends_at` datetime(3) DEFAULT NULL,
        `ended_at` datetime(3) DEFAULT NULL,
        PRIMARY KEY (`qualification_id`, `user_id`),
        CONSTRAINT `fk_fivenet_qualifications_exam_users_quali_id` FOREIGN KEY (`qualification_id`) REFERENCES `fivenet_qualifications` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
        CONSTRAINT `fk_fivenet_qualifications_exam_users_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
    ) ENGINE = InnoDB;

-- Table: fivenet_qualifications_exam_responses
CREATE TABLE
    IF NOT EXISTS `fivenet_qualifications_exam_responses` (
        `question_id` bigint(20) unsigned NOT NULL,
        `user_id` int(11) NOT NULL,
        `response` longtext,
        PRIMARY KEY (`question_id`, `user_id`),
        CONSTRAINT `fk_fivenet_qualifications_exam_responses_question_id` FOREIGN KEY (`question_id`) REFERENCES `fivenet_qualifications_exam_questions` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
    ) ENGINE = InnoDB;

COMMIT;
