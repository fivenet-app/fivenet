BEGIN;

-- Delete perm ListLabels
DELETE FROM `fivenet_rbac_permissions`
WHERE `namespace` = 'citizens' AND `service` = 'LabelsService' AND `name` = 'ListLabels' LIMIT 1;

COMMIT;
