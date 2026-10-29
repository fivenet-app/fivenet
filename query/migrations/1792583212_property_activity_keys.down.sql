BEGIN;

ALTER TABLE `fivenet_vehicle_activity`
  DROP KEY `idx_fivenet_vehicle_activity_plate_key_created_at`,
  DROP COLUMN `key`;

ALTER TABLE `fivenet_user_activity`
  DROP KEY `idx_fivenet_user_activity_target_key_created_at`,
  DROP COLUMN `key`;

COMMIT;
