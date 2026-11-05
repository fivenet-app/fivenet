BEGIN;

ALTER TABLE `fivenet_property_definitions`
  DROP COLUMN `reason_mode`;

COMMIT;
