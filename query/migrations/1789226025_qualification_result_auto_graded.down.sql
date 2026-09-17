BEGIN;

ALTER TABLE `fivenet_qualifications_results`
    DROP COLUMN `auto_graded`;

COMMIT;
