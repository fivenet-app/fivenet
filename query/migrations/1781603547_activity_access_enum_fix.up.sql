BEGIN;

UPDATE
	`fivenet_documents_activity`
SET
	`data` = REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(`data`, '"access":"ACCESS_LEVEL_UNSPECIFIED"', '"access":0'), '"access": "ACCESS_LEVEL_UNSPECIFIED"', '"access":0'), '"access":"ACCESS_LEVEL_BLOCKED"', '"access":1'), '"access": "ACCESS_LEVEL_BLOCKED"', '"access":1'), '"access":"ACCESS_LEVEL_VIEW"', '"access":2'), '"access": "ACCESS_LEVEL_VIEW"', '"access":2'), '"access":"ACCESS_LEVEL_COMMENT"', '"access":3'), '"access": "ACCESS_LEVEL_COMMENT"', '"access":3'), '"access":"ACCESS_LEVEL_STATUS"', '"access":4'), '"access": "ACCESS_LEVEL_STATUS"', '"access":4'), '"access":"ACCESS_LEVEL_ACCESS"', '"access":5'), '"access": "ACCESS_LEVEL_ACCESS"', '"access":5'), '"access":"ACCESS_LEVEL_EDIT"', '"access":6'), '"access": "ACCESS_LEVEL_EDIT"', '"access":6')
WHERE
	`data` LIKE '%ACCESS_LEVEL_%';

UPDATE
	`fivenet_wiki_pages_activity`
SET
	`data` = REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(`data`, '"access":"ACCESS_LEVEL_UNSPECIFIED"', '"access":0'), '"access": "ACCESS_LEVEL_UNSPECIFIED"', '"access":0'), '"access":"ACCESS_LEVEL_BLOCKED"', '"access":1'), '"access": "ACCESS_LEVEL_BLOCKED"', '"access":1'), '"access":"ACCESS_LEVEL_VIEW"', '"access":2'), '"access": "ACCESS_LEVEL_VIEW"', '"access":2'), '"access":"ACCESS_LEVEL_ACCESS"', '"access":3'), '"access": "ACCESS_LEVEL_ACCESS"', '"access":3'), '"access":"ACCESS_LEVEL_EDIT"', '"access":4'), '"access": "ACCESS_LEVEL_EDIT"', '"access":4')
WHERE
	`data` LIKE '%ACCESS_LEVEL_%';

COMMIT;
