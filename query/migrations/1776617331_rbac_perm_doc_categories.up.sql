BEGIN;

SET @source_perm_id := (
    SELECT `id`
    FROM `fivenet_rbac_permissions`
    WHERE `category` = 'completor.CompletorService' AND `name` = 'CompleteDocumentCategories'
    LIMIT 1
);
SET @target_perm_id := (
    SELECT `id`
    FROM `fivenet_rbac_permissions`
    WHERE `category` = 'documents.DocumentsService' AND `name` = 'ListCategories'
    LIMIT 1
);

-- Move role and job permission assignments from CompleteDocumentCategories to ListCategories
INSERT INTO `fivenet_rbac_roles_permissions` (`role_id`, `permission_id`, `val`)
SELECT
    frp.`role_id`,
    @target_perm_id,
    frp.`val`
FROM `fivenet_rbac_roles_permissions` frp
WHERE frp.`permission_id` = @source_perm_id
    AND @source_perm_id IS NOT NULL
    AND @target_perm_id IS NOT NULL
ON DUPLICATE KEY UPDATE `val` = VALUES(`val`);

INSERT INTO `fivenet_rbac_job_permissions` (`job`, `permission_id`, `val`)
SELECT
    fjp.`job`,
    @target_perm_id,
    fjp.`val`
FROM `fivenet_rbac_job_permissions` fjp
WHERE fjp.`permission_id` = @source_perm_id
    AND @source_perm_id IS NOT NULL
    AND @target_perm_id IS NOT NULL
ON DUPLICATE KEY UPDATE `val` = VALUES(`val`);

-- Move the Job(s) attribute from CompleteDocumentCategories to ListCategories.
-- If a Jobs/Job attribute already exists on ListCategories, merge values into it and remove the old one.
SET @source_attr_id := (
    SELECT `id`
    FROM `fivenet_rbac_attrs`
    WHERE `permission_id` = @source_perm_id
        AND `key` IN ('Jobs', 'Job')
    LIMIT 1
);
SET @target_attr_id := (
    SELECT `id`
    FROM `fivenet_rbac_attrs`
    WHERE `permission_id` = @target_perm_id
        AND `key` IN ('Jobs', 'Job')
    LIMIT 1
);

UPDATE `fivenet_rbac_attrs`
SET `permission_id` = @target_perm_id
WHERE `id` = @source_attr_id
    AND @source_attr_id IS NOT NULL
    AND @target_perm_id IS NOT NULL
    AND @target_attr_id IS NULL;

INSERT INTO `fivenet_rbac_job_attrs` (`job`, `attr_id`, `max_values`)
SELECT
    fja.`job`,
    @target_attr_id,
    fja.`max_values`
FROM `fivenet_rbac_job_attrs` fja
WHERE fja.`attr_id` = @source_attr_id
    AND @source_attr_id IS NOT NULL
    AND @target_attr_id IS NOT NULL
ON DUPLICATE KEY UPDATE `max_values` = VALUES(`max_values`);

INSERT INTO `fivenet_rbac_roles_attrs` (`role_id`, `created_at`, `updated_at`, `attr_id`, `value`)
SELECT
    fra.`role_id`,
    fra.`created_at`,
    fra.`updated_at`,
    @target_attr_id,
    fra.`value`
FROM `fivenet_rbac_roles_attrs` fra
WHERE fra.`attr_id` = @source_attr_id
    AND @source_attr_id IS NOT NULL
    AND @target_attr_id IS NOT NULL
ON DUPLICATE KEY UPDATE `value` = VALUES(`value`), `updated_at` = VALUES(`updated_at`);

DELETE FROM `fivenet_rbac_job_attrs`
WHERE `attr_id` = @source_attr_id
    AND @source_attr_id IS NOT NULL
    AND @target_attr_id IS NOT NULL;

DELETE FROM `fivenet_rbac_roles_attrs`
WHERE `attr_id` = @source_attr_id
    AND @source_attr_id IS NOT NULL
    AND @target_attr_id IS NOT NULL;

DELETE FROM `fivenet_rbac_attrs`
WHERE `id` = @source_attr_id
    AND @source_attr_id IS NOT NULL
    AND @target_attr_id IS NOT NULL;

-- Drop old CompleteDocumentCategories permission once all grants/attrs have been moved.
DELETE FROM `fivenet_rbac_permissions`
WHERE `id` = @source_perm_id
    AND @source_perm_id IS NOT NULL;

COMMIT;
