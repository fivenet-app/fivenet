BEGIN;

-- Table: fivenet_jobs_conduct

CREATE TABLE
    IF NOT EXISTS `fivenet_jobs_conduct` (
        `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
        `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
        `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3),
        `job` varchar(20) DEFAULT NULL,
        `type` smallint(2) NOT NULL,
        `message` longtext,
        `expires_at` datetime(3) DEFAULT NULL,
        `target_user_id` int(11) NULL DEFAULT NULL,
        `creator_id` int(11) NULL DEFAULT NULL,
        PRIMARY KEY (`id`),
        KEY (`type`),
        KEY (`created_at`),
        KEY (`target_user_id`),
        CONSTRAINT `fk_fivenet_jobs_conduct_target_user_id` FOREIGN KEY (`target_user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
        CONSTRAINT `fk_fivenet_jobs_conduct_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

COMMIT;
