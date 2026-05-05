BEGIN;

ALTER TABLE `fivenet_rbac_permissions` ADD `namespace` varchar(48) NULL AFTER `created_at`;
ALTER TABLE `fivenet_rbac_permissions` CHANGE `category` `service` varchar(64) NOT NULL;
ALTER TABLE `fivenet_rbac_permissions` MODIFY COLUMN `service` varchar(64) NOT NULL;
ALTER TABLE `fivenet_rbac_permissions` MODIFY COLUMN `name` varchar(128) NOT NULL;

ALTER TABLE `fivenet_rbac_permissions` DROP KEY `idx_category_name_unique`;
ALTER TABLE `fivenet_rbac_permissions` ADD UNIQUE KEY `idx_namespace_category_name_unique` (`namespace`, `service`, `name`);

UPDATE `fivenet_rbac_permissions`
SET
	`namespace` = SUBSTRING_INDEX(`service`, '.', 1),
	`service` = SUBSTRING_INDEX(`service`, '.', -1)
WHERE `service` LIKE '%.%';

UPDATE fivenet_rbac_permissions SET `service` = 'DispatchesService', `guard_name` = 'centrum-dispatchesservice-createdispatch' WHERE `namespace` = 'centrum' AND `service` = 'CentrumService' AND `name` = 'CreateDispatch' LIMIT 1;
UPDATE fivenet_rbac_permissions SET `service` = 'DispatchesService', `guard_name` = 'centrum-dispatchesservice-updatedispatch' WHERE `namespace` = 'centrum' AND `service` = 'CentrumService' AND `name` = 'UpdateDispatch' LIMIT 1;
UPDATE fivenet_rbac_permissions SET `service` = 'DispatchesService', `guard_name` = 'centrum-dispatchesservice-takedispatch' WHERE `namespace` = 'centrum' AND `service` = 'CentrumService' AND `name` = 'TakeDispatch' LIMIT 1;
UPDATE fivenet_rbac_permissions SET `service` = 'DispatchesService', `guard_name` = 'centrum-dispatchesservice-deletedispatch' WHERE `namespace` = 'centrum' AND `service` = 'CentrumService' AND `name` = 'DeleteDispatch' LIMIT 1;

UPDATE fivenet_rbac_permissions SET `service` = 'UnitsService', `guard_name` = 'centrum-unitsservice-createorupdateunit' WHERE `namespace` = 'centrum' AND `service` = 'CentrumService' AND `name` = 'CreateOrUpdateUnit' LIMIT 1;
UPDATE fivenet_rbac_permissions SET `service` = 'UnitsService', `guard_name` = 'centrum-unitsservice-deleteunit' WHERE `namespace` = 'centrum' AND `service` = 'CentrumService' AND `name` = 'DeleteUnit' LIMIT 1;

UPDATE fivenet_rbac_permissions SET `service` = 'ColleaguesService', `guard_name` = 'jobs-colleaguesservice-listcolleagues' WHERE `namespace` = 'jobs' AND `service` = 'JobsService' AND `name` = 'ListColleagues' LIMIT 1;
UPDATE fivenet_rbac_permissions SET `service` = 'ColleaguesService', `guard_name` = 'jobs-colleaguesservice-getcolleague' WHERE `namespace` = 'jobs' AND `service` = 'JobsService' AND `name` = 'GetColleague' LIMIT 1;
UPDATE fivenet_rbac_permissions SET `service` = 'ColleaguesService', `guard_name` = 'jobs-colleaguesservice-listcolleagueactivity' WHERE `namespace` = 'jobs' AND `service` = 'JobsService' AND `name` = 'ListColleagueActivity' LIMIT 1;
UPDATE fivenet_rbac_permissions SET `service` = 'ColleaguesService', `guard_name` = 'jobs-colleaguesservice-setcolleagueprops' WHERE `namespace` = 'jobs' AND `service` = 'JobsService' AND `name` = 'SetColleagueProps' LIMIT 1;
UPDATE fivenet_rbac_permissions SET `service` = 'ColleaguesService', `guard_name` = 'jobs-colleaguesservice-managelabels' WHERE `namespace` = 'jobs' AND `service` = 'JobsService' AND `name` = 'ManageLabels' LIMIT 1;

UPDATE fivenet_rbac_permissions SET `service` = 'TemplatesService', `guard_name` = 'documents-templatesservice-listtemplates' WHERE `namespace` = 'documents' AND `service` = 'DocumentsService' AND `name` = 'ListTemplates' LIMIT 1;
UPDATE fivenet_rbac_permissions SET `service` = 'TemplatesService', `guard_name` = 'documents-templatesservice-createtemplate' WHERE `namespace` = 'documents' AND `service` = 'DocumentsService' AND `name` = 'CreateTemplate' LIMIT 1;
UPDATE fivenet_rbac_permissions SET `service` = 'TemplatesService', `guard_name` = 'documents-templatesservice-deletetemplate' WHERE `namespace` = 'documents' AND `service` = 'DocumentsService' AND `name` = 'DeleteTemplate' LIMIT 1;

UPDATE fivenet_rbac_permissions SET `service` = 'CategoriesService', `guard_name` = 'documents-categoriesservice-listcategories' WHERE `namespace` = 'documents' AND `service` = 'DocumentsService' AND `name` = 'ListCategories' LIMIT 1;
UPDATE fivenet_rbac_permissions SET `service` = 'CategoriesService', `guard_name` = 'documents-categoriesservice-createorupdatecategory' WHERE `namespace` = 'documents' AND `service` = 'DocumentsService' AND `name` = 'CreateOrUpdateCategory' LIMIT 1;
UPDATE fivenet_rbac_permissions SET `service` = 'CategoriesService', `guard_name` = 'documents-categoriesservice-deletecategory' WHERE `namespace` = 'documents' AND `service` = 'DocumentsService' AND `name` = 'DeleteCategory' LIMIT 1;

UPDATE fivenet_rbac_permissions SET `service` = 'CommentsService', `guard_name` = 'documents-commentsservice-deletecomment' WHERE `namespace` = 'documents' AND `service` = 'DocumentsService' AND `name` = 'DeleteComment' LIMIT 1;


COMMIT;
