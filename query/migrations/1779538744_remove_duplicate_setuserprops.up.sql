BEGIN;

DELETE FROM `fivenet_rbac_permissions` WHERE `namespace` = 'citizens' AND `service` = 'LabelsService' AND `name` = 'SetUserProps' LIMIT 1;

COMMIT;
