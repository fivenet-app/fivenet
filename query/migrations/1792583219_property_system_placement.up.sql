BEGIN;

ALTER TABLE `fivenet_property_system_overrides`
  ADD COLUMN `group_key` varchar(64) DEFAULT NULL AFTER `enabled`;

COMMIT;
