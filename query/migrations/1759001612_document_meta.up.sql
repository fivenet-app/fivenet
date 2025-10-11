BEGIN;

-- Document Meta Table
CREATE TABLE IF NOT EXISTS `fivenet_documents_meta` (
  `document_id` bigint unsigned NOT NULL,

  `recomputed_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),

  -- Overall booleans (fast list flags)
  `approved` tinyint(1) NOT NULL DEFAULT 0,
  `signed` tinyint(1) NOT NULL DEFAULT 0,

  -- Signatures rollup (aggregated across required signature policies)
  `sig_required_remaining` int NOT NULL DEFAULT 0,
  `sig_required_total` int NOT NULL DEFAULT 0,
  `sig_collected_valid` int NOT NULL DEFAULT 0,
  `sig_policies_active` int NOT NULL DEFAULT 0,

  -- Approvals rollup (aggregated across required approval policies)
  `ap_required_total` int NOT NULL DEFAULT 0,
  `ap_collected_approved` int NOT NULL DEFAULT 0,
  `ap_required_remaining` int NOT NULL DEFAULT 0,
  `ap_declined_count` int NOT NULL DEFAULT 0,
  `ap_pending_count` int NOT NULL DEFAULT 0,
  `ap_any_declined` tinyint(1) NOT NULL DEFAULT 0,
  `ap_policies_active` int NOT NULL DEFAULT 0,

  PRIMARY KEY (`document_id`),

  KEY `idx_fivenet_documents_meta_approved` (`approved`),
  KEY `idx_fivenet_documents_meta_signed` (`signed`),
  KEY `idx_documents_meta_ap_required_remaining` (`ap_required_remaining`),
  KEY `idx_documents_meta_sig_required_remaining` (`sig_required_remaining`),

  CONSTRAINT `fk_fivenet_documents_meta_document_id` FOREIGN KEY (`document_id`) REFERENCES `fivenet_documents` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

COMMIT;
