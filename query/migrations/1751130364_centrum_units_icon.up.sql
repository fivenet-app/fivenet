BEGIN;

-- Table: `fivenet_centrum_units` - Add `icon` column
ALTER TABLE `fivenet_centrum_units` ADD `icon` varchar(128) DEFAULT 'MapMarkerIcon' AFTER `color`;

COMMIT;
