BEGIN;

-- Table: fivenet_centrum_settings

CREATE TABLE
    IF NOT EXISTS `fivenet_centrum_settings` (
        `job` varchar(50) NOT NULL,
        `enabled` tinyint(1) DEFAULT 0 NOT NULL,
        `mode` mediumint(2) DEFAULT 0 NULL,
        `fallback_mode` mediumint(2) DEFAULT 0 NULL,
        PRIMARY KEY (`job`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- Table: fivenet_centrum_users

CREATE TABLE
    IF NOT EXISTS `fivenet_centrum_users` (
        `job` varchar(50) NOT NULL,
        `user_id` int(11) NOT NULL,
        `identifier` varchar(64) NOT NULL,
        PRIMARY KEY (`user_id`),
        KEY `idx_fivenet_centrum_users_job` (`job`),
        CONSTRAINT `fk_fivenet_centrum_users_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
        CONSTRAINT `fk_fivenet_centrum_users_identifier` FOREIGN KEY (`identifier`) REFERENCES `fivenet_user_locations` (`identifier`) ON DELETE CASCADE ON UPDATE CASCADE
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- Table: fivenet_centrum_units

CREATE TABLE
    IF NOT EXISTS `fivenet_centrum_units` (
        `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
        `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
        `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3),
        `job` varchar(20) DEFAULT NULL,
        `name` varchar(128) NOT NULL,
        `initials` varchar(4) NOT NULL,
        `color` char(6) NOT NULL,
        `description` varchar(255) NULL DEFAULT NULL,
        PRIMARY KEY (`id`),
        UNIQUE KEY (`job`, `name`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- Table: fivenet_centrum_units_users

CREATE TABLE
    IF NOT EXISTS `fivenet_centrum_units_users` (
        `unit_id` bigint(20) unsigned NOT NULL,
        `user_id` int(11) NOT NULL,
        `identifier` varchar(64) NOT NULL,
        PRIMARY KEY (`unit_id`, `user_id`),
        KEY `idx_fivenet_centrum_units_users_unit_id` (`unit_id`),
        UNIQUE KEY `idx_fivenet_centrum_units_users_user_id` (`user_id`),
        CONSTRAINT `fk_fivenet_centrum_units_users_unit_id` FOREIGN KEY (`unit_id`) REFERENCES `fivenet_centrum_units` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
        CONSTRAINT `fk_fivenet_centrum_units_users_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
        CONSTRAINT `fk_fivenet_centrum_units_users_identifier` FOREIGN KEY (`identifier`) REFERENCES `fivenet_user_locations` (`identifier`) ON DELETE CASCADE ON UPDATE CASCADE
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- Table: fivenet_centrum_units_status

CREATE TABLE
    IF NOT EXISTS `fivenet_centrum_units_status` (
        `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
        `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
        `unit_id` bigint(20) unsigned NOT NULL,
        `status` smallint(2) NOT NULL,
        `reason` varchar(255) NULL DEFAULT NULL,
        `code` varchar(20) NULL DEFAULT NULL,
        `user_id` int(11) NULL DEFAULT NULL,
        `x` decimal(24, 14) DEFAULT NULL,
        `y` decimal(24, 14) DEFAULT NULL,
        `creator_id` int(11) NULL DEFAULT NULL,
        PRIMARY KEY (`id`),
        KEY `idx_fivenet_centrum_units_status_unit_id` (`unit_id`),
        KEY `idx_fivenet_centrum_units_status_user_id` (`user_id`),
        CONSTRAINT `fk_fivenet_centrum_units_status_unit_id` FOREIGN KEY (`unit_id`) REFERENCES `fivenet_centrum_units` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
        CONSTRAINT `fk_fivenet_centrum_units_status_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
        CONSTRAINT `fk_fivenet_centrum_units_status_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
    );

-- Table: fivenet_centrum_dispatches

CREATE TABLE
    IF NOT EXISTS `fivenet_centrum_dispatches` (
        `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
        `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
        `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3),
        `job` varchar(20) DEFAULT NULL,
        `message` varchar(255) NOT NULL,
        `description` varchar(1024) NULL DEFAULT NULL,
        `attributes` varchar(2048) NULL DEFAULT NULL,
        `x` decimal(24, 14) DEFAULT NULL,
        `y` decimal(24, 14) DEFAULT NULL,
        `anon` tinyint(1) DEFAULT 0,
        `user_id` int(11) NOT NULL,
        PRIMARY KEY (`id`),
        KEY `idx_fivenet_centrum_dispatches_job` (`job`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- Table: fivenet_centrum_dispatches_asgmts

CREATE TABLE
    IF NOT EXISTS `fivenet_centrum_dispatches_asgmts` (
        `dispatch_id` bigint(20) unsigned NOT NULL,
        `unit_id` bigint(20) unsigned NOT NULL,
        `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
        `expires_at` datetime(3) NULL DEFAULT NULL,
        PRIMARY KEY (`dispatch_id`, `unit_id`),
        KEY `idx_fivenet_centrum_dispatches_asgmts_dispatch_id` (`dispatch_id`),
        KEY `idx_fivenet_centrum_dispatches_asgmts_unit_id` (`unit_id`),
        CONSTRAINT `fk_fivenet_centrum_dispatches_asgmts_dispatch_id` FOREIGN KEY (`dispatch_id`) REFERENCES `fivenet_centrum_dispatches` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
        CONSTRAINT `fk_fivenet_centrum_dispatches_asgmts_unit_id` FOREIGN KEY (`unit_id`) REFERENCES `fivenet_centrum_units` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- Table: fivenet_centrum_dispatches_status

CREATE TABLE
    IF NOT EXISTS `fivenet_centrum_dispatches_status` (
        `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
        `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
        `dispatch_id` bigint(20) unsigned NOT NULL,
        `unit_id` bigint(20) unsigned NULL DEFAULT NULL,
        `status` smallint(2) NOT NULL,
        `reason` varchar(255) DEFAULT NULL,
        `code` varchar(20) NULL DEFAULT NULL,
        `x` decimal(24, 14) DEFAULT NULL,
        `y` decimal(24, 14) DEFAULT NULL,
        `user_id` int(11) NULL DEFAULT NULL,
        PRIMARY KEY (`id`),
        KEY `idx_fivenet_centrum_dispatches_status_dispatch_id` (`dispatch_id`),
        KEY `idx_fivenet_centrum_dispatches_status_status` (`status`),
        CONSTRAINT `fk_fivenet_centrum_dispatches_status_dispatch_id` FOREIGN KEY (`dispatch_id`) REFERENCES `fivenet_centrum_dispatches` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
        CONSTRAINT `fk_fivenet_centrum_dispatches_status_unit_id` FOREIGN KEY (`unit_id`) REFERENCES `fivenet_centrum_units` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
        CONSTRAINT `fk_fivenet_centrum_dispatches_status_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- Table: fivenet_centrum_markers

CREATE TABLE
    IF NOT EXISTS `fivenet_centrum_markers` (
        `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
        `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
        `job` varchar(20) DEFAULT NULL,
        `name` varchar(255) NOT NULL,
        `description` varchar(1024) NULL DEFAULT NULL,
        `x` decimal(24, 14) DEFAULT NULL,
        `y` decimal(24, 14) DEFAULT NULL,
        `color` char(6) DEFAULT "FFFC31",
        `icon` varchar(128) DEFAULT NULL,
        `marker_type` smallint(2) NOT NULL,
        `marker_data` longtext DEFAULT NULL,
        `creator_id` int(11) NULL DEFAULT NULL,
        PRIMARY KEY (`id`),
        KEY `idx_fivenet_centrum_markers_job` (`job`),
        CONSTRAINT `fk_fivenet_centrum_markers_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- Table: fivenet_centrum_dispatches_attrs

-- TODO

-- Table: fivenet_centrum_codes

-- TODO

COMMIT;
