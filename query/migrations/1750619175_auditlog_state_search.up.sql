BEGIN;

ALTER TABLE `fivenet_audit_log` ADD KEY `idx_state` (`state`);

COMMIT;
