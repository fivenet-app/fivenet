BEGIN;

-- Document Approval System

CREATE TABLE `fivenet_documents_approval_policies` (
  `document_id` bigint(20) unsigned NOT NULL,
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,

  `on_edit_behavior` smallint(2) NOT NULL,

  PRIMARY KEY (`document_id`),

  CONSTRAINT `fk_policy_doc` FOREIGN KEY (`document_id`) REFERENCES `fivenet_documents` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB;

CREATE TABLE `fivenet_documents_approval_stages` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `document_id` bigint(20) unsigned NOT NULL,
  `snapshot_date` datetime(3) NOT NULL,

  `name` varchar(120) DEFAULT NULL,
  `order` int NOT NULL,
  `selector_json` longtext NOT NULL,

  `require_all` tinyint(1) NOT NULL DEFAULT 0,
  `quorum_any` int DEFAULT NULL,
  `due_at` datetime(3) DEFAULT NULL,

  `assigned_count` int NOT NULL DEFAULT 0,
  `approved_count` int NOT NULL DEFAULT 0,
  `declined_count` int NOT NULL DEFAULT 0,
  `pending_count` int NOT NULL DEFAULT 0,

  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,

  PRIMARY KEY (`id`),

  KEY `idx_stage_doc_snapshot` (`document_id`, `snapshot_date`, `order`),
  KEY `idx_stage_policy` (`policy_id`),

  CONSTRAINT `fk_stage_doc` FOREIGN KEY (`document_id`) REFERENCES `fivenet_documents` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_stage_policy_doc` FOREIGN KEY (`document_id`) REFERENCES `fivenet_documents_approval_policies` (`document_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB;

CREATE TABLE `fivenet_documents_approval_tasks` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `stage_id` bigint(20) unsigned NOT NULL,
  `document_id` bigint(20) unsigned NOT NULL,
  `snapshot_date` datetime(3) NOT NULL,
  -- Who is the task for? 1=USER, 2=FACTION_GRADE
  `assignee_kind` smallint(2) NOT NULL,
  -- USER assignment
  `user_id` int(11) DEFAULT NULL,
  -- set when assignee_kind=USER
  -- FACTION_GRADE assignment (claimable by any eligible user)
  `job` varchar(20) DEFAULT NULL,
  `min_grade` int DEFAULT NULL,
  -- Snapshot of reviewer context WHEN DECIDED (stable audit for group tasks)
  `decided_by_user_id` int(11) DEFAULT NULL,
  `decided_by_job` varchar(20) DEFAULT NULL,
  `decided_by_user_grade` int DEFAULT NULL,
  `status` smallint(2) NOT NULL,
  `comment` varchar(500) DEFAULT NULL,
  `created_at` datetime(3) NOT NULL,
  `decided_at` datetime(3) DEFAULT NULL,
  `due_at` datetime(3) DEFAULT NULL,

  PRIMARY KEY (`id`),

  -- Prevent duplicates for the same user within a stage
  UNIQUE KEY `uq_stage_user` (`stage_id`, `user_id`),
  -- And prevent duplicates for the same group target within a stage
  UNIQUE KEY `uq_stage_group` (`stage_id`, `assignee_kind`, `job`, `min_grade`),

  KEY `idx_task_user_status_created` (`user_id`, `status`, `created_at`),
  KEY `idx_task_doc_snapshot_status` (`document_id`, `snapshot_date`, `status`),
  KEY `idx_task_stage_status` (`stage_id`, `status`),

  CONSTRAINT `fk_task_stage` FOREIGN KEY (`stage_id`) REFERENCES `fivenet_documents_approval_stages`(`id`),
  CONSTRAINT `fk_task_doc` FOREIGN KEY (`document_id`) REFERENCES `fivenet_documents` (`id`)
) ENGINE = InnoDB;

-- Document Signing System

CREATE TABLE `fivenet_documents_signature_requirements` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `document_id` bigint(20) unsigned NOT NULL,
  `snapshot_date` datetime(3) NOT NULL,

  `sequence_order` int DEFAULT NULL,
  `required` tinyint(1) NOT NULL DEFAULT 1,
  `label` varchar(120) DEFAULT NULL,

  `selector_json` longtext NOT NULL,
  `binding_mode` smallint(2) NOT NULL,
  `allowed_types_mask` smallint(2) NOT NULL DEFAULT 7,

  `collected_count` int NOT NULL DEFAULT 0,
  `required_count` int NOT NULL DEFAULT 1,

  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,

  PRIMARY KEY (`id`),

  KEY `idx_sigreq_doc_snapshot` (`document_id`, `snapshot_date`),
  KEY `idx_sigreq_doc_snap_required` (`document_id`, `snapshot_date`, `required`),

  CONSTRAINT `fk_sigreq_doc` FOREIGN KEY (`document_id`) REFERENCES `fivenet_documents` (`id`)
) ENGINE = InnoDB;

CREATE TABLE `fivenet_documents_signatures` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `document_id` bigint(20) unsigned NOT NULL,
  `snapshot_date` datetime(3) NOT NULL,
  `requirement_id` bigint(20) unsigned DEFAULT NULL,
  `user_id` int(11) NOT NULL,
  `user_job` varchar(20) NOT NULL,

  `type` smallint(2) NOT NULL,
  `payload_json` longtext NOT NULL,
  `stamp_id` bigint(20) unsigned DEFAULT NULL,

  `status` smallint(2) NOT NULL,
  `reason` varchar(255) DEFAULT NULL,

  `created_at` datetime(3) NOT NULL,
  `revoked_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),

  UNIQUE KEY `uq_req_user` (`requirement_id`, `user_id`),
  KEY `idx_sig_doc_snapshot_status` (`document_id`, `snapshot_date`, `status`),
  KEY `idx_sig_user_created` (`user_id`, `created_at`),
  KEY `idx_sig_requirement_status` (`requirement_id`, `status`),

  CONSTRAINT `fk_sig_doc` FOREIGN KEY (`document_id`) REFERENCES `fivenet_documents` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_sig_req` FOREIGN KEY (`requirement_id`) REFERENCES `fivenet_documents_signature_requirements` (`id`) ON DELETE SET NULL ON UPDATE CASCADE,
  CONSTRAINT `fk_sig_stamp` FOREIGN KEY (`stamp_id`) REFERENCES `fivenet_documents_signatures_stamps` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE = InnoDB;

CREATE TABLE `fivenet_documents_signatures_stamps` (
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

  `created_at` datetime(3) NOT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,

  PRIMARY KEY (`id`),

  KEY `idx_stamp_user` (`user_id`),
  KEY `idx_stamp_sort_key` (`sort_key`);
  KEY `idx_stamp_created` (`created_at`),
  KEY `idx_stamp_deleted_at` (`deleted_at`)
) ENGINE = InnoDB;

CREATE TABLE `fivenet_documents_signatures_stamps_access` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `target_id` bigint(20) unsigned NOT NULL,
  `job` varchar(40) NOT NULL,
  `minimum_grade` int(11) NOT NULL,
  `access` smallint(2) NOT NULL,

  PRIMARY KEY (`id`),

  UNIQUE KEY `idx_documents_stamps_pages_access_unique_access` (`target_id`, `job`, `minimum_grade`),

  KEY `fk_documents_stamps_pages_access_access` (`access`),

  CONSTRAINT `fk_documents_stamps_pages_access_target_id` FOREIGN KEY (`target_id`) REFERENCES `fivenet_documents_signatures_stamps` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB;

-- Document Meta/Approval/Signing Table
CREATE TABLE `fivenet_documents_meta` (
  `document_id` bigint(20) unsigned NOT NULL,

  `recomputed_at` datetime(3) NOT NULL,

  `approved` tinyint(1) NOT NULL DEFAULT 0,
  `signed` tinyint(1) NOT NULL DEFAULT 0,

  `approval_current_order` int DEFAULT NULL,
  `approval_pending_tasks` int NOT NULL DEFAULT 0,
  `approval_any_declined` tinyint(1) NOT NULL DEFAULT 0,

  `sig_required_remaining` int NOT NULL DEFAULT 0,
  `sig_required_total` int NOT NULL DEFAULT 0,
  `sig_collected_valid` int NOT NULL DEFAULT 0,

  PRIMARY KEY (`document_id`),

  KEY `idx_documents_meta_approved` (`approved`),
  KEY `idx_documents_meta_signed` (`signed`),

  CONSTRAINT `fk_sum_doc` FOREIGN KEY (`document_id`) REFERENCES `fivenet_documents` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB;

COMMIT;
