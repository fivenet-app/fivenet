BEGIN;

-- Table: fivenet_calendar_rsvp_occurrence
CREATE TABLE
    IF NOT EXISTS `fivenet_calendar_rsvp_occurrence` (
        `entry_id` bigint(20) unsigned NOT NULL,
        `occurrence_key` varchar(128) NOT NULL,
        `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
        `user_id` int(11) NOT NULL,
        `response` smallint(2) DEFAULT 0,
        PRIMARY KEY (`entry_id`, `occurrence_key`, `user_id`),
        KEY `idx_fivenet_calendar_rsvp_occurrence_response` (`entry_id`, `response`),
        KEY `idx_fivenet_calendar_rsvp_occurrence_user` (`user_id`),
        CONSTRAINT `fk_fivenet_calendar_rsvp_occurrence_entry_id` FOREIGN KEY (`entry_id`) REFERENCES `fivenet_calendar_entries` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
        CONSTRAINT `fk_fivenet_calendar_rsvp_occurrence_user_id` FOREIGN KEY (`user_id`) REFERENCES `{{.UsersTableName}}` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
    ) ENGINE = InnoDB;

COMMIT;
