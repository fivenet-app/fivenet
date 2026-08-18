BEGIN;

ALTER TABLE `fivenet_calendar_entries`
    ADD COLUMN `all_day` tinyint(1) NOT NULL DEFAULT 0 AFTER `end_time`;

UPDATE `fivenet_calendar_entries`
SET `all_day` = 1
WHERE `end_time` IS NULL;

UPDATE `fivenet_calendar_entries` AS `entry`
INNER JOIN `fivenet_calendar` AS `calendar`
    ON `calendar`.`id` = `entry`.`calendar_id`
SET `entry`.`all_day` = 1
WHERE `calendar`.`system_kind` = 1
  AND `entry`.`deleted_at` IS NULL;

COMMIT;
