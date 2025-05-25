BEGIN;

UPDATE `fivenet_job_props` SET `quick_buttons` = REPLACE(`quick_buttons`, ', "bodyCheckup":true', '') WHERE `quick_buttons` LIKE '%"bodyCheckup"%';
UPDATE `fivenet_job_props` SET `quick_buttons` = REPLACE(`quick_buttons`, ',"bodyCheckup":true', '') WHERE `quick_buttons` LIKE '%"bodyCheckup"%';
UPDATE `fivenet_job_props` SET `quick_buttons` = REPLACE(`quick_buttons`, '"bodyCheckup":true, ', '') WHERE `quick_buttons` LIKE '%"bodyCheckup"%';
UPDATE `fivenet_job_props` SET `quick_buttons` = REPLACE(`quick_buttons`, '"bodyCheckup":true,', '') WHERE `quick_buttons` LIKE '%"bodyCheckup"%';
UPDATE `fivenet_job_props` SET `quick_buttons` = REPLACE(`quick_buttons`, '"bodyCheckup":true', '') WHERE `quick_buttons` LIKE '%"bodyCheckup"%';

DELETE FROM `fivenet_rbac_permissions`
  WHERE category = 'citizens.CitizensService' AND name = 'ManageLabels'
  LIMIT 1;

UPDATE `fivenet_rbac_permissions`
	SET name='ManageLabels', guard_name='citizens-citizensservice-managelabels'
	WHERE category = 'citizens.CitizensService' AND name = 'ManageCitizenLabels';

COMMIT;
