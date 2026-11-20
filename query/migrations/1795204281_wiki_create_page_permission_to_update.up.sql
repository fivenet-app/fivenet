BEGIN;

SET @create_page_permission_id := (
  SELECT `id`
  FROM `fivenet_rbac_permissions`
  WHERE `namespace` = 'wiki' AND `service` = 'WikiService' AND `name` = 'CreatePage'
  LIMIT 1
);
SET @update_page_permission_id := (
  SELECT `id`
  FROM `fivenet_rbac_permissions`
  WHERE `namespace` = 'wiki' AND `service` = 'WikiService' AND `name` = 'UpdatePage'
  LIMIT 1
);

-- Preserve explicit CreatePage grants as UpdatePage grants before retiring
-- CreatePage. A true grant wins over any existing false override.
INSERT INTO `fivenet_rbac_roles_permissions` (`role_id`, `permission_id`, `val`)
SELECT `role_id`, @update_page_permission_id, TRUE
FROM `fivenet_rbac_roles_permissions`
WHERE `permission_id` = @create_page_permission_id AND `val` = TRUE
ON DUPLICATE KEY UPDATE `val` = TRUE;

INSERT INTO `fivenet_rbac_job_permissions` (`job`, `permission_id`, `val`)
SELECT `job`, @update_page_permission_id, TRUE
FROM `fivenet_rbac_job_permissions`
WHERE `permission_id` = @create_page_permission_id AND `val` = TRUE
ON DUPLICATE KEY UPDATE `val` = TRUE;

DELETE FROM `fivenet_rbac_roles_permissions` WHERE `permission_id` = @create_page_permission_id;
DELETE FROM `fivenet_rbac_job_permissions` WHERE `permission_id` = @create_page_permission_id;
DELETE FROM `fivenet_rbac_permissions` WHERE `id` = @create_page_permission_id;

COMMIT;
