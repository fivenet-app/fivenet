BEGIN;

ALTER TABLE `fivenet_property_definitions`
  ADD COLUMN `implicit` tinyint(1) NOT NULL DEFAULT 0 AFTER `display_mode`;

COMMIT;
