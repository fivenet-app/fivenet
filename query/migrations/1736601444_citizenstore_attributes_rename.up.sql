BEGIN;

UPDATE `fivenet_permissions`
SET
  `name` = 'ManageCitizenLabels',
  `guard_name` = 'citizenstoreservice-managecitizenlabels'
WHERE
  `name` = 'ManageCitizenAttributes';

UPDATE `fivenet_permissions`
SET
  `name` = 'CompleteCitizenLabels',
  `guard_name` = 'completorservice-completecitizenlabels'
WHERE
  `name` = 'CompleteCitizenAttributes';

UPDATE
	fivenet_attrs fa
SET
	fa.valid_values = REPLACE(fa.valid_values, 'Attributes', 'Labels')
WHERE
	fa.permission_id IN (
	SELECT
		id
	FROM
		fivenet_permissions fp
	WHERE
		(fp.category = 'CitizenStoreService'
			AND fp.name = 'SetUserProps')
		OR
		(fp.category = 'CitizenStoreService'
			AND fp.name = 'ListCitizens'));

UPDATE
	fivenet_role_attrs fra
SET
	fra.value = REPLACE(fra.value, 'Attributes', 'Labels')
WHERE
	fra.attr_id IN (
	SELECT
		id
	FROM
		fivenet_attrs fa
	WHERE
		fa.permission_id IN (
		SELECT
			id
		FROM
			fivenet_permissions fp
		WHERE
			(fp.category = 'CitizenStoreService'
				AND fp.name = 'SetUserProps')
			OR
		(fp.category = 'CitizenStoreService'
				AND fp.name = 'ListCitizens')
			)
		);

RENAME TABLE `fivenet_user_citizen_attributes` TO `fivenet_user_citizen_labels`;
RENAME TABLE `fivenet_job_citizen_attributes` TO `fivenet_job_citizen_labels`;

ALTER TABLE `fivenet_job_props` DROP COLUMN `citizen_attributes`;

COMMIT;
