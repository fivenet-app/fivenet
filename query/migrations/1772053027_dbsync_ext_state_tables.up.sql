BEGIN;

CREATE TABLE `fivenet_sync_user` (
  `user_id` int NOT NULL,
  `identifier` varchar(64) NOT NULL,
  -- Metadata about the data source
  `source_updated_at` datetime(3) DEFAULT NULL,
  `last_synced_at` datetime(3) NOT NULL,
  `data_json` json NOT NULL,
  `data_hash` bigint unsigned DEFAULT NULL,

  PRIMARY KEY (`user_id`),

  UNIQUE KEY `ux_user_external_identifier` (`identifier`),

  CONSTRAINT `fk_user_external_user_id` FOREIGN KEY (`user_id`) REFERENCES `fivenet_user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

COMMIT;
