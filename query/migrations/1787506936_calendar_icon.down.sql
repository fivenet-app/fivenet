BEGIN;

ALTER TABLE `fivenet_calendar_entries`
    DROP COLUMN `icon`;

ALTER TABLE `fivenet_calendar`
    DROP COLUMN `icon`;

COMMIT;
