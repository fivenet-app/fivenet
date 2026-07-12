BEGIN;

ALTER TABLE `fivenet_config`
    DROP COLUMN `setup_complete`;

COMMIT;
