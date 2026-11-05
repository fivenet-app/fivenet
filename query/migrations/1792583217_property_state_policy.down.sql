BEGIN;

ALTER TABLE `fivenet_property_definitions`
  DROP COLUMN `value_constraints`,
  DROP COLUMN `state_policy`;

COMMIT;
