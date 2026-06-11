BEGIN;

ALTER TABLE `fivenet_calendar_entries`
    ADD COLUMN `recurring_until` datetime(3) DEFAULT NULL AFTER `recurring`;

ALTER TABLE `fivenet_calendar_entries`
    ADD INDEX `idx_fivenet_calendar_entries_calendar_start_end` (
      `calendar_id`,
      `start_time`,
      `end_time`
    ),
    ADD INDEX `idx_fivenet_calendar_entries_calendar_recurring_until_start` (
      `calendar_id`,
      `recurring_until`,
      `start_time`
    );

COMMIT;
