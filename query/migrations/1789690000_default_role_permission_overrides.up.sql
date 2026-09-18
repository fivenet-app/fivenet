BEGIN;

-- Default permissions are unconditional. Remove redundant or conflicting
-- permission entries from job roles so they inherit the default role.
DELETE role_permission
FROM `fivenet_rbac_roles_permissions` AS role_permission
INNER JOIN `fivenet_rbac_roles` AS job_role
  ON job_role.id = role_permission.role_id
INNER JOIN `fivenet_rbac_roles_permissions` AS default_permission
  ON default_permission.permission_id = role_permission.permission_id
INNER JOIN `fivenet_rbac_roles` AS default_role
  ON default_role.id = default_permission.role_id
WHERE default_role.job = '__default__'
  AND default_permission.val = TRUE
  AND job_role.job <> '__default__';

COMMIT;
