BEGIN;

ALTER TABLE `fivenet_qualifications_results`
    ADD COLUMN `auto_graded` TINYINT(1) NOT NULL DEFAULT 0 AFTER `summary`;

COMMIT;
