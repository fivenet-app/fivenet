BEGIN;

ALTER TABLE `fivenet_property_system_overrides`
  DROP COLUMN `group_key`;

COMMIT;
