BEGIN;

ALTER TABLE `fivenet_property_definitions`
  ADD COLUMN `empty_value_display` tinyint unsigned NOT NULL DEFAULT 1 AFTER `reason_mode`;

COMMIT;
