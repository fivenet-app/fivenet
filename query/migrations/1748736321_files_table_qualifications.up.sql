BEGIN;

DROP TABLE IF EXISTS `fivenet_qualifications_exam_questions_files`;

-- Table: `fivenet_qualifications` - Add draft field
ALTER TABLE `fivenet_qualifications` ADD COLUMN `draft` tinyint(1) DEFAULT '0';
ALTER TABLE `fivenet_qualifications` CHANGE `draft` `draft` tinyint(1) DEFAULT '0' AFTER `closed`;

ALTER TABLE `fivenet_qualifications` ADD INDEX `idx_draft` (`draft`);

ALTER TABLE `fivenet_qualifications` ADD COLUMN `content_type` smallint(2) NOT NULL AFTER `description`;

UPDATE `fivenet_qualifications` SET `content_type` = 1 WHERE `content_type` = 0;

-- Perm renames
-- Rename `CreatePage` to `UpdatePage` in `wiki.WikiService`
UPDATE `fivenet_rbac_permissions` SET `name` = 'UpdatePage', `guard_name` = 'wiki-wikiservice-updatepage' WHERE `category` = 'wiki.WikiService' AND `name` = 'CreatePage' LIMIT 1;

-- Rename `CreateDocument` to `UpdateDocument` in `documents.DocumentsService`
-- 1) Look up perm IDs
SET @create_perm_id = (
  SELECT id
  FROM fivenet_rbac_permissions
  WHERE name = 'CreateDocument'
);

SET @update_perm_id = (
  SELECT id
  FROM fivenet_rbac_permissions
  WHERE name = 'UpdateDocument'
);

-- 2) Grant `UpdateDocument` to every role that has `CreateDocument`, skipping duplicates:
INSERT INTO fivenet_rbac_roles_permissions (role_id, permission_id, val)
SELECT
  rp.role_id,
  @update_perm_id,
  1
FROM
  fivenet_rbac_roles_permissions AS rp
WHERE
  rp.permission_id = @create_perm_id
ON DUPLICATE KEY UPDATE
  -- do nothing; this prevents duplicate key errors if (role_id,@update_perm_id) already exists
  role_id = rp.role_id;

-- 3) Remove `CreateDocument` entirely:
DELETE rp
FROM fivenet_rbac_roles_permissions AS rp
WHERE rp.permission_id = @create_perm_id;

-- Rename `qualifications.QualificationsService` `CreateQualification` -> `UpdateQualification`
UPDATE
	fivenet_rbac_attrs
SET
	permission_id = (
	SELECT
		id
	FROM
		fivenet_rbac_permissions
	WHERE
		category = 'qualifications.QualificationsService'
		AND name = 'UpdateQualification'
	)
WHERE
	permission_id = (SELECT
		id
	FROM
		fivenet_rbac_permissions
	WHERE
		category = 'qualifications.QualificationsService'
		AND name = 'CreateQualification');

-- 1) Look up perm IDs
SET @create_perm_id = (
  SELECT id
  FROM fivenet_rbac_permissions
  WHERE name = 'CreateQualification'
);

SET @update_perm_id = (
  SELECT id
  FROM fivenet_rbac_permissions
  WHERE name = 'UpdateQualification'
);

SELECT @create_perm_id, @update_perm_id;

-- 2) Grant `UpdateQualification` to every role that has `CreateQualification`, skipping duplicates:
INSERT INTO fivenet_rbac_roles_permissions (role_id, permission_id, val)
SELECT
  rp.role_id,
  @update_perm_id,
  1
FROM
  fivenet_rbac_roles_permissions AS rp
WHERE
  rp.permission_id = @create_perm_id
ON DUPLICATE KEY UPDATE
  -- do nothing; this prevents duplicate key errors if (role_id,@update_perm_id) already exists
  role_id = rp.role_id;

-- 3) Remove `CreateQualification` entirely:
DELETE rp
FROM fivenet_rbac_roles_permissions AS rp
WHERE rp.permission_id = @create_perm_id;

COMMIT;
