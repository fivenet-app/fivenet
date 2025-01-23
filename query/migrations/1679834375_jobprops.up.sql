BEGIN;

-- Table: fivenet_job_props
CREATE TABLE IF NOT EXISTS `fivenet_job_props` (
  `job` varchar(20) NOT NULL,
  `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3),
  `theme` varchar(20) DEFAULT "defaultTheme",
  `livemap_marker_color` char(7) DEFAULT "#5c7aff",
  `quick_buttons` varchar(255) DEFAULT NULL,
  `radio_frequency` varchar(24) DEFAULT NULL,
  `discord_guild_id` varchar(128) DEFAULT NULL,
  `discord_last_sync` datetime(3) DEFAULT NULL,
  `discord_sync_settings` longtext DEFAULT NULL,
  `discord_sync_changes` longtext DEFAULT NULL,
  `motd` text DEFAULT NULL,
  `logo_url` varchar(128) DEFAULT NULL,
  `settings` longtext DEFAULT NULL,
  `citizen_attributes` text DEFAULT NULL,
  UNIQUE KEY `idx_fivenet_job_props_unique` (`job`),
  KEY `idx_fivenet_job_props_discord_guild_id` (`discord_guild_id`),
  KEY `idx_fivenet_job_props_logo_url` (`logo_url`)
) ENGINE=InnoDB;

-- Table: fivenet_job_citizen_attributes
CREATE TABLE IF NOT EXISTS `fivenet_job_citizen_attributes` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `job` varchar(20) NOT NULL,
  `name` varchar(32) NOT NULL,
  `color` char(7) DEFAULT "#5c7aff",
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_job_citizen_attributes_unique` (`job`, `name`),
  KEY `idx_fivenet_job_citizen_attributes_name` (`name`)
) ENGINE=InnoDB;

-- Table: fivenet_user_citizen_attributes
CREATE TABLE IF NOT EXISTS `fivenet_user_citizen_attributes` (
  `user_id` int(11) NOT NULL,
  `attribute_id` bigint(20) unsigned NOT NULL,
  UNIQUE KEY `idx_fivenet_user_citizen_attributes_unique` (`user_id`, `attribute_id`),
  CONSTRAINT `fk_fivenet_user_citizen_attributes_user_id` FOREIGN KEY (`user_id`) REFERENCES `{{.UsersTableName}}` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_user_citizen_attributes_attribute_id` FOREIGN KEY (`attribute_id`) REFERENCES `fivenet_job_citizen_attributes` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

COMMIT;
