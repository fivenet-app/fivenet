BEGIN;

-- Table: fivenet_job_props
CREATE TABLE IF NOT EXISTS `fivenet_job_props` (
  `job` varchar(20) NOT NULL,
  `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3),
  `theme` varchar(20) DEFAULT "defaultTheme",
  `livemap_marker_color` char(6) DEFAULT "5C7AFF",
  `quick_buttons` varchar(255) DEFAULT NULL,
  `radio_frequency` varchar(24) DEFAULT NULL,
  `discord_guild_id` bigint(20) unsigned DEFAULT NULL,
  `discord_last_sync` datetime(3) DEFAULT NULL,
  `discord_sync_settings` longtext DEFAULT NULL,
  `jobs_motd` text DEFAULT NULL,
  UNIQUE KEY `idx_fivenet_job_props_unique` (`job`),
  KEY `idx_fivenet_job_props_discord_guild_id` (`discord_guild_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

COMMIT;
