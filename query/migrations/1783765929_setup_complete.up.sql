BEGIN;

ALTER TABLE `fivenet_config`
    ADD COLUMN `setup_complete` boolean NOT NULL DEFAULT false AFTER `updated_at`;

UPDATE `fivenet_config`
SET `setup_complete` = true;

COMMIT;
