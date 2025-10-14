BEGIN;

-- Approval system changes
ALTER TABLE `fivenet_documents_approvals` DROP KEY `uq_fivenet_doc_approval_user_round`;
ALTER TABLE `fivenet_documents_approvals` DROP FOREIGN KEY `fk_fivenet_doc_approvals_policy_id`;
ALTER TABLE `fivenet_documents_approvals` DROP INDEX `idx_fivenet_doc_approval_policy_status`;
ALTER TABLE `fivenet_documents_approvals` DROP COLUMN `policy_id`;
ALTER TABLE `fivenet_documents_approvals` ADD UNIQUE KEY `uq_fivenet_doc_approval_user_round` (`document_id`, `snapshot_date`, `user_id`);

ALTER TABLE `fivenet_documents_approval_tasks` DROP KEY `uq_fivenet_doc_approval_task_pol_user`;
ALTER TABLE `fivenet_documents_approval_tasks` DROP FOREIGN KEY `fk_fivenet_doc_apptsk_task_policy_id`;
ALTER TABLE `fivenet_documents_approval_tasks` DROP INDEX `idx_fivenet_doc_apptsk_policy_status`;
ALTER TABLE `fivenet_documents_approval_tasks` DROP KEY `uq_fivenet_doc_approval_task_job_pol_slot`;
ALTER TABLE `fivenet_documents_approval_tasks` DROP COLUMN `policy_id`;

ALTER TABLE `fivenet_documents_approval_tasks` ADD UNIQUE KEY `uq_fivenet_doc_approval_task_user` (`document_id`, `snapshot_date`, `assignee_kind`, `user_id`);
ALTER TABLE `fivenet_documents_approval_tasks` ADD UNIQUE KEY `uq_fivenet_doc_approval_task_job_pol_slot` (`document_id`, `snapshot_date`, `assignee_kind`, `job`, `minimum_grade`, `slot_no`);
ALTER TABLE `fivenet_documents_approval_tasks` ADD KEY `idx_fivenet_doc_apptsk_doc_status` (`document_id`, `status`);

ALTER TABLE `fivenet_documents_approval_policies` DROP INDEX uq_fivenet_doc_approval_pol_doc;
ALTER TABLE `fivenet_documents_approval_policies` DROP COLUMN id;
ALTER TABLE `fivenet_documents_approval_policies` ADD PRIMARY KEY (`document_id`);

ALTER TABLE `fivenet_documents_approval_tasks` ADD COLUMN `label` varchar(120) DEFAULT NULL AFTER `minimum_grade`;

-- Signature system changes
ALTER TABLE `fivenet_documents_signature_policies` DROP COLUMN `label`;
ALTER TABLE `fivenet_documents_signature_policies` DROP INDEX `idx_fivenet_doc_sig_pol_doc_snap_required`;
ALTER TABLE `fivenet_documents_signature_policies` DROP COLUMN `required`;

ALTER TABLE `fivenet_documents_signatures` DROP KEY `uq_fivenet_doc_signature_user_round`;
ALTER TABLE `fivenet_documents_signatures` DROP FOREIGN KEY `fk_fivenet_doc_signatures_policy_id`;
ALTER TABLE `fivenet_documents_signatures` DROP COLUMN `policy_id`;
ALTER TABLE `fivenet_documents_signatures` ADD UNIQUE KEY `uq_fivenet_doc_signature_user_round` (`snapshot_date`, `user_id`);

ALTER TABLE `fivenet_documents_signature_tasks` DROP INDEX `idx_fivenet_doc_sigtsks_policy_id_status`;
ALTER TABLE `fivenet_documents_signature_tasks` DROP FOREIGN KEY `fk_fivenet_doc_signature_tasks_policy_id`;
ALTER TABLE `fivenet_documents_signature_tasks` DROP KEY `uq_fivenet_doc_sig_tsk_user_round`;
ALTER TABLE `fivenet_documents_signature_tasks` DROP KEY `uq_fivenet_doc_sig_tsk_group_round`;
ALTER TABLE `fivenet_documents_signature_tasks` DROP COLUMN `policy_id`;
ALTER TABLE `fivenet_documents_signature_tasks` ADD KEY `idx_fivenet_doc_sigtsks_status` (`status`);
ALTER TABLE `fivenet_documents_signature_tasks` ADD UNIQUE KEY `uq_fivenet_doc_sig_tsk_user_round` (`document_id`, `snapshot_date`, `assignee_kind`, `user_id`);
ALTER TABLE `fivenet_documents_signature_tasks` ADD UNIQUE KEY `uq_fivenet_doc_sig_tsk_group_round` (`document_id`, `snapshot_date`, `assignee_kind`, `job`, `minimum_grade`, `slot_no`);

ALTER TABLE `fivenet_documents_signature_policies` DROP COLUMN id;
ALTER TABLE `fivenet_documents_signature_policies` ADD PRIMARY KEY (`document_id`);
ALTER TABLE `fivenet_documents_signature_policies` DROP INDEX `idx_fivenet_doc_sig_pol_doc_snap_required`;
ALTER TABLE `fivenet_documents_signature_policies` DROP COLUMN `label`;
ALTER TABLE `fivenet_documents_signature_policies` DROP COLUMN `required`;


ALTER TABLE `fivenet_documents_signature_tasks` ADD COLUMN `label` varchar(120) DEFAULT NULL AFTER `minimum_grade`;

ALTER TABLE `fivenet_documents_signature_policies` ADD COLUMN `rule_kind` SMALLINT(2) DEFAULT 1 AFTER `binding_mode`;
ALTER TABLE `fivenet_documents_signature_policies` ADD COLUMN `required_count` int NOT NULL DEFAULT 1 AFTER `rule_kind`;
ALTER TABLE `fivenet_documents_signature_policies` ADD COLUMN `due_at` datetime(3) DEFAULT NULL AFTER `required_count`;
ALTER TABLE `fivenet_documents_signature_policies` CHANGE `allowed_types_mask` `allowed_types_mask` varchar(120) DEFAULT '[]' NOT NULL AFTER `required_count`;
ALTER TABLE `fivenet_documents_signature_policies` ADD COLUMN `assigned_count` int NOT NULL DEFAULT 0 AFTER `due_at`;
ALTER TABLE `fivenet_documents_signature_policies` ADD COLUMN `approved_count` int NOT NULL DEFAULT 0 AFTER `assigned_count`;
ALTER TABLE `fivenet_documents_signature_policies` ADD COLUMN `declined_count` int NOT NULL DEFAULT 0 AFTER `approved_count`;
ALTER TABLE `fivenet_documents_signature_policies` ADD COLUMN `pending_count` int NOT NULL DEFAULT 0 AFTER `declined_count`;
ALTER TABLE `fivenet_documents_signature_policies` ADD COLUMN `any_declined` tinyint(1) NOT NULL DEFAULT 1 AFTER `pending_count`;
-- started_at and completed_at columns + index
ALTER TABLE `fivenet_documents_signature_policies` ADD COLUMN `started_at` DATETIME(3) AFTER `due_at`;
ALTER TABLE `fivenet_documents_signature_policies` ADD COLUMN `completed_at` DATETIME(3) AFTER `started_at`;
ALTER TABLE `fivenet_documents_signature_policies` ADD KEY `idx_fivenet_doc_sig_pol_started_at` (`started_at`);
ALTER TABLE `fivenet_documents_signature_policies` ADD KEY `idx_fivenet_doc_sig_pol_completed_at` (`completed_at`);

-- Table: fivenet_documents_meta
ALTER TABLE `fivenet_documents_meta` ADD COLUMN `sig_declined_count` int NOT NULL DEFAULT 0 AFTER `sig_required_remaining`;
ALTER TABLE `fivenet_documents_meta` ADD COLUMN `sig_pending_count` int NOT NULL DEFAULT 0 AFTER `sig_declined_count`;
ALTER TABLE `fivenet_documents_meta` ADD COLUMN `sig_any_declined` tinyint(1) NOT NULL DEFAULT 0 AFTER `sig_pending_count`;

COMMIT:
