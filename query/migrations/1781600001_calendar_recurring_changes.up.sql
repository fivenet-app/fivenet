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

-- Remove current RSVP occurences
DELETE FROM `fivenet_calendar_rsvp_occurrence`;

ALTER TABLE `fivenet_calendar_rsvp_occurrence`
  ADD COLUMN `recurrence_id` DATETIME(3) NOT NULL,
  ADD COLUMN `recurrence_version` INT NOT NULL DEFAULT 1,
  ADD UNIQUE KEY `uq_fivenet_calendar_rsvp_occurrence` (
    `entry_id`,
    `recurrence_version`,
    `recurrence_id`,
    `user_id`
  ),
  ADD INDEX `idx_calendar_rsvp_occurrence_cleanup` (
    `entry_id`,
    `recurrence_version`,
    `created_at`
  );

ALTER TABLE `fivenet_calendar_entries`
  ADD COLUMN `recurrence_version` INT NOT NULL DEFAULT 1;

COMMIT;
