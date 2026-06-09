BEGIN;

ALTER TABLE `fivenet_calendar`
    ADD COLUMN `discord_settings` longtext DEFAULT NULL AFTER `system_kind`;

CREATE TABLE IF NOT EXISTS `fivenet_calendar_discord_reminder_sends` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
    `calendar_id` bigint(20) unsigned NOT NULL,
    `entry_id` bigint(20) unsigned NOT NULL,
    `occurrence_key` varchar(128) NOT NULL,
    `step_at_minute` int(11) NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_fivenet_calendar_discord_reminder_sends_unique` (`calendar_id`, `entry_id`, `occurrence_key`, `step_at_minute`),
    KEY `idx_fivenet_calendar_discord_reminder_sends_created_at` (`created_at`),
    KEY `idx_fivenet_calendar_discord_reminder_sends_calendar_id` (`calendar_id`),
    KEY `idx_fivenet_calendar_discord_reminder_sends_entry_id` (`entry_id`),
    CONSTRAINT `fk_fivenet_calendar_discord_reminder_sends_calendar_id` FOREIGN KEY (`calendar_id`) REFERENCES `fivenet_calendar` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT `fk_fivenet_calendar_discord_reminder_sends_entry_id` FOREIGN KEY (`entry_id`) REFERENCES `fivenet_calendar_entries` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB;

COMMIT;
