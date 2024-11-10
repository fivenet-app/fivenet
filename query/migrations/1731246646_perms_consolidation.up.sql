BEGIN;

DELETE FROM fivenet_permissions WHERE category = 'DocStoreService' AND name = 'GetDocument' LIMIT 1;
DELETE FROM fivenet_permissions WHERE category = 'QualificationsService' AND name = 'GetQualification' LIMIT 1;
DELETE FROM fivenet_permissions WHERE category = 'WikiService' AND name = 'GetPage' LIMIT 1;

COMMIT;
