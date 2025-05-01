BEGIN;

DELETE FROM `fivenet_permissions` WHERE `category` = 'CalendarService' AND `name` = 'CreateOrUpdateCalendarEntry' LIMIT 1;

-- Ensure the previous permissions are removed for sure
DELETE FROM `fivenet_permissions` WHERE `category` = 'DocStoreService' AND `name` = 'PostComment' LIMIT 1;

DELETE FROM `fivenet_permissions` WHERE `category` = 'WikiService' AND `name` = 'UpdatePage' LIMIT 1;

COMMIT;
