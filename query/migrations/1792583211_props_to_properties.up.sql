BEGIN;

-- Lifecycle fields for managed vehicle properties. The scalar value itself
-- remains in fivenet_vehicle_property_values, so permissions and filters use
-- the canonical property-value table.
CREATE TABLE IF NOT EXISTS `fivenet_vehicle_property_state` (
  `plate` varchar(12) NOT NULL,
  `property_key` varchar(64) NOT NULL,
  `reason` varchar(255) DEFAULT NULL,
  `active_at` datetime(3) DEFAULT NULL,
  `expires_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`plate`, `property_key`),
  KEY `idx_property_expiry` (`property_key`, `expires_at`),
  KEY `idx_property_active` (`property_key`, `active_at`),
  CONSTRAINT `fk_fivenet_vehicle_property_state_plate` FOREIGN KEY (`plate`) REFERENCES `fivenet_vehicle` (`plate`)
    ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Preserve an explicit false value for every existing vehicle-props row. This
-- lets VehicleProps remain a compatibility projection while the old table is
-- removed.
INSERT INTO `fivenet_vehicle_property_values`
  (`plate`, `property_key`, `value_type`, `value_bool`, `updated_at`)
SELECT `plate`, 'vehicle.wanted', 4, `wanted`, COALESCE(`updated_at`, CURRENT_TIMESTAMP(3))
FROM `fivenet_vehicles_props`
ON DUPLICATE KEY UPDATE
  `value_type` = VALUES(`value_type`),
  `value_bool` = VALUES(`value_bool`),
  `updated_at` = VALUES(`updated_at`);

INSERT INTO `fivenet_vehicle_property_state`
  (`plate`, `property_key`, `reason`, `active_at`, `expires_at`, `updated_at`)
SELECT `plate`, 'vehicle.wanted', `wanted_reason`, `wanted_at`, `wanted_till`, COALESCE(`updated_at`, CURRENT_TIMESTAMP(3))
FROM `fivenet_vehicles_props`
WHERE `wanted` = TRUE
ON DUPLICATE KEY UPDATE
  `reason` = VALUES(`reason`),
  `active_at` = VALUES(`active_at`),
  `expires_at` = VALUES(`expires_at`),
  `updated_at` = VALUES(`updated_at`);

DROP TABLE `fivenet_vehicles_props`;

-- Keep lifecycle data for managed user properties separately from the scalar
-- value, just as for vehicle.wanted.
CREATE TABLE IF NOT EXISTS `fivenet_user_property_state` (
  `user_id` int NOT NULL,
  `property_key` varchar(64) NOT NULL,
  `reason` varchar(255) DEFAULT NULL,
  `active_at` datetime(3) DEFAULT NULL,
  `expires_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`user_id`, `property_key`),
  KEY `idx_property_expiry` (`property_key`, `expires_at`),
  KEY `idx_property_active` (`property_key`, `active_at`),
  CONSTRAINT `fk_fivenet_user_property_state_user_id` FOREIGN KEY (`user_id`) REFERENCES `{{.UsersTableName}}` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Files are relationships, rather than scalar properties. A user may have one
-- file per kind (currently avatar and mugshot), while the file itself retains
-- its normal lifecycle and cleanup behaviour.
CREATE TABLE IF NOT EXISTS `fivenet_user_files` (
  `user_id` int NOT NULL,
  `file_id` bigint unsigned NOT NULL,
  `kind` varchar(32) NOT NULL,
  PRIMARY KEY (`user_id`, `kind`),
  KEY `idx_file_id` (`file_id`),
  CONSTRAINT `fk_fivenet_user_files_user_id` FOREIGN KEY (`user_id`) REFERENCES `{{.UsersTableName}}` (`id`)
    ON DELETE CASCADE,
  CONSTRAINT `fk_fivenet_user_files_file_id` FOREIGN KEY (`file_id`) REFERENCES `fivenet_files` (`id`)
    ON DELETE RESTRICT
) ENGINE=InnoDB;

-- Migrate the user wanted flag into the canonical user-property value table.
INSERT INTO `fivenet_user_property_values`
  (`user_id`, `property_key`, `value_type`, `value_bool`, `updated_at`)
SELECT `user_id`, 'user.wanted', 4, `wanted`, COALESCE(`updated_at`, CURRENT_TIMESTAMP(3))
FROM `fivenet_user_props`
ON DUPLICATE KEY UPDATE
  `value_type` = VALUES(`value_type`),
  `value_bool` = VALUES(`value_bool`),
  `updated_at` = VALUES(`updated_at`);

-- The former traffic_infraction_points_updated_at field is the managed
-- property's state timestamp. No active/expiry fields apply to points.
INSERT INTO `fivenet_user_property_state`
  (`user_id`, `property_key`, `updated_at`)
SELECT `user_id`, 'user.traffic_infraction_points',
  COALESCE(`traffic_infraction_points_updated_at`, `updated_at`, CURRENT_TIMESTAMP(3))
FROM `fivenet_user_props`
WHERE `traffic_infraction_points` IS NOT NULL
ON DUPLICATE KEY UPDATE
  `updated_at` = VALUES(`updated_at`);

INSERT INTO `fivenet_user_property_state`
  (`user_id`, `property_key`, `active_at`, `expires_at`, `updated_at`)
SELECT `user_id`, 'user.wanted', `wanted_at`, `wanted_till`, COALESCE(`updated_at`, CURRENT_TIMESTAMP(3))
FROM `fivenet_user_props`
WHERE `wanted` = TRUE
ON DUPLICATE KEY UPDATE
  `active_at` = VALUES(`active_at`),
  `expires_at` = VALUES(`expires_at`),
  `updated_at` = VALUES(`updated_at`);

-- Migrate the remaining legacy fields that have registered user-property keys.
INSERT INTO `fivenet_user_property_values`
  (`user_id`, `property_key`, `value_type`, `value_string`, `updated_at`)
SELECT `user_id`, 'user.job', 1, `job`, COALESCE(`updated_at`, CURRENT_TIMESTAMP(3))
FROM `fivenet_user_props`
WHERE `job` IS NOT NULL
ON DUPLICATE KEY UPDATE
  `value_type` = VALUES(`value_type`),
  `value_string` = VALUES(`value_string`),
  `updated_at` = VALUES(`updated_at`);

INSERT INTO `fivenet_user_property_values`
  (`user_id`, `property_key`, `value_type`, `value_integer`, `updated_at`)
SELECT `user_id`, 'user.job_grade', 2, `job_grade`, COALESCE(`updated_at`, CURRENT_TIMESTAMP(3))
FROM `fivenet_user_props`
WHERE `job_grade` IS NOT NULL
ON DUPLICATE KEY UPDATE
  `value_type` = VALUES(`value_type`),
  `value_integer` = VALUES(`value_integer`),
  `updated_at` = VALUES(`updated_at`);

INSERT INTO `fivenet_user_property_values`
  (`user_id`, `property_key`, `value_type`, `value_integer`, `updated_at`)
SELECT `user_id`, 'user.traffic_infraction_points', 2, `traffic_infraction_points`,
  COALESCE(`traffic_infraction_points_updated_at`, `updated_at`, CURRENT_TIMESTAMP(3))
FROM `fivenet_user_props`
WHERE `traffic_infraction_points` IS NOT NULL
ON DUPLICATE KEY UPDATE
  `value_type` = VALUES(`value_type`),
  `value_integer` = VALUES(`value_integer`),
  `updated_at` = VALUES(`updated_at`);

INSERT INTO `fivenet_user_files` (`user_id`, `file_id`, `kind`)
SELECT `user_id`, `avatar_file_id`, 'avatar'
FROM `fivenet_user_props`
WHERE `avatar_file_id` IS NOT NULL
ON DUPLICATE KEY UPDATE `file_id` = VALUES(`file_id`);

INSERT INTO `fivenet_user_files` (`user_id`, `file_id`, `kind`)
SELECT `user_id`, `mugshot_file_id`, 'mugshot'
FROM `fivenet_user_props`
WHERE `mugshot_file_id` IS NOT NULL
ON DUPLICATE KEY UPDATE `file_id` = VALUES(`file_id`);

INSERT INTO `fivenet_user_property_values`
  (`user_id`, `property_key`, `value_type`, `value_integer`, `updated_at`)
SELECT `user_id`, 'user.open_fines', 2, `open_fines`, COALESCE(`updated_at`, CURRENT_TIMESTAMP(3))
FROM `fivenet_user_props`
WHERE `open_fines` IS NOT NULL
ON DUPLICATE KEY UPDATE
  `value_type` = VALUES(`value_type`),
  `value_integer` = VALUES(`value_integer`),
  `updated_at` = VALUES(`updated_at`);

INSERT INTO `fivenet_user_property_values`
  (`user_id`, `property_key`, `value_type`, `value_string`, `updated_at`)
SELECT `user_id`, 'user.blood_type', 1, `blood_type`, COALESCE(`updated_at`, CURRENT_TIMESTAMP(3))
FROM `fivenet_user_props`
WHERE `blood_type` IS NOT NULL
ON DUPLICATE KEY UPDATE
  `value_type` = VALUES(`value_type`),
  `value_string` = VALUES(`value_string`),
  `updated_at` = VALUES(`updated_at`);

INSERT INTO `fivenet_user_property_values`
  (`user_id`, `property_key`, `value_type`, `value_string`, `updated_at`)
SELECT `user_id`, 'user.mugshot', 1, `mug_shot`, COALESCE(`updated_at`, CURRENT_TIMESTAMP(3))
FROM `fivenet_user_props`
WHERE `mug_shot` IS NOT NULL
ON DUPLICATE KEY UPDATE
  `value_type` = VALUES(`value_type`),
  `value_string` = VALUES(`value_string`),
  `updated_at` = VALUES(`updated_at`);

INSERT INTO `fivenet_user_property_values`
  (`user_id`, `property_key`, `value_type`, `value_string`, `updated_at`)
SELECT `user_id`, 'user.email', 1, `email`, COALESCE(`updated_at`, CURRENT_TIMESTAMP(3))
FROM `fivenet_user_props`
WHERE `email` IS NOT NULL
ON DUPLICATE KEY UPDATE
  `value_type` = VALUES(`value_type`),
  `value_string` = VALUES(`value_string`),
  `updated_at` = VALUES(`updated_at`);

-- Keep the legacy projection during the UserProps compatibility window.

COMMIT;
