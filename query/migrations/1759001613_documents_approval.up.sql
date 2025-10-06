BEGIN;

-- Document Approval System
CREATE TABLE IF NOT EXISTS `fivenet_documents_approval_policies` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `document_id` bigint(20) unsigned NOT NULL,
  `snapshot_date` datetime(3) NOT NULL,

  `on_edit_behavior` smallint(2) NOT NULL,
  `rule_kind` SMALLINT(2) DEFAULT 1,
  `required_count` int NOT NULL DEFAULT 1,
  `due_at` datetime(3) DEFAULT NULL,

  `assigned_count` int NOT NULL DEFAULT 0,
  `approved_count` int NOT NULL DEFAULT 0,
  `declined_count` int NOT NULL DEFAULT 0,
  `pending_count` int NOT NULL DEFAULT 0,
  `any_declined` TINYINT(1) NOT NULL DEFAULT 0,

  `started_at` DATETIME(3),
  `completed_at` DATETIME(3),

  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` datetime(3) DEFAULT NULL,

  PRIMARY KEY (`id`),

  -- For now "limit" to one policy per document
  UNIQUE KEY `uq_fivenet_doc_approval_pol_doc` (`document_id`),

  KEY `idx_fivenet_doc_approval_pol_doc_snapshot` (`document_id`, `snapshot_date`),
  KEY `idx_fivenet_doc_approval_pol_started_at` (`started_at`),
  KEY `idx_fivenet_doc_approval_pol_completed_at` (`completed_at`),
  KEY `idx_fivenet_doc_approval_pol_deleted_at` (`deleted_at`),

  CONSTRAINT `fk_fivenet_documents_approval_policies_doc` FOREIGN KEY (`document_id`) REFERENCES `fivenet_documents` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `fivenet_documents_approval_tasks` (
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
  `slot_no` INT NOT NULL DEFAULT 1,

  `status` smallint(2) NOT NULL,
  `comment` varchar(500) DEFAULT NULL,
  `due_at` datetime(3) DEFAULT NULL,
  `decision_count` int NOT NULL DEFAULT 0,

  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `decided_at` datetime(3) DEFAULT NULL,

  `approval_id` bigint(20) unsigned DEFAULT NULL,

  `creator_id` int(11) DEFAULT NULL,
  `creator_job` varchar(20) DEFAULT NULL,

  PRIMARY KEY (`id`),

  -- Prevent duplicates for the same user within a task
  UNIQUE KEY `uq_fivenet_doc_approval_task_pol_user` (`policy_id`, `document_id`, `snapshot_date`, `assignee_kind`, `user_id`),
  -- And prevent duplicates for the same group target within a task
  UNIQUE KEY `uq_fivenet_doc_approval_task_job_pol_slot` (`policy_id`, `document_id`, `snapshot_date`, `assignee_kind`, `job`, `minimum_grade`, `slot_no`),

  KEY `idx_fivenet_doc_apptsk_doc_snap_status` (`document_id`, `snapshot_date`, `status`),
  KEY `idx_fivenet_doc_apptsk_policy_status` (`policy_id`, `status`),
  KEY `idx_fivenet_doc_apptsk_user_status_created` (`user_id`, `status`, `created_at`),
  KEY `idx_fivenet_doc_apptsk_access_check` (`document_id`, `snapshot_date`, `assignee_kind`, `job`, `minimum_grade`, `status`),

  CONSTRAINT `fk_fivenet_doc_apptsk_task_doc_id` FOREIGN KEY (`document_id`) REFERENCES `fivenet_documents` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_doc_apptsk_task_policy_id` FOREIGN KEY (`policy_id`) REFERENCES `fivenet_documents_approval_policies` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_doc_apptsk_task_user_id` FOREIGN KEY (`user_id`) REFERENCES `{{.UsersTableName}}` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `fivenet_documents_approvals` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `document_id` bigint(20) unsigned NOT NULL,
  `snapshot_date` datetime(3) NOT NULL,
  `policy_id` bigint(20) unsigned DEFAULT NULL,
  -- Link to the task that produced it (if any)
  `task_id` bigint(20) unsigned DEFAULT NULL,

  `user_id` int(11) NOT NULL,
  `user_job` varchar(20) NOT NULL,
  `user_job_grade` int DEFAULT NULL,

  -- 1=APPROVED, 2=DECLINED, 3=REVOKED (optional)
  `status` smallint(2) NOT NULL,
  `comment` varchar(500) DEFAULT NULL,

  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `revoked_at` datetime(3) DEFAULT NULL,

  PRIMARY KEY (`id`),

  UNIQUE KEY `uq_fivenet_doc_approval_user_round` (`policy_id`, `snapshot_date`, `user_id`),

  KEY `idx_fivenet_doc_approval_user_created` (`user_id`, `created_at`),
  KEY `idx_fivenet_doc_approval_user_doc` (`user_id`, `document_id`),
  KEY `idx_fivenet_doc_approval_policy_status` (`policy_id`,`status`),
  KEY `idx_fivenet_doc_approval_doc_snap_status` (`document_id`, `snapshot_date`, `status`),

  CONSTRAINT `fk_fivenet_doc_approvals_doc_id` FOREIGN KEY (`document_id`) REFERENCES `fivenet_documents` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_doc_approvals_policy_id` FOREIGN KEY (`policy_id`) REFERENCES `fivenet_documents_approval_policies`(`id`) ON DELETE SET NULL ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_doc_approvals_task_id` FOREIGN KEY (`task_id`) REFERENCES `fivenet_documents_approval_tasks` (`id`) ON DELETE SET NULL ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_doc_approvals_user_id` FOREIGN KEY (`user_id`) REFERENCES `{{.UsersTableName}}` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

ALTER TABLE `fivenet_documents_approval_tasks`
  ADD CONSTRAINT `fk_fivenet_doc_apptsk_task_approval_id` FOREIGN KEY (`approval_id`) REFERENCES `fivenet_documents_approvals` (`id`) ON DELETE SET NULL ON UPDATE CASCADE;

COMMIT;
