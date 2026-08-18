BEGIN;

UPDATE `fivenet_calendar_entries` AS `entry`
INNER JOIN `fivenet_calendar` AS `calendar`
    ON `calendar`.`id` = `entry`.`calendar_id`
SET `entry`.`all_day` = 0
WHERE `calendar`.`system_kind` = 1
  AND `entry`.`deleted_at` IS NULL;

ALTER TABLE `fivenet_calendar_entries`
    DROP COLUMN `all_day`;

COMMIT;
