BEGIN;

-- Reverse the citizen label access shift from VIEW/GIVE/REMOVE = 2/3/4 to 1/2/3.
-- This assumes no ACCESS_LEVEL_BLOCKED rows were created after the up migration.
UPDATE `fivenet_user_labels_job_job_access`
SET `access` = `access` + 10
WHERE `access` BETWEEN 2 AND 4;

UPDATE `fivenet_user_labels_job_job_access`
SET `access` = `access` - 11
WHERE `access` BETWEEN 12 AND 14;

UPDATE `fivenet_user_labels_job_visibility_subject`
SET `access` = `access` + 10
WHERE `access` BETWEEN 2 AND 4;

UPDATE `fivenet_user_labels_job_visibility_subject`
SET `access` = `access` - 11
WHERE `access` BETWEEN 12 AND 14;

COMMIT;
