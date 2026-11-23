BEGIN;

-- 1792583210 converted role values from the legacy Fields attribute, but
-- copied the corresponding job max_values verbatim. Job limits therefore
-- retained the StringList wire shape and legacy field names while their
-- attributes became PropertyList. Convert those ceilings to property keys.
UPDATE `fivenet_rbac_job_attrs` job_attr
INNER JOIN (
  SELECT old_job_attr.`job`, old_job_attr.`attr_id`, CONCAT(
    '{"propertyList":{"keys":[',
    COALESCE(GROUP_CONCAT(JSON_QUOTE(CASE legacy_value.`value`
      WHEN 'Wanted' THEN 'user.wanted'
      WHEN 'UserProps.Wanted' THEN 'user.wanted'
      WHEN 'PhoneNumber' THEN 'user.phone_number'
      WHEN 'Licenses' THEN 'user.licenses'
      WHEN 'BloodType' THEN 'user.blood_type'
      WHEN 'UserProps.BloodType' THEN 'user.blood_type'
      WHEN 'Email' THEN 'user.email'
      WHEN 'UserProps.Email' THEN 'user.email'
      WHEN 'Job' THEN 'user.job'
      WHEN 'UserProps.Job' THEN 'user.job'
      WHEN 'Labels' THEN 'user.labels'
      WHEN 'UserProps.Labels' THEN 'user.labels'
      WHEN 'MugShot' THEN 'user.mugshot'
      WHEN 'Mugshot' THEN 'user.mugshot'
      WHEN 'UserProps.MugShot' THEN 'user.mugshot'
      WHEN 'UserProps.Mugshot' THEN 'user.mugshot'
      WHEN 'OpenFines' THEN 'user.open_fines'
      WHEN 'UserProps.OpenFines' THEN 'user.open_fines'
      WHEN 'TrafficInfractionPoints' THEN 'user.traffic_infraction_points'
      WHEN 'UserProps.TrafficInfractionPoints' THEN 'user.traffic_infraction_points'
      ELSE legacy_value.`value`
    END) ORDER BY legacy_value.`value` SEPARATOR ','), ''),
    ']}}'
  ) AS `max_values`
  FROM `fivenet_rbac_job_attrs` old_job_attr
  INNER JOIN `fivenet_rbac_attrs` attr ON attr.`id` = old_job_attr.`attr_id`
  INNER JOIN `fivenet_rbac_permissions` perm ON perm.`id` = attr.`permission_id`
  LEFT JOIN JSON_TABLE(
    old_job_attr.`max_values`,
    '$.stringList.strings[*]' COLUMNS (`value` varchar(128) PATH '$')
  ) AS legacy_value ON TRUE
  WHERE attr.`key` = 'Properties'
    AND attr.`type` = 'PropertyList'
    AND perm.`namespace` = 'citizens'
    AND perm.`service` = 'CitizensService'
    AND perm.`name` IN ('ListCitizens', 'SetUserProps')
    AND JSON_TYPE(JSON_EXTRACT(old_job_attr.`max_values`, '$.stringList')) = 'OBJECT'
  GROUP BY old_job_attr.`job`, old_job_attr.`attr_id`
) converted
  ON converted.`job` = job_attr.`job`
 AND converted.`attr_id` = job_attr.`attr_id`
SET job_attr.`max_values` = converted.`max_values`;

COMMIT;
