BEGIN;

RENAME TABLE fivenet_documents_signatures_stamps TO fivenet_documents_stamps;
RENAME TABLE fivenet_documents_signatures_stamps_access TO fivenet_documents_stamps_access;

ALTER TABLE `fivenet_documents_signature_tasks` DROP FOREIGN KEY `fk_fivenet_doc_sigtsks_signature_id`;

DROP TABLE IF EXISTS `fivenet_documents_signature_policies`;
DROP TABLE IF EXISTS `fivenet_documents_signatures`;
DROP TABLE IF EXISTS `fivenet_documents_signature_tasks`;

ALTER TABLE `fivenet_documents_approval_tasks` ADD `signature_required` tinyint(1) DEFAULT 0 NOT NULL AFTER `label`;
ALTER TABLE `fivenet_documents_approval_policies` ADD `signature_required` tinyint(1) DEFAULT 0 NOT NULL AFTER `required_count`;

ALTER TABLE `fivenet_documents_approvals` ADD COLUMN `payload_svg` longtext NOT NULL AFTER `user_job_grade`;
ALTER TABLE `fivenet_documents_approvals` ADD COLUMN `stamp_id` bigint(20) unsigned DEFAULT NULL AFTER `payload_svg`;
ALTER TABLE `fivenet_documents_approvals` ADD CONSTRAINT `fk_fivenet_doc_signatures_stamp_id` FOREIGN KEY (`stamp_id`) REFERENCES `fivenet_documents_stamps` (`id`) ON DELETE SET NULL ON UPDATE CASCADE;

-- Rename SigningService to StampsService: Only for stamp methods
UPDATE `fivenet_rbac_permissions` SET `category` = 'documents.StampsService' WHERE `category` = 'documents.SigningService' AND `name` IN ('UpsertStamp', 'DeleteStamp', 'ListUsableStamps');
-- Remove old Signing service perms
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'documents.SigningService' LIMIT 6;

COMMIT;
