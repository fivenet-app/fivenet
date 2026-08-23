BEGIN;

ALTER TABLE `fivenet_calendar`
    ADD COLUMN `icon` varchar(128) DEFAULT NULL AFTER `color`;

ALTER TABLE `fivenet_calendar_entries`
    ADD COLUMN `icon` varchar(128) DEFAULT NULL AFTER `rsvp_open`;

COMMIT;
