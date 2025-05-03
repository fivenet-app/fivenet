BEGIN;

ALTER TABLE `fivenet_user_props` ADD COLUMN `traffic_infraction_points_updated_at` datetime(3) AFTER `traffic_infraction_points`;

COMMIT;
