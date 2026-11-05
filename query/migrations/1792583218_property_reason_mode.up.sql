BEGIN;

ALTER TABLE `fivenet_property_definitions`
  ADD COLUMN `reason_mode` tinyint unsigned NOT NULL DEFAULT 1 AFTER `value_constraints`;

COMMIT;
