BEGIN;

-- Table: `fivenet_documents_workflow_state`
ALTER TABLE `fivenet_documents_workflow_state` ADD COLUMN `reminder_count` int(11) NOT NULL DEFAULT 0 AFTER `auto_close_time`;

ALTER TABLE `fivenet_documents_workflow_state` ADD KEY `idx_fivenet_documents_workflow_state_reminder_count` (`reminder_count`);

-- Table: `fivenet_documents_workflow_users`
ALTER TABLE `fivenet_documents_workflow_users` ADD COLUMN `reminder_count` int(11) NOT NULL DEFAULT 0 AFTER `manual_reminder_message`;
ALTER TABLE `fivenet_documents_workflow_users` ADD COLUMN `max_reminder_count` int(11) NOT NULL DEFAULT 10 AFTER `reminder_count`;

ALTER TABLE `fivenet_documents_workflow_users` ADD KEY `idx_fivenet_documents_workflow_users_reminder_count` (`reminder_count`, `max_reminder_count`);

COMMIT;
