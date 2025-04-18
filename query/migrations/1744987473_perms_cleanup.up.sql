BEGIN;

DELETE FROM `fivenet_permissions` WHERE `category` = 'QualificationsService' AND `name` IN (
  'CreateOrUpdateQualificationReq',
  'DeleteQualificationReq',
  'CreateOrUpdateQualificationResult',
  'DeleteQualificationResult'
) LIMIT 4;

UPDATE `fivenet_permissions` SET `name` = 'CreateCalendar', `guard_name` = 'calendarservice-createcalendar' WHERE `category` = 'CalendarService' AND `name` = 'CreateOrUpdateCalendar' LIMIT 1;
DELETE FROM `fivenet_permissions` WHERE `category` = 'CalendarService' AND `name` = 'UpdateCalendar' LIMIT 1;

DELETE FROM `fivenet_permissions` WHERE `category` = 'DocStoreService' AND `name` = 'PostComment' LIMIT 1;

DELETE FROM `fivenet_permissions` WHERE `category` = 'WikiService' AND `name` = 'UpdatePage' LIMIT 1;

COMMIT;
