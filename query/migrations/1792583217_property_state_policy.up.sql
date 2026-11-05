BEGIN;

ALTER TABLE `fivenet_property_definitions`
  ADD COLUMN `state_policy` text DEFAULT NULL AFTER `implicit`,
  ADD COLUMN `value_constraints` text DEFAULT NULL AFTER `state_policy`;

COMMIT;
