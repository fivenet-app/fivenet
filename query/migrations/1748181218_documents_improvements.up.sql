BEGIN;

-- Table: `fivenet_documents_pins` - User document pins
ALTER TABLE `fivenet_documents_pins` ADD `user_id` int(11) NULL;
ALTER TABLE `fivenet_documents_pins` CHANGE `user_id` `user_id` int(11) NULL AFTER `job`;

ALTER TABLE `fivenet_documents_pins` DROP INDEX `PRIMARY`;
ALTER TABLE `fivenet_documents_pins` MODIFY COLUMN `job` varchar(20) NULL;
ALTER TABLE `fivenet_documents_pins`
  ADD UNIQUE KEY `idx_document_id_job` (`document_id`, `job`),
  ADD UNIQUE KEY `idx_document_id_user_id` (`document_id`, `user_id`),
  ADD CONSTRAINT `fk_fivenet_documents_pins_user_id` FOREIGN KEY (`user_id`) REFERENCES `{{.UsersTableName}}` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- Table: `fivenet_documents` - Document Draft System
ALTER TABLE `fivenet_documents` DROP INDEX idx_fivenet_documents_content;
ALTER TABLE `fivenet_documents` DROP INDEX idx_fivenet_documents_title;

ALTER TABLE `fivenet_documents` ADD COLUMN `draft` tinyint(1) DEFAULT '0', ALGORITHM=INPLACE;
ALTER TABLE `fivenet_documents` CHANGE `draft` `draft` tinyint(1) DEFAULT '0' AFTER `closed`, ALGORITHM=INPLACE;

ALTER TABLE `fivenet_documents` ADD INDEX `idx_draft` (`draft`);

ALTER TABLE `fivenet_documents`
  ADD FULLTEXT KEY `idx_fivenet_documents_title` (`title`);
ALTER TABLE `fivenet_documents`
  ADD FULLTEXT KEY `idx_fivenet_documents_content` (`content`);

COMMIT;
