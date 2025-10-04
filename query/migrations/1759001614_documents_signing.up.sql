BEGIN;

-- Document Signing Stamps
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

  KEY `idx_fivenet_doc_signs_stamps_user_id` (`user_id`),
  KEY `idx_fivenet_doc_signs_stamps_sort_key` (`sort_key`),
  KEY `idx_fivenet_doc_signs_stamps_created_at` (`created_at`),
  KEY `idx_fivenet_doc_signs_stamps_deleted_at` (`deleted_at`),

  CONSTRAINT `fk_fivenet_documents_signatures_stamp_user` FOREIGN KEY (`user_id`) REFERENCES `{{.UsersTableName}}` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `fivenet_documents_signatures_stamps_access` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `target_id` bigint(20) unsigned NOT NULL,
  `job` varchar(40) NOT NULL,
  `minimum_grade` int(11) NOT NULL,
  `access` smallint(2) NOT NULL,

  PRIMARY KEY (`id`),

  UNIQUE KEY `uq_fivenet_doc_sig_stamps_access_unique` (`target_id`, `job`, `minimum_grade`),

  KEY `idx_fivenet_doc_sig_stamps_access_access` (`access`),

  CONSTRAINT `fk_fivenet_doc_sig_stamps_access_target_id` FOREIGN KEY (`target_id`) REFERENCES `fivenet_documents_signatures_stamps` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Document Signing System
CREATE TABLE IF NOT EXISTS `fivenet_documents_signature_policies` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `document_id` bigint(20) unsigned NOT NULL,
  `snapshot_date` datetime(3) NOT NULL,

  `label` varchar(120) DEFAULT NULL,
  `required` tinyint(1) NOT NULL DEFAULT 1,
  `binding_mode` smallint(2) NOT NULL,
  `allowed_types_mask` varchar(120) NOT NULL DEFAULT '[]',

  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` datetime(3) DEFAULT NULL,

  PRIMARY KEY (`id`),

  KEY `idx_fivenet_doc_sig_pol_doc_snapshot` (`document_id`, `snapshot_date`),
  KEY `idx_fivenet_doc_sig_pol_doc_snap_required` (`document_id`, `snapshot_date`, `required`),
  KEY `idx_fivenet_doc_sig_pol_deleted_at` (`deleted_at`),

  CONSTRAINT `fk_fivenet_documents_sig_pol_document_id` FOREIGN KEY (`document_id`) REFERENCES `fivenet_documents` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `fivenet_documents_signature_tasks` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `document_id` bigint(20) unsigned NOT NULL,
  `snapshot_date` datetime(3) NOT NULL,
  `policy_id` bigint(20) unsigned NOT NULL,

  -- Who is the task for? 1=USER, 2=JOB
  `assignee_kind` smallint(2) NOT NULL,

  -- User assignment
  `user_id` int(11) DEFAULT NULL,
  -- Job assignment
  `job` varchar(20) DEFAULT NULL,
  `minimum_grade` int DEFAULT NULL,
  `slot_no`       int NOT NULL DEFAULT 1,

  `status` smallint(2) NOT NULL,
  `comment` varchar(500) DEFAULT NULL,

  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `completed_at` datetime(3) DEFAULT NULL,
  `due_at` datetime(3) DEFAULT NULL,

  `signature_id` bigint(20) unsigned DEFAULT NULL,

  `creator_id` int(11) DEFAULT NULL,
  `creator_job` varchar(20) DEFAULT NULL,

  PRIMARY KEY (`id`),

  -- Prevent duplicates for the same user within a task
  UNIQUE KEY `uq_fivenet_doc_sig_tsk_user_round` (`policy_id`, `document_id`, `snapshot_date`, `assignee_kind`, `user_id`),
  -- And prevent duplicates for the same group target within a task
  UNIQUE KEY `uq_fivenet_doc_sig_tsk_group_round` (`policy_id`, `document_id`, `snapshot_date`, `assignee_kind`, `job`, `minimum_grade`, `slot_no`),

  KEY `idx_fivenet_doc_sigtsks_doc_snap_status` (`document_id`, `snapshot_date`, `status`),
  KEY `idx_fivenet_doc_sigtsks_policy_id_status` (`policy_id`, `status`),
  KEY `idx_fivenet_doc_sigtsks_user_status_create` (`user_id`, `status`, `created_at`),
  KEY `idx_fivenet_doc_sigtsks_access_check` (`document_id`, `snapshot_date`, `assignee_kind`, `job`, `minimum_grade`, `status`),

  CONSTRAINT `fk_fivenet_doc_signature_tasks_doc_id` FOREIGN KEY (`document_id`) REFERENCES `fivenet_documents` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_doc_signature_tasks_policy_id` FOREIGN KEY (`policy_id`) REFERENCES `fivenet_documents_signature_policies` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_doc_signature_tasks_user_id` FOREIGN KEY (`user_id`) REFERENCES `{{.UsersTableName}}` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `fivenet_documents_signatures` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `document_id` bigint(20) unsigned NOT NULL,
  `snapshot_date` datetime(3) NOT NULL,
  `policy_id` bigint(20) unsigned DEFAULT NULL,
  -- Link to the task that produced it (if any)
  `task_id` bigint(20) unsigned DEFAULT NULL,

  `user_id` int(11) NOT NULL,
  `user_job` varchar(20) NOT NULL,
  `user_job_grade` int DEFAULT NULL,

  `type` smallint(2) NOT NULL,
  `payload_svg` longtext NOT NULL,
  `stamp_id` bigint(20) unsigned DEFAULT NULL,

  -- 1=VALID, 2=REVOKED, 3=EXPIRED, 4=INVALID
  `status` smallint(2) NOT NULL,
  `comment` varchar(500) DEFAULT NULL,

  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `revoked_at` datetime(3) DEFAULT NULL,

  PRIMARY KEY (`id`),

  UNIQUE KEY `uq_fivenet_doc_signature_user_round` (`policy_id`, `snapshot_date`, `user_id`),

  KEY `idx_fivenet_doc_sig_doc_snapshot_status` (`document_id`, `snapshot_date`, `status`),
  KEY `idx_fivenet_doc_sig_user_created` (`user_id`, `created_at`),

  CONSTRAINT `fk_fivenet_doc_signatures_doc_id` FOREIGN KEY (`document_id`) REFERENCES `fivenet_documents` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_doc_signatures_policy_id` FOREIGN KEY (`policy_id`) REFERENCES `fivenet_documents_signature_policies` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_doc_signatures_task_id` FOREIGN KEY (`task_id`) REFERENCES `fivenet_documents_signature_tasks` (`id`) ON DELETE SET NULL ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_doc_signatures_user_id` FOREIGN KEY (`user_id`) REFERENCES `{{.UsersTableName}}` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_doc_signatures_stamp_id` FOREIGN KEY (`stamp_id`) REFERENCES `fivenet_documents_signatures_stamps` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB;

ALTER TABLE `fivenet_documents_signature_tasks`
  ADD CONSTRAINT `fk_fivenet_doc_sigtsks_signature_id` FOREIGN KEY (`signature_id`) REFERENCES `fivenet_documents_signatures` (`id`) ON DELETE SET NULL ON UPDATE CASCADE;

COMMIT;
