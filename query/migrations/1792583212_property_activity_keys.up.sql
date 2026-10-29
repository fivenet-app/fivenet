BEGIN;

-- Generic property activity is keyed independently from its structured data so
-- activity feeds can filter it without database-specific JSON expressions.
ALTER TABLE `fivenet_user_activity`
  ADD COLUMN `key` varchar(64) DEFAULT NULL AFTER `type`,
  ADD KEY `idx_fivenet_user_activity_target_key_created_at` (`target_user_id`, `key`, `created_at`);

ALTER TABLE `fivenet_vehicle_activity`
  ADD COLUMN `key` varchar(64) DEFAULT NULL AFTER `type`,
  ADD KEY `idx_fivenet_vehicle_activity_plate_key_created_at` (`plate`, `key`, `created_at`);

COMMIT;
