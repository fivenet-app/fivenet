BEGIN;

CREATE TABLE IF NOT EXISTS `fivenet_vehicle_files` (
  `plate` varchar(12) NOT NULL,
  `file_id` bigint unsigned NOT NULL,
  `kind` varchar(32) NOT NULL DEFAULT 'attachment',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`plate`, `file_id`),
  KEY `idx_vehicle_files_file_id` (`file_id`),
  KEY `idx_vehicle_files_plate_kind` (`plate`, `kind`),
  CONSTRAINT `fk_fivenet_vehicle_files_plate` FOREIGN KEY (`plate`) REFERENCES `fivenet_vehicle` (`plate`)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_vehicle_files_file_id` FOREIGN KEY (`file_id`) REFERENCES `fivenet_files` (`id`)
    ON DELETE RESTRICT
) ENGINE=InnoDB;

COMMIT;
