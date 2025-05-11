BEGIN;

UPDATE `fivenet_role_permissions` SET `permission_id` = (
	SELECT
		id
	FROM
		fivenet_permissions
	WHERE
		category = 'DocStoreService'
		AND name = 'CreateOrUpdateCategory'
	LIMIT 1
)
WHERE
	`permission_id` = (SELECT `id` FROM `fivenet_permissions` WHERE `category` = 'DocStoreService' AND name = 'CreateCategory' LIMIT 1);

UPDATE `fivenet_job_permissions` SET `permission_id` = (
	SELECT
		id
	FROM
		fivenet_permissions
	WHERE
		category = 'DocStoreService'
		AND name = 'CreateOrUpdateCategory'
	LIMIT 1
)
WHERE
	`permission_id` = (SELECT `id` FROM `fivenet_permissions` WHERE `category` = 'DocStoreService' AND name = 'CreateCategory' LIMIT 1);

DELETE FROM `fivenet_permissions` WHERE `category` = 'DocStoreService' AND `name` = 'CreateCategory' LIMIT 1;

COMMIT;
