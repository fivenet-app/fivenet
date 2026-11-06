BEGIN;

ALTER TABLE `fivenet_property_definitions`
  DROP COLUMN `empty_value_display`;

COMMIT;
