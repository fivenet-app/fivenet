BEGIN;

-- Table: fivenet_documents_pins
CREATE TABLE IF NOT EXISTS `fivenet_documents_pins` (
  `document_id` bigint(20) unsigned NOT NULL,
  `job` varchar(50) NOT NULL,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `state` tinyint(1) DEFAULT 1,
  `creator_id` int(11) NOT NULL,
  PRIMARY KEY (`document_id`, `job`),
  KEY `idx_fivenet_documents_pins_document_id` (`document_id`),
  KEY `idx_fivenet_documents_pins_creator_id` (`creator_id`),
  CONSTRAINT `fk_fivenet_documents_pins_document_id` FOREIGN KEY (`document_id`) REFERENCES `fivenet_documents` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_documents_pins_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `{{.UsersTableName}}` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

COMMIT;
