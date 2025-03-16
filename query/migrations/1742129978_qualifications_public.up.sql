BEGIN;

ALTER TABLE `fivenet_qualifications` ADD COLUMN `public` tinyint(1) DEFAULT 0 AFTER `closed`;

COMMIT;
