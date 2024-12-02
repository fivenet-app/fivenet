BEGIN;

ALTER TABLE `fivenet_documents_templates` ADD `workflow` mediumtext NULL AFTER `schema`;

-- Table: fivenet_documents_workflow_state
CREATE TABLE `fivenet_documents_workflow_state` (
  `document_id` bigint unsigned NOT NULL,
  `next_reminder_time` datetime(3) DEFAULT NULL,
  `next_reminder_count` int(5) DEFAULT NULL,
  `auto_close_time` datetime(3) DEFAULT NULL,
  UNIQUE KEY `idx_fivenet_documents_workflow_state_document_id` (`document_id`),
  KEY `idx_fivenet_documents_workflow_state_next_reminder_time` (`next_reminder_time`),
  KEY `idx_fivenet_documents_workflow_state_auto_close_time` (`auto_close_time`),
  CONSTRAINT `fk_fivenet_documents_workflow_state_document_id` FOREIGN KEY (`document_id`) REFERENCES `fivenet_documents` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Table: fivenet_documents_workflow_users
CREATE TABLE `fivenet_documents_workflow_users` (
  `document_id` bigint unsigned NOT NULL,
  `user_id` int(11) NOT NULL,
  `manual_reminder_time` datetime(3) DEFAULT NULL,
  `manual_reminder_message` varchar(255) DEFAULT NULL,
  UNIQUE KEY `idx_fivenet_documents_workflow_users_document_id` (`document_id`, `user_id`),
  KEY `idx_fivenet_documents_workflow_users_manual_reminder_time` (`manual_reminder_time`),
  CONSTRAINT `fk_fivenet_documents_workflow_users_document_id` FOREIGN KEY (`document_id`) REFERENCES `fivenet_documents` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_documents_workflow_users_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

COMMIT;
