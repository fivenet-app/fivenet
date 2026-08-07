BEGIN;

UPDATE
	`fivenet_rbac_permissions`
SET
	`name` = 'ManageLabels',
	`guard_name` = 'jobs-colleaguesservice-managelabels'
WHERE
	`service` = 'ColleaguesService'
	AND `name` = 'CreateOrUpdateLabel'
	AND `guard_name` = 'jobs-colleaguesservice-createorupdatelabel'
LIMIT 1;

COMMIT;
