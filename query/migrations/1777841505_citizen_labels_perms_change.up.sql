BEGIN;

-- Rename `completor.CompletorService/CompleteCitizenLabels` to `citizens.LabelsService/CompleteCitizenLabels`
DELETE
FROM
	`fivenet_rbac_attrs`
WHERE
	`permission_id` = (SELECT `id` FROM `fivenet_rbac_permissions` WHERE `category` = 'completor.CompletorService' AND `name` = 'CompleteCitizenLabels' LIMIT 1)
	AND `key` = 'Jobs'
LIMIT 1;

UPDATE
	`fivenet_rbac_permissions`
SET
	`category` = 'citizens.LabelsService',
	`name` = 'ListLabels',
	`guard_name` = 'citizens-labelsservice-listlabels',
  `order` = 3200,
  `icon` = 'i-mdi-label-multiple'
WHERE
	`category` = 'completor.CompletorService'
	AND `name` = 'CompleteCitizenLabels'
LIMIT 1;

COMMIT;
