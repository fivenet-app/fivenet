BEGIN;

-- GROUP_CONCAT is used only for the bounded legacy system-property mappings
-- below. Custom property definitions are intentionally not copied into role
-- grants by this migration, and the current system-property lists are small
-- enough for the server's default group_concat_max_len.

-- Property-backed fields are represented by namespaced system property keys.
-- Legacy Fields values are converted to typed, namespaced property keys.
SET @user_system_properties := '{"propertyList":{"keys":["user.phone_number","user.licenses","user.blood_type","user.email","user.job","user.labels","user.mugshot","user.open_fines","user.traffic_infraction_points","user.wanted"]}}';
SET @vehicle_system_properties := '{"propertyList":{"keys":["vehicle.wanted"]}}';

-- User/citizen properties: create the typed property-list attribute.
INSERT INTO `fivenet_rbac_attrs` (`permission_id`, `key`, `type`, `valid_values`)
SELECT p.`id`, 'Properties', 'PropertyList', @user_system_properties
FROM `fivenet_rbac_permissions` p
WHERE p.`namespace` = 'citizens'
  AND p.`service` = 'CitizensService'
  AND p.`name` = 'ListCitizens'
  AND NOT EXISTS (
    SELECT 1 FROM `fivenet_rbac_attrs` a
    WHERE a.`permission_id` = p.`id` AND a.`key` = 'Properties'
  );

UPDATE `fivenet_rbac_attrs` a
INNER JOIN `fivenet_rbac_permissions` p ON p.`id` = a.`permission_id`
SET a.`type` = 'PropertyList', a.`valid_values` = @user_system_properties
WHERE p.`namespace` = 'citizens'
  AND p.`service` = 'CitizensService'
  AND p.`name` = 'ListCitizens'
  AND a.`key` = 'Properties';

-- Move each role's legacy citizen property grants into the new property-list
-- attribute without widening access.
INSERT INTO `fivenet_rbac_roles_attrs`
  (`role_id`, `created_at`, `updated_at`, `attr_id`, `value`)
SELECT old_value.`role_id`, old_value.`created_at`, old_value.`updated_at`,
       new_attr.`id`, CONCAT(
         '{"propertyList":{"keys":[',
         COALESCE(GROUP_CONCAT(JSON_QUOTE(property_map.`property_key`)
           ORDER BY property_map.`property_key` SEPARATOR ','), ''),
         ']}}'
       )
FROM `fivenet_rbac_roles_attrs` old_value
INNER JOIN `fivenet_rbac_attrs` old_attr ON old_attr.`id` = old_value.`attr_id`
INNER JOIN `fivenet_rbac_permissions` perm ON perm.`id` = old_attr.`permission_id`
INNER JOIN `fivenet_rbac_attrs` new_attr
  ON new_attr.`permission_id` = perm.`id` AND new_attr.`key` = 'Properties'
LEFT JOIN (
  SELECT 'PhoneNumber' AS `legacy_value`, 'user.phone_number' AS `property_key`
  UNION ALL SELECT 'Licenses', 'user.licenses'
  UNION ALL SELECT 'UserProps.BloodType', 'user.blood_type'
  UNION ALL SELECT 'UserProps.Email', 'user.email'
  UNION ALL SELECT 'UserProps.Job', 'user.job'
  UNION ALL SELECT 'UserProps.Labels', 'user.labels'
  UNION ALL SELECT 'UserProps.Mugshot', 'user.mugshot'
  UNION ALL SELECT 'UserProps.OpenFines', 'user.open_fines'
  UNION ALL SELECT 'UserProps.TrafficInfractionPoints', 'user.traffic_infraction_points'
  UNION ALL SELECT 'UserProps.Wanted', 'user.wanted'
) property_map ON JSON_CONTAINS(
  old_value.`value`, JSON_QUOTE(property_map.`legacy_value`), '$.stringList.strings'
)
WHERE old_attr.`key` = 'Fields'
  AND old_attr.`type` = 'StringList'
  AND perm.`namespace` = 'citizens'
  AND perm.`service` = 'CitizensService'
  AND perm.`name` = 'ListCitizens'
GROUP BY old_value.`role_id`, old_value.`created_at`, old_value.`updated_at`, new_attr.`id`
ON DUPLICATE KEY UPDATE
  `value` = VALUES(`value`), `updated_at` = VALUES(`updated_at`);

INSERT INTO `fivenet_rbac_job_attrs` (`job`, `attr_id`, `max_values`)
SELECT old_job_attr.`job`, new_attr.`id`, old_job_attr.`max_values`
FROM `fivenet_rbac_job_attrs` old_job_attr
INNER JOIN `fivenet_rbac_attrs` old_attr ON old_attr.`id` = old_job_attr.`attr_id`
INNER JOIN `fivenet_rbac_permissions` perm ON perm.`id` = old_attr.`permission_id`
INNER JOIN `fivenet_rbac_attrs` new_attr
  ON new_attr.`permission_id` = perm.`id` AND new_attr.`key` = 'Properties'
WHERE old_attr.`key` = 'Fields'
  AND old_attr.`type` = 'StringList'
  AND perm.`namespace` = 'citizens'
  AND perm.`service` = 'CitizensService'
  AND perm.`name` = 'ListCitizens'
ON DUPLICATE KEY UPDATE `max_values` = VALUES(`max_values`);

DELETE role_attr
FROM `fivenet_rbac_roles_attrs` role_attr
INNER JOIN `fivenet_rbac_attrs` attr ON attr.`id` = role_attr.`attr_id`
INNER JOIN `fivenet_rbac_permissions` perm ON perm.`id` = attr.`permission_id`
WHERE perm.`namespace` = 'citizens' AND perm.`service` = 'CitizensService'
  AND perm.`name` = 'ListCitizens' AND attr.`key` = 'Fields';
DELETE job_attr
FROM `fivenet_rbac_job_attrs` job_attr
INNER JOIN `fivenet_rbac_attrs` attr ON attr.`id` = job_attr.`attr_id`
INNER JOIN `fivenet_rbac_permissions` perm ON perm.`id` = attr.`permission_id`
WHERE perm.`namespace` = 'citizens' AND perm.`service` = 'CitizensService'
  AND perm.`name` = 'ListCitizens' AND attr.`key` = 'Fields';
DELETE attr
FROM `fivenet_rbac_attrs` attr
INNER JOIN `fivenet_rbac_permissions` perm ON perm.`id` = attr.`permission_id`
WHERE perm.`namespace` = 'citizens' AND perm.`service` = 'CitizensService'
  AND perm.`name` = 'ListCitizens' AND attr.`key` = 'Fields';

-- The generated permissions already use Properties, while older installations
-- still have a Fields attribute. Merge the old attribute into the existing
-- Properties row instead of renaming it into the unique-key collision.
INSERT INTO `fivenet_rbac_attrs` (`permission_id`, `key`, `type`, `valid_values`)
SELECT p.`id`, 'Properties', 'PropertyList', @user_system_properties
FROM `fivenet_rbac_permissions` p
WHERE p.`namespace` = 'citizens'
  AND p.`service` = 'CitizensService'
  AND p.`name` IN ('SetUserProps', 'GetUserProperties', 'SetUserProperty', 'DeleteUserProperty')
  AND NOT EXISTS (
    SELECT 1 FROM `fivenet_rbac_attrs` existing
    WHERE existing.`permission_id` = p.`id` AND existing.`key` = 'Properties'
  );

INSERT INTO `fivenet_rbac_roles_attrs`
  (`role_id`, `created_at`, `updated_at`, `attr_id`, `value`)
SELECT old_value.`role_id`, old_value.`created_at`, old_value.`updated_at`,
       new_attr.`id`, CONCAT(
         '{"propertyList":{"keys":[',
         COALESCE(GROUP_CONCAT(JSON_QUOTE(property_map.`property_key`)
           ORDER BY property_map.`property_key` SEPARATOR ','), ''),
         ']}}'
       )
FROM `fivenet_rbac_roles_attrs` old_value
INNER JOIN `fivenet_rbac_attrs` old_attr ON old_attr.`id` = old_value.`attr_id`
INNER JOIN `fivenet_rbac_permissions` perm ON perm.`id` = old_attr.`permission_id`
INNER JOIN `fivenet_rbac_attrs` new_attr
  ON new_attr.`permission_id` = perm.`id` AND new_attr.`key` = 'Properties'
LEFT JOIN (
  SELECT 'Wanted' AS `legacy_value`, 'user.wanted' AS `property_key`
  UNION ALL SELECT 'PhoneNumber', 'user.phone_number'
  UNION ALL SELECT 'Licenses', 'user.licenses'
  UNION ALL SELECT 'UserProps.BloodType', 'user.blood_type'
  UNION ALL SELECT 'UserProps.Email', 'user.email'
  UNION ALL SELECT 'UserProps.Job', 'user.job'
  UNION ALL SELECT 'UserProps.Labels', 'user.labels'
  UNION ALL SELECT 'UserProps.Mugshot', 'user.mugshot'
  UNION ALL SELECT 'UserProps.OpenFines', 'user.open_fines'
  UNION ALL SELECT 'UserProps.TrafficInfractionPoints', 'user.traffic_infraction_points'
) property_map ON JSON_CONTAINS(
  old_value.`value`, JSON_QUOTE(property_map.`legacy_value`), '$.stringList.strings'
)
WHERE old_attr.`key` = 'Fields'
  AND old_attr.`type` = 'StringList'
  AND perm.`namespace` = 'citizens'
  AND perm.`service` = 'CitizensService'
  AND perm.`name` IN ('SetUserProps', 'GetUserProperties', 'SetUserProperty', 'DeleteUserProperty')
GROUP BY old_value.`role_id`, old_value.`created_at`, old_value.`updated_at`, new_attr.`id`
ON DUPLICATE KEY UPDATE
  `value` = VALUES(`value`), `updated_at` = VALUES(`updated_at`);

INSERT INTO `fivenet_rbac_job_attrs` (`job`, `attr_id`, `max_values`)
SELECT old_job_attr.`job`, new_attr.`id`, old_job_attr.`max_values`
FROM `fivenet_rbac_job_attrs` old_job_attr
INNER JOIN `fivenet_rbac_attrs` old_attr ON old_attr.`id` = old_job_attr.`attr_id`
INNER JOIN `fivenet_rbac_permissions` perm ON perm.`id` = old_attr.`permission_id`
INNER JOIN `fivenet_rbac_attrs` new_attr
  ON new_attr.`permission_id` = perm.`id` AND new_attr.`key` = 'Properties'
WHERE old_attr.`key` = 'Fields'
  AND old_attr.`type` = 'StringList'
  AND perm.`namespace` = 'citizens'
  AND perm.`service` = 'CitizensService'
  AND perm.`name` IN ('SetUserProps', 'GetUserProperties', 'SetUserProperty', 'DeleteUserProperty')
ON DUPLICATE KEY UPDATE `max_values` = VALUES(`max_values`);

UPDATE `fivenet_rbac_attrs` a
INNER JOIN `fivenet_rbac_permissions` p ON p.`id` = a.`permission_id`
SET a.`type` = 'PropertyList', a.`valid_values` = @user_system_properties
WHERE p.`namespace` = 'citizens'
  AND p.`service` = 'CitizensService'
  AND p.`name` IN ('SetUserProps', 'GetUserProperties', 'SetUserProperty', 'DeleteUserProperty')
  AND a.`key` = 'Properties';

DELETE role_attr
FROM `fivenet_rbac_roles_attrs` role_attr
INNER JOIN `fivenet_rbac_attrs` attr ON attr.`id` = role_attr.`attr_id`
INNER JOIN `fivenet_rbac_permissions` perm ON perm.`id` = attr.`permission_id`
WHERE perm.`namespace` = 'citizens' AND perm.`service` = 'CitizensService'
  AND perm.`name` IN ('SetUserProps', 'GetUserProperties', 'SetUserProperty', 'DeleteUserProperty')
  AND attr.`key` = 'Fields';
DELETE job_attr
FROM `fivenet_rbac_job_attrs` job_attr
INNER JOIN `fivenet_rbac_attrs` attr ON attr.`id` = job_attr.`attr_id`
INNER JOIN `fivenet_rbac_permissions` perm ON perm.`id` = attr.`permission_id`
WHERE perm.`namespace` = 'citizens' AND perm.`service` = 'CitizensService'
  AND perm.`name` IN ('SetUserProps', 'GetUserProperties', 'SetUserProperty', 'DeleteUserProperty')
  AND attr.`key` = 'Fields';
DELETE attr
FROM `fivenet_rbac_attrs` attr
INNER JOIN `fivenet_rbac_permissions` perm ON perm.`id` = attr.`permission_id`
WHERE perm.`namespace` = 'citizens' AND perm.`service` = 'CitizensService'
  AND perm.`name` IN ('SetUserProps', 'GetUserProperties', 'SetUserProperty', 'DeleteUserProperty')
  AND attr.`key` = 'Fields';

-- Vehicle property endpoints now use Properties/PropertyList.
UPDATE `fivenet_rbac_attrs` a
INNER JOIN `fivenet_rbac_permissions` p ON p.`id` = a.`permission_id`
SET a.`valid_values` = @vehicle_system_properties
WHERE p.`namespace` = 'vehicles'
  AND p.`service` = 'VehiclesService'
  AND p.`name` = 'ListVehicles'
  AND a.`key` = 'Properties';

INSERT INTO `fivenet_rbac_attrs` (`permission_id`, `key`, `type`, `valid_values`)
SELECT p.`id`, 'Properties', 'PropertyList', @vehicle_system_properties
FROM `fivenet_rbac_permissions` p
WHERE p.`namespace` = 'vehicles'
  AND p.`service` = 'VehiclesService'
  AND p.`name` IN ('SetVehicleProperty', 'DeleteVehicleProperty')
  AND NOT EXISTS (
    SELECT 1 FROM `fivenet_rbac_attrs` existing
    WHERE existing.`permission_id` = p.`id` AND existing.`key` = 'Properties'
  );

INSERT INTO `fivenet_rbac_roles_attrs`
  (`role_id`, `created_at`, `updated_at`, `attr_id`, `value`)
SELECT old_value.`role_id`, old_value.`created_at`, old_value.`updated_at`,
       new_attr.`id`, '{"propertyList":{"keys":["vehicle.wanted"]}}'
FROM `fivenet_rbac_roles_attrs` old_value
INNER JOIN `fivenet_rbac_attrs` old_attr ON old_attr.`id` = old_value.`attr_id`
INNER JOIN `fivenet_rbac_permissions` perm ON perm.`id` = old_attr.`permission_id`
INNER JOIN `fivenet_rbac_attrs` new_attr
  ON new_attr.`permission_id` = perm.`id` AND new_attr.`key` = 'Properties'
WHERE old_attr.`key` = 'Fields'
  AND old_attr.`type` = 'StringList'
  AND perm.`namespace` = 'vehicles'
  AND perm.`service` = 'VehiclesService'
  AND perm.`name` IN ('SetVehicleProperty', 'DeleteVehicleProperty')
  AND JSON_CONTAINS(old_value.`value`, '"Wanted"', '$.stringList.strings')
ON DUPLICATE KEY UPDATE
  `value` = VALUES(`value`), `updated_at` = VALUES(`updated_at`);

INSERT INTO `fivenet_rbac_job_attrs` (`job`, `attr_id`, `max_values`)
SELECT old_job_attr.`job`, new_attr.`id`, old_job_attr.`max_values`
FROM `fivenet_rbac_job_attrs` old_job_attr
INNER JOIN `fivenet_rbac_attrs` old_attr ON old_attr.`id` = old_job_attr.`attr_id`
INNER JOIN `fivenet_rbac_permissions` perm ON perm.`id` = old_attr.`permission_id`
INNER JOIN `fivenet_rbac_attrs` new_attr
  ON new_attr.`permission_id` = perm.`id` AND new_attr.`key` = 'Properties'
WHERE old_attr.`key` = 'Fields'
  AND old_attr.`type` = 'StringList'
  AND perm.`namespace` = 'vehicles'
  AND perm.`service` = 'VehiclesService'
  AND perm.`name` IN ('SetVehicleProperty', 'DeleteVehicleProperty')
ON DUPLICATE KEY UPDATE `max_values` = VALUES(`max_values`);

UPDATE `fivenet_rbac_attrs` a
INNER JOIN `fivenet_rbac_permissions` p ON p.`id` = a.`permission_id`
SET a.`type` = 'PropertyList',
    a.`valid_values` = @vehicle_system_properties
WHERE p.`namespace` = 'vehicles'
  AND p.`service` = 'VehiclesService'
  AND p.`name` IN ('SetVehicleProperty', 'DeleteVehicleProperty')
  AND a.`key` = 'Properties';

DELETE role_attr
FROM `fivenet_rbac_roles_attrs` role_attr
INNER JOIN `fivenet_rbac_attrs` attr ON attr.`id` = role_attr.`attr_id`
INNER JOIN `fivenet_rbac_permissions` perm ON perm.`id` = attr.`permission_id`
WHERE perm.`namespace` = 'vehicles' AND perm.`service` = 'VehiclesService'
  AND perm.`name` IN ('SetVehicleProperty', 'DeleteVehicleProperty')
  AND attr.`key` = 'Fields';
DELETE job_attr
FROM `fivenet_rbac_job_attrs` job_attr
INNER JOIN `fivenet_rbac_attrs` attr ON attr.`id` = job_attr.`attr_id`
INNER JOIN `fivenet_rbac_permissions` perm ON perm.`id` = attr.`permission_id`
WHERE perm.`namespace` = 'vehicles' AND perm.`service` = 'VehiclesService'
  AND perm.`name` IN ('SetVehicleProperty', 'DeleteVehicleProperty')
  AND attr.`key` = 'Fields';
DELETE attr
FROM `fivenet_rbac_attrs` attr
INNER JOIN `fivenet_rbac_permissions` perm ON perm.`id` = attr.`permission_id`
WHERE perm.`namespace` = 'vehicles' AND perm.`service` = 'VehiclesService'
  AND perm.`name` IN ('SetVehicleProperty', 'DeleteVehicleProperty')
  AND attr.`key` = 'Fields';

INSERT INTO `fivenet_rbac_attrs` (`permission_id`, `key`, `type`, `valid_values`)
SELECT p.`id`, 'Properties', 'PropertyList', @vehicle_system_properties
FROM `fivenet_rbac_permissions` p
WHERE p.`namespace` = 'vehicles'
  AND p.`service` = 'VehiclesService'
  AND p.`name` = 'ListVehicles'
  AND NOT EXISTS (
    SELECT 1 FROM `fivenet_rbac_attrs` existing
    WHERE existing.`permission_id` = p.`id` AND existing.`key` = 'Properties'
  );

INSERT INTO `fivenet_rbac_attrs` (`permission_id`, `key`, `type`, `valid_values`)
SELECT p.`id`, 'Properties', 'PropertyList', @vehicle_system_properties
FROM `fivenet_rbac_permissions` p
WHERE p.`namespace` = 'vehicles'
  AND p.`service` = 'VehiclesService'
  AND p.`name` = 'SetVehicleProps'
  AND NOT EXISTS (
    SELECT 1 FROM `fivenet_rbac_attrs` existing
    WHERE existing.`permission_id` = p.`id` AND existing.`key` = 'Properties'
  );

INSERT INTO `fivenet_rbac_job_attrs` (`job`, `attr_id`, `max_values`)
SELECT old_job_attr.`job`, new_attr.`id`, old_job_attr.`max_values`
FROM `fivenet_rbac_job_attrs` old_job_attr
INNER JOIN `fivenet_rbac_attrs` old_attr ON old_attr.`id` = old_job_attr.`attr_id`
INNER JOIN `fivenet_rbac_permissions` perm ON perm.`id` = old_attr.`permission_id`
INNER JOIN `fivenet_rbac_attrs` new_attr
  ON new_attr.`permission_id` = perm.`id` AND new_attr.`key` = 'Properties'
WHERE old_attr.`key` = 'Fields'
  AND old_attr.`type` = 'StringList'
  AND perm.`namespace` = 'vehicles'
  AND perm.`service` = 'VehiclesService'
  AND perm.`name` IN ('ListVehicles', 'SetVehicleProps')
ON DUPLICATE KEY UPDATE `max_values` = VALUES(`max_values`);

-- Migrate the fixed Wanted grant for ListVehicles when it had no generic
-- Properties marker.
INSERT INTO `fivenet_rbac_roles_attrs`
  (`role_id`, `created_at`, `updated_at`, `attr_id`, `value`)
SELECT old_value.`role_id`, old_value.`created_at`, old_value.`updated_at`,
       property_attr.`id`, '{"propertyList":{"keys":["vehicle.wanted"]}}'
FROM `fivenet_rbac_roles_attrs` old_value
INNER JOIN `fivenet_rbac_attrs` old_attr ON old_attr.`id` = old_value.`attr_id`
INNER JOIN `fivenet_rbac_permissions` perm ON perm.`id` = old_attr.`permission_id`
INNER JOIN `fivenet_rbac_attrs` property_attr
  ON property_attr.`permission_id` = perm.`id` AND property_attr.`key` = 'Properties'
WHERE old_attr.`key` = 'Fields'
  AND perm.`namespace` = 'vehicles' AND perm.`service` = 'VehiclesService'
  AND perm.`name` = 'ListVehicles'
  AND JSON_CONTAINS(old_value.`value`, '"Wanted"', '$.stringList.strings')
  AND NOT JSON_CONTAINS(old_value.`value`, '"Properties"', '$.stringList.strings')
ON DUPLICATE KEY UPDATE `value` = VALUES(`value`), `updated_at` = VALUES(`updated_at`);

-- SetVehicleProps now authorizes Wanted through its typed property list.
INSERT INTO `fivenet_rbac_roles_attrs`
  (`role_id`, `created_at`, `updated_at`, `attr_id`, `value`)
SELECT old_value.`role_id`, old_value.`created_at`, old_value.`updated_at`,
       property_attr.`id`, '{"propertyList":{"keys":["vehicle.wanted"]}}'
FROM `fivenet_rbac_roles_attrs` old_value
INNER JOIN `fivenet_rbac_attrs` old_attr ON old_attr.`id` = old_value.`attr_id`
INNER JOIN `fivenet_rbac_permissions` perm ON perm.`id` = old_attr.`permission_id`
INNER JOIN `fivenet_rbac_attrs` property_attr
  ON property_attr.`permission_id` = perm.`id` AND property_attr.`key` = 'Properties'
WHERE old_attr.`key` = 'Fields'
  AND perm.`namespace` = 'vehicles' AND perm.`service` = 'VehiclesService'
  AND perm.`name` = 'SetVehicleProps'
  AND JSON_CONTAINS(old_value.`value`, '"Wanted"', '$.stringList.strings')
ON DUPLICATE KEY UPDATE `value` = VALUES(`value`), `updated_at` = VALUES(`updated_at`);

-- Move the legacy ListVehicles Properties marker to the new property-list
-- attribute before removing it from Fields.
INSERT INTO `fivenet_rbac_roles_attrs`
  (`role_id`, `created_at`, `updated_at`, `attr_id`, `value`)
SELECT old_value.`role_id`, old_value.`created_at`, old_value.`updated_at`,
       property_attr.`id`, @vehicle_system_properties
FROM `fivenet_rbac_roles_attrs` old_value
INNER JOIN `fivenet_rbac_attrs` old_attr ON old_attr.`id` = old_value.`attr_id`
INNER JOIN `fivenet_rbac_permissions` perm ON perm.`id` = old_attr.`permission_id`
INNER JOIN `fivenet_rbac_attrs` property_attr
  ON property_attr.`permission_id` = perm.`id` AND property_attr.`key` = 'Properties'
WHERE old_attr.`key` = 'Fields'
  AND perm.`namespace` = 'vehicles'
  AND perm.`service` = 'VehiclesService'
  AND perm.`name` = 'ListVehicles'
  AND JSON_CONTAINS(old_value.`value`, '"Properties"', '$.stringList.strings')
ON DUPLICATE KEY UPDATE
  `value` = VALUES(`value`), `updated_at` = VALUES(`updated_at`);

UPDATE `fivenet_rbac_roles_attrs` role_attr
INNER JOIN `fivenet_rbac_attrs` attr ON attr.`id` = role_attr.`attr_id`
INNER JOIN `fivenet_rbac_permissions` perm ON perm.`id` = attr.`permission_id`
SET role_attr.`value` = REPLACE(
  REPLACE(role_attr.`value`, ',"Properties"', ''),
  '"Properties",', ''
)
WHERE attr.`key` = 'Fields'
  AND perm.`namespace` = 'vehicles'
  AND perm.`service` = 'VehiclesService'
  AND perm.`name` = 'ListVehicles';

-- SetVehicleProps has no generic property endpoint; remove only the old
-- generic marker while retaining its Wanted field grant.
UPDATE `fivenet_rbac_roles_attrs` role_attr
INNER JOIN `fivenet_rbac_attrs` attr ON attr.`id` = role_attr.`attr_id`
INNER JOIN `fivenet_rbac_permissions` perm ON perm.`id` = attr.`permission_id`
SET role_attr.`value` = REPLACE(
  REPLACE(role_attr.`value`, ',"Properties"', ''),
  '"Properties",', ''
)
WHERE attr.`key` = 'Fields'
  AND perm.`namespace` = 'vehicles'
  AND perm.`service` = 'VehiclesService'
  AND perm.`name` = 'SetVehicleProps';

-- GetVehicleProperties was folded into ListVehicles. Preserve its grants and
-- convert its old property marker to the new property-list representation.
SET @old_vehicle_props_perm_id := (
  SELECT `id` FROM `fivenet_rbac_permissions`
  WHERE `namespace` = 'vehicles' AND `service` = 'VehiclesService' AND `name` = 'GetVehicleProperties'
  LIMIT 1
);
SET @list_vehicles_perm_id := (
  SELECT `id` FROM `fivenet_rbac_permissions`
  WHERE `namespace` = 'vehicles' AND `service` = 'VehiclesService' AND `name` = 'ListVehicles'
  LIMIT 1
);
SET @list_vehicles_properties_attr_id := (
  SELECT a.`id` FROM `fivenet_rbac_attrs` a
  WHERE a.`permission_id` = @list_vehicles_perm_id AND a.`key` = 'Properties'
  LIMIT 1
);

INSERT INTO `fivenet_rbac_roles_permissions` (`role_id`, `permission_id`, `val`)
SELECT `role_id`, @list_vehicles_perm_id, `val`
FROM `fivenet_rbac_roles_permissions`
WHERE `permission_id` = @old_vehicle_props_perm_id
  AND @old_vehicle_props_perm_id IS NOT NULL
  AND @list_vehicles_perm_id IS NOT NULL
ON DUPLICATE KEY UPDATE `val` = VALUES(`val`);

INSERT INTO `fivenet_rbac_roles_attrs` (`role_id`, `created_at`, `updated_at`, `attr_id`, `value`)
SELECT old_value.`role_id`, old_value.`created_at`, old_value.`updated_at`,
       @list_vehicles_properties_attr_id, @vehicle_system_properties
FROM `fivenet_rbac_roles_attrs` old_value
INNER JOIN `fivenet_rbac_attrs` old_attr ON old_attr.`id` = old_value.`attr_id`
WHERE old_attr.`permission_id` = @old_vehicle_props_perm_id
  AND old_attr.`key` = 'Fields'
  AND old_attr.`type` = 'StringList'
  AND JSON_CONTAINS(old_value.`value`, '"Properties"', '$.stringList.strings')
  AND @old_vehicle_props_perm_id IS NOT NULL
  AND @list_vehicles_properties_attr_id IS NOT NULL
ON DUPLICATE KEY UPDATE `value` = VALUES(`value`), `updated_at` = VALUES(`updated_at`);

INSERT INTO `fivenet_rbac_job_attrs` (`job`, `attr_id`, `max_values`)
SELECT old_job_attr.`job`, @list_vehicles_properties_attr_id, old_job_attr.`max_values`
FROM `fivenet_rbac_job_attrs` old_job_attr
INNER JOIN `fivenet_rbac_attrs` old_attr ON old_attr.`id` = old_job_attr.`attr_id`
WHERE old_attr.`permission_id` = @old_vehicle_props_perm_id
  AND old_attr.`key` = 'Fields'
  AND old_attr.`type` = 'StringList'
  AND @old_vehicle_props_perm_id IS NOT NULL
  AND @list_vehicles_properties_attr_id IS NOT NULL
ON DUPLICATE KEY UPDATE `max_values` = VALUES(`max_values`);

DELETE old_value
FROM `fivenet_rbac_roles_attrs` old_value
INNER JOIN `fivenet_rbac_attrs` old_attr ON old_attr.`id` = old_value.`attr_id`
WHERE old_attr.`permission_id` = @old_vehicle_props_perm_id;
DELETE old_value
FROM `fivenet_rbac_job_attrs` old_value
INNER JOIN `fivenet_rbac_attrs` old_attr ON old_attr.`id` = old_value.`attr_id`
WHERE old_attr.`permission_id` = @old_vehicle_props_perm_id;
DELETE FROM `fivenet_rbac_attrs` WHERE `permission_id` = @old_vehicle_props_perm_id;
DELETE FROM `fivenet_rbac_roles_permissions` WHERE `permission_id` = @old_vehicle_props_perm_id;
DELETE FROM `fivenet_rbac_job_permissions` WHERE `permission_id` = @old_vehicle_props_perm_id;
DELETE FROM `fivenet_rbac_permissions` WHERE `id` = @old_vehicle_props_perm_id;

-- Fields is no longer part of either vehicle permission.
DELETE role_attr
FROM `fivenet_rbac_roles_attrs` role_attr
INNER JOIN `fivenet_rbac_attrs` attr ON attr.`id` = role_attr.`attr_id`
INNER JOIN `fivenet_rbac_permissions` perm ON perm.`id` = attr.`permission_id`
WHERE perm.`namespace` = 'vehicles' AND perm.`service` = 'VehiclesService'
  AND perm.`name` IN ('ListVehicles', 'SetVehicleProps') AND attr.`key` = 'Fields';
DELETE job_attr
FROM `fivenet_rbac_job_attrs` job_attr
INNER JOIN `fivenet_rbac_attrs` attr ON attr.`id` = job_attr.`attr_id`
INNER JOIN `fivenet_rbac_permissions` perm ON perm.`id` = attr.`permission_id`
WHERE perm.`namespace` = 'vehicles' AND perm.`service` = 'VehiclesService'
  AND perm.`name` IN ('ListVehicles', 'SetVehicleProps') AND attr.`key` = 'Fields';
DELETE attr
FROM `fivenet_rbac_attrs` attr
INNER JOIN `fivenet_rbac_permissions` perm ON perm.`id` = attr.`permission_id`
WHERE perm.`namespace` = 'vehicles' AND perm.`service` = 'VehiclesService'
  AND perm.`name` IN ('ListVehicles', 'SetVehicleProps') AND attr.`key` = 'Fields';

-- Seed the default role with every valid remaining string-list field.
INSERT INTO `fivenet_rbac_roles_attrs` (`role_id`, `created_at`, `updated_at`, `attr_id`, `value`)
SELECT default_role.`id`, CURRENT_TIMESTAMP(3), NULL, attr.`id`, attr.`valid_values`
FROM `fivenet_rbac_roles` default_role
INNER JOIN `fivenet_rbac_attrs` attr
  ON attr.`type` = 'StringList'
 AND attr.`key` = 'Fields'
WHERE default_role.`job` = '__default__'
  AND attr.`valid_values` IS NOT NULL
  AND attr.`valid_values` <> '{}'
ON DUPLICATE KEY UPDATE `value` = VALUES(`value`), `updated_at` = VALUES(`updated_at`);

COMMIT;
