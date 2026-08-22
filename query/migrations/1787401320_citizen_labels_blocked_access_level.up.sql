BEGIN;

-- Citizen label access levels gained ACCESS_LEVEL_BLOCKED at 1.
-- Shift existing persisted levels from VIEW/GIVE/REMOVE = 1/2/3 to 2/3/4.
-- Use a temporary offset to avoid unique-key collisions between adjacent levels.
UPDATE `fivenet_user_labels_job_job_access`
SET `access` = `access` + 10
WHERE `access` BETWEEN 1 AND 3;

UPDATE `fivenet_user_labels_job_job_access`
SET `access` = `access` - 9
WHERE `access` BETWEEN 11 AND 13;

UPDATE `fivenet_user_labels_job_visibility_subject`
SET `access` = `access` + 10
WHERE `access` BETWEEN 1 AND 3;

UPDATE `fivenet_user_labels_job_visibility_subject`
SET `access` = `access` - 9
WHERE `access` BETWEEN 11 AND 13;

COMMIT;
