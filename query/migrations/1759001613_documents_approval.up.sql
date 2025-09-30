BEGIN;

-- Document Approval System
CREATE TABLE IF NOT EXISTS `fivenet_documents_approval_policies` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `document_id` bigint(20) unsigned NOT NULL,

  `on_edit_behavior` smallint(2) NOT NULL,
  `rule_kind` SMALLINT(2) DEFAULT 1,
  `required_count` int NOT NULL DEFAULT 1,
  `quorum_any` int DEFAULT NULL,
  `due_at` datetime(3) DEFAULT NULL,

  `active_snapshot_date` datetime(3) NOT NULL,

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

  KEY `idx_policy_doc_snapshot` (`document_id`, `active_snapshot_date`),
  KEY `idx_policy_started_at` (`started_at`),
  KEY `idx_policy_completed_at` (`completed_at`),
  KEY `idx_policy_deleted_at` (`deleted_at`),

  CONSTRAINT `fk_policy_doc` FOREIGN KEY (`document_id`) REFERENCES `fivenet_documents` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `fivenet_documents_approval_access` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `target_id` bigint unsigned NOT NULL,
  `user_id` int DEFAULT NULL,
  `job` varchar(40) DEFAULT NULL,
  `minimum_grade` int DEFAULT NULL,
  `access` smallint NOT NULL,

  PRIMARY KEY (`id`),

  UNIQUE KEY `idx_user_id_access_unique` (`target_id`, `user_id`),
  UNIQUE KEY `idx_job_minimum_grade_access_unique` (`target_id`, `job`, `minimum_grade`),

  KEY `fk_documents_approval_stage_access_user_id` (`user_id`),
  KEY `fk_documents_approval_stage_access_access` (`access`),
  KEY `idx_job_minimum_grade` (`job`, `minimum_grade`),
  KEY `idx_access_target_access_user` (`target_id`, `access`, `user_id`),
  KEY `idx_access_target_access_job_grade` (`target_id`, `access`, `job`, `minimum_grade`),

  CONSTRAINT `fk_fivenet_doc_approval_stages_access_target_id` FOREIGN KEY (`target_id`) REFERENCES `fivenet_documents_approval_policies` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_doc_approval_stages_access_user_id` FOREIGN KEY (`user_id`) REFERENCES `{{.UsersTableName}}` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `fivenet_documents_approval_tasks` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `document_id` bigint(20) unsigned NOT NULL,
  `snapshot_date` datetime(3) NOT NULL,
  -- Who is the task for? 1=USER, 2=FACTION_GRADE
  `assignee_kind` smallint(2) NOT NULL,
  -- USER assignment
  `user_id` int(11) DEFAULT NULL,
  -- set when assignee_kind=USER
  -- FACTION_GRADE assignment (claimable by any eligible user)
  `job` varchar(20) DEFAULT NULL,
  `minimum_grade` int DEFAULT NULL,
  -- Snapshot of reviewer context WHEN DECIDED (stable audit for group tasks)
  `decided_by_user_id` int(11) DEFAULT NULL,
  `decided_by_job` varchar(20) DEFAULT NULL,
  `decided_by_user_grade` int DEFAULT NULL,
  `status` smallint(2) NOT NULL,
  `comment` varchar(500) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `decided_at` datetime(3) DEFAULT NULL,
  `due_at` datetime(3) DEFAULT NULL,
  `decision_count` int NOT NULL DEFAULT 0,

  PRIMARY KEY (`id`),

  -- Prevent duplicates for the same user within a stage
  UNIQUE KEY `uq_stage_user` (`document_id`, `snapshot_date`, `assignee_kind`, `user_id`),
  -- And prevent duplicates for the same group target within a stage
  UNIQUE KEY `uq_stage_group` (`document_id`, `snapshot_date`, `assignee_kind`, `job`, `minimum_grade`),

  KEY `idx_task_user_status_created` (`user_id`, `status`, `created_at`),
  KEY `idx_task_doc_snapshot_status` (`document_id`, `snapshot_date`, `status`),
  KEY `idx_task_status` (`status`),

  CONSTRAINT `fk_task_doc` FOREIGN KEY (`document_id`) REFERENCES `fivenet_documents` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

COMMIT;
