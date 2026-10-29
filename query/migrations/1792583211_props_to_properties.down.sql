BEGIN;

-- This rollback is intentionally best-effort. Custom property values and
-- post-migration changes have no representation in the legacy schema and
-- cannot be restored. Use it only when accepting that data loss.

CREATE TABLE IF NOT EXISTS `fivenet_vehicles_props` (
  `plate` varchar(12) NOT NULL,
  `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3),
  `wanted` tinyint(1) NOT NULL DEFAULT 0,
  `wanted_reason` varchar(255) DEFAULT NULL,
  `wanted_at` datetime(3) DEFAULT NULL,
  `wanted_till` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`plate`),
  KEY `idx_fivenet_vehicles_props_wanted` (`wanted`),
  KEY `idx_wanted_at_till` (`wanted_at`, `wanted_till`),
  CONSTRAINT `fk_fivenet_vehicles_props_plate` FOREIGN KEY (`plate`) REFERENCES `fivenet_vehicle` (`plate`)
    ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

INSERT INTO `fivenet_vehicles_props`
  (`plate`, `updated_at`, `wanted`, `wanted_reason`, `wanted_at`, `wanted_till`)
SELECT value.`plate`, value.`updated_at`, value.`value_bool`, state.`reason`, state.`active_at`, state.`expires_at`
FROM `fivenet_vehicle_property_values` value
LEFT JOIN `fivenet_vehicle_property_state` state
  ON state.`plate` = value.`plate` AND state.`property_key` = value.`property_key`
WHERE value.`property_key` = 'vehicle.wanted';

DELETE FROM `fivenet_vehicle_property_values` WHERE `property_key` = 'vehicle.wanted';
DROP TABLE `fivenet_vehicle_property_state`;

DELETE FROM `fivenet_user_property_values`
WHERE `property_key` IN (
  'user.wanted',
  'user.job',
  'user.job_grade',
  'user.traffic_infraction_points',
  'user.open_fines',
  'user.blood_type',
  'user.mugshot',
  'user.email'
);
DROP TABLE `fivenet_user_property_state`;
DROP TABLE `fivenet_user_files`;

COMMIT;
