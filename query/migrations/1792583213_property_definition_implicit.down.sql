BEGIN;

ALTER TABLE `fivenet_property_definitions`
  DROP COLUMN `implicit`;

COMMIT;
