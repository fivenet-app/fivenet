BEGIN;

DELETE FROM `fivenet_user_activity` WHERE `type` = 9 AND `data` LIKE '%ACCESS_LEVEL_%' LIMIT 1;

COMMIT;
