BEGIN;

ALTER TABLE `fivenet_user_props` ADD COLUMN `wanted_at` datetime(3) AFTER `wanted`;
ALTER TABLE `fivenet_user_props` ADD COLUMN `wanted_till` datetime(3) AFTER `wanted_at`;

ALTER TABLE `fivenet_user_props` ADD INDEX `idx_wanted_at_till` (`wanted_at`, `wanted_till`);

ALTER TABLE `fivenet_vehicles_props` ADD COLUMN `wanted_at` datetime(3) AFTER `wanted`;
ALTER TABLE `fivenet_vehicles_props` ADD COLUMN `wanted_till` datetime(3) AFTER `wanted_at`;

ALTER TABLE `fivenet_vehicles_props` ADD INDEX `idx_wanted_at_till` (`wanted_at`, `wanted_till`);

COMMIT;
