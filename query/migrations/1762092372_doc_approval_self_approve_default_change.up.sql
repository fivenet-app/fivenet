BEGIN;

ALTER TABLE `fivenet_documents_approval_policies` MODIFY COLUMN `self_approve_allowed` tinyint(1) DEFAULT 0 NOT NULL;

COMMIT;
