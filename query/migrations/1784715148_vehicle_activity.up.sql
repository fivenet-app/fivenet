BEGIN;

-- Table: fivenet_vehicles_activity
CREATE TABLE IF NOT EXISTS `fivenet_vehicles_activity` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `creator_id` int(11) DEFAULT NULL,
  `plate` varchar(12) NOT NULL,
  `type` smallint(2) NOT NULL,
  `creator_job` varchar(20) NOT NULL DEFAULT '',
  `reason` varchar(255) DEFAULT NULL,
  `data` longtext,
  PRIMARY KEY (`id`),
  KEY `idx_fivenet_vehicles_activity_creator_id` (`creator_id`),
  KEY `idx_fivenet_vehicles_activity_plate` (`plate`),
  KEY `idx_fivenet_vehicles_activity_created_at` (`created_at`),
  KEY `idx_fivenet_vehicles_activity_type` (`type`),
  CONSTRAINT `fk_fivenet_vehicles_activity_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `{{.UsersTableName}}` (`id`) ON DELETE SET NULL ON UPDATE SET NULL,
  CONSTRAINT `fk_fivenet_vehicles_activity_plate` FOREIGN KEY (`plate`) REFERENCES `{{- if .ESXCompat }}owned_vehicles{{ else }}fivenet_owned_vehicles{{ end -}}` (`plate`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

COMMIT;
