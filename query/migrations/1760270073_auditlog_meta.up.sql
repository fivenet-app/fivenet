BEGIN;

ALTER TABLE `fivenet_audit_log` ADD COLUMN `result` smallint NOT NULL AFTER `state`,
	ADD COLUMN `meta` varchar(256) AFTER `result`;

ALTER TABLE `fivenet_audit_log` RENAME COLUMN `state` TO `action`;
ALTER TABLE `fivenet_audit_log` RENAME INDEX `idx_state` TO `idx_action`;

ALTER TABLE `fivenet_audit_log` ADD INDEX `idx_result` (`result`);

UPDATE `fivenet_audit_log` SET `action` = 2, `result` = 3 WHERE `action` = 1 AND `result` = 0;
UPDATE `fivenet_audit_log` SET `result` = 1 WHERE `result` = 0;

UPDATE `fivenet_audit_log` SET `data` = '{"aud.msg":"No request data"}' WHERE `data` = 'No data' OR `data` = '{}';
UPDATE `fivenet_audit_log` SET `data` = '{"aud.err":"Failed to marshal data"}' WHERE `data` = 'Failed to marshal data';

COMMIT;
