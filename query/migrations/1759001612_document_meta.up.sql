BEGIN;

-- Document Meta Table
CREATE TABLE IF NOT EXISTS `fivenet_documents_meta` (
  `document_id` bigint(20) unsigned NOT NULL,

  `recomputed_at` datetime(3) NOT NULL,

  `approved` tinyint(1) NOT NULL DEFAULT 0,
  `signed` tinyint(1) NOT NULL DEFAULT 0,

  `sig_required_remaining` int NOT NULL DEFAULT 0,
  `sig_required_total` int NOT NULL DEFAULT 0,
  `sig_collected_valid` int NOT NULL DEFAULT 0,

  PRIMARY KEY (`document_id`),

  KEY `idx_documents_meta_approved` (`approved`),
  KEY `idx_documents_meta_signed` (`signed`),

  CONSTRAINT `fk_meta_doc` FOREIGN KEY (`document_id`) REFERENCES `fivenet_documents` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

COMMIT;
