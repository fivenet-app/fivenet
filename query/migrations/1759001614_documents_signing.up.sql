BEGIN;

-- Document Signing System
CREATE TABLE IF NOT EXISTS `fivenet_documents_signature_policies` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `document_id` bigint(20) unsigned NOT NULL,
  `snapshot_date` datetime(3) NOT NULL,

  `label` varchar(120) DEFAULT NULL,
  `required` tinyint(1) NOT NULL DEFAULT 1,

  `binding_mode` smallint(2) NOT NULL,
  `allowed_types_mask` varchar(120) NOT NULL DEFAULT '[]',

  `collected_count` int NOT NULL DEFAULT 0,
  `required_count` int NOT NULL DEFAULT 1,

  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` datetime(3) DEFAULT NULL,

  PRIMARY KEY (`id`),

  KEY `idx_sigreq_doc_snapshot` (`document_id`, `snapshot_date`),
  KEY `idx_sigreq_doc_snap_required` (`document_id`, `snapshot_date`, `required`),
  KEY `idx_sigreq_deleted_at` (`deleted_at`),

  CONSTRAINT `fk_sigreq_doc` FOREIGN KEY (`document_id`) REFERENCES `fivenet_documents` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `fivenet_documents_signature_policies_access` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `target_id` bigint unsigned NOT NULL,
  `user_id` int DEFAULT NULL,
  `job` varchar(40) DEFAULT NULL,
  `minimum_grade` int DEFAULT NULL,
  `access` smallint NOT NULL,

  PRIMARY KEY (`id`),

  UNIQUE KEY `idx_user_id_access_unique` (`target_id`, `user_id`),
  UNIQUE KEY `idx_job_minimum_grade_access_unique` (`target_id`,`job`,`minimum_grade`),

  KEY `fk_documents_signature_reqs_access_user_id` (`user_id`),
  KEY `fk_documents_signature_reqs_access_access` (`access`),
  KEY `idx_job_minimum_grade` (`job`,`minimum_grade`),
  KEY `idx_access_target_access_user` (`target_id`,`access`),
  KEY `idx_access_target_access_job_grade` (`target_id`,`access`,`job`,`minimum_grade`),

  CONSTRAINT `fk_fivenet_doc_signature_reqs_access_target_id` FOREIGN KEY (`target_id`) REFERENCES `fivenet_documents_signature_policies` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_doc_signature_reqs_access_user_id` FOREIGN KEY (`user_id`) REFERENCES `{{.UsersTableName}}` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `fivenet_documents_signatures_stamps` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(120) NOT NULL,

  `user_id` int(11) DEFAULT NULL,

  `svg_template` mediumtext NOT NULL,
  -- Parameterized SVG with slots (if any)
  `variants_json` longtext DEFAULT NULL,

  `sort_key` varchar(255) GENERATED ALWAYS AS ((
    CASE
      WHEN (REGEXP_SUBSTR(`name`, '[0-9]+') IS NOT NULL) THEN
        REGEXP_REPLACE(`name`, '[0-9]+', LPAD(REGEXP_SUBSTR(`name`, '[0-9]+'), 8, '0'))
      ELSE `name`
    END
  )) STORED,

  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `deleted_at` datetime(3) DEFAULT NULL,

  PRIMARY KEY (`id`),

  KEY `idx_stamp_user` (`user_id`),
  KEY `idx_stamp_sort_key` (`sort_key`),
  KEY `idx_stamp_created` (`created_at`),
  KEY `idx_stamp_deleted_at` (`deleted_at`),

  CONSTRAINT `fk_stamp_user` FOREIGN KEY (`user_id`) REFERENCES `{{.UsersTableName}}` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `fivenet_documents_signatures_stamps_access` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `target_id` bigint(20) unsigned NOT NULL,
  `job` varchar(40) NOT NULL,
  `minimum_grade` int(11) NOT NULL,
  `access` smallint(2) NOT NULL,

  PRIMARY KEY (`id`),

  UNIQUE KEY `idx_documents_stamps_access_unique_access` (`target_id`, `job`, `minimum_grade`),

  KEY `fk_documents_stamps_access_access` (`access`),

  CONSTRAINT `fk_documents_signatures_stamps_access_target_id` FOREIGN KEY (`target_id`) REFERENCES `fivenet_documents_signatures_stamps` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `fivenet_documents_signatures` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `document_id` bigint(20) unsigned NOT NULL,
  `snapshot_date` datetime(3) NOT NULL,
  `policy_id` bigint(20) unsigned DEFAULT NULL,
  `user_id` int(11) NOT NULL,
  `user_job` varchar(20) NOT NULL,

  `type` smallint(2) NOT NULL,
  `payload_json` longtext NOT NULL,
  `stamp_id` bigint(20) unsigned DEFAULT NULL,

  `status` smallint(2) NOT NULL,
  `reason` varchar(255) DEFAULT NULL,

  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `revoked_at` datetime(3) DEFAULT NULL,

  PRIMARY KEY (`id`),

  UNIQUE KEY `uq_req_user` (`policy_id`, `user_id`),
  KEY `idx_sig_doc_snapshot_status` (`document_id`, `snapshot_date`, `status`),
  KEY `idx_sig_user_created` (`user_id`, `created_at`),
  KEY `idx_sig_requirement_status` (`policy_id`, `status`),

  CONSTRAINT `fk_sig_doc` FOREIGN KEY (`document_id`) REFERENCES `fivenet_documents` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_sig_req` FOREIGN KEY (`policy_id`) REFERENCES `fivenet_documents_signature_policies` (`id`) ON DELETE SET NULL ON UPDATE CASCADE,
  CONSTRAINT `fk_sig_stamp` FOREIGN KEY (`stamp_id`) REFERENCES `fivenet_documents_signatures_stamps` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB;

CREATE TABLE `fivenet_documents_signature_tasks` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `document_id` bigint(20) unsigned NOT NULL,
  `snapshot_date` datetime(3) NOT NULL,

  `policy_id` bigint(20) unsigned NOT NULL,

  `assignee_kind` smallint(2) NOT NULL, -- 1=USER, 2=FACTION_GRADE
  `user_id` int(11) DEFAULT NULL,
  `job` varchar(20) DEFAULT NULL,
  `minimum_grade` int DEFAULT NULL,

  `status` smallint(2) NOT NULL,
  `comment` varchar(500) DEFAULT NULL,
  `due_at` datetime(3) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `completed_at` datetime(3) DEFAULT NULL,

  `signature_id` bigint(20) unsigned DEFAULT NULL,

  PRIMARY KEY (`id`),

  UNIQUE KEY `uq_sigtsk_user_round` (`policy_id`,`document_id`,`snapshot_date`,`assignee_kind`,`user_id`),
  UNIQUE KEY `uq_sigtsk_group_round` (`policy_id`,`document_id`,`snapshot_date`,`assignee_kind`,`job`,`minimum_grade`),

  KEY `idx_sigtsk_doc_snap_status` (`document_id`,`snapshot_date`,`status`),
  KEY `idx_sigtsk_policy_status` (`policy_id`,`status`),
  KEY `idx_sigtsk_user_status_created` (`user_id`,`status`,`created_at`),

  CONSTRAINT `fk_sigtsk_doc` FOREIGN KEY (`document_id`) REFERENCES `fivenet_documents` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_sigtsk_policy` FOREIGN KEY (`policy_id`) REFERENCES `fivenet_documents_signature_policies` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_sigtsk_signature` FOREIGN KEY (`signature_id`) REFERENCES `fivenet_documents_signatures` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB;

COMMIT;
