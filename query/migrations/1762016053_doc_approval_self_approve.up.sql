BEGIN;

ALTER TABLE `fivenet_documents_approval_policies` ADD COLUMN `self_approve_allowed` tinyint(1) NOT NULL DEFAULT 1 AFTER `signature_required`;

COMMIT;
