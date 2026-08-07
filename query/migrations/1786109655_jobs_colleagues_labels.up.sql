BEGIN;

UPDATE
	`fivenet_rbac_permissions`
SET
	`name` = 'CreateOrUpdateLabel',
	`guard_name` = 'jobs-colleaguesservice-createorupdatelabel'
WHERE
	`service` = 'ColleaguesService'
	AND `name` = 'ManageLabels'
	AND `guard_name` = 'jobs-colleaguesservice-managelabels'
LIMIT 1;

COMMIT;
