BEGIN;

DELETE t
FROM `fivenet_calendar_access` t
INNER JOIN (
  SELECT `target_id`, `subject_id`, MAX(`access`) AS `access`
  FROM `fivenet_calendar_access`
  WHERE `effect` = 0
  GROUP BY `target_id`, `subject_id`
  HAVING COUNT(*) > 1
) keep_rows
  ON keep_rows.`target_id` = t.`target_id`
  AND keep_rows.`subject_id` = t.`subject_id`
WHERE t.`effect` = 0
  AND t.`access` < keep_rows.`access`;

DELETE t
FROM `fivenet_centrum_units_access` t
INNER JOIN (
  SELECT `target_id`, `subject_id`, MAX(`access`) AS `access`
  FROM `fivenet_centrum_units_access`
  WHERE `effect` = 0
  GROUP BY `target_id`, `subject_id`
  HAVING COUNT(*) > 1
) keep_rows
  ON keep_rows.`target_id` = t.`target_id`
  AND keep_rows.`subject_id` = t.`subject_id`
WHERE t.`effect` = 0
  AND t.`access` < keep_rows.`access`;

DELETE t
FROM `fivenet_documents_access` t
INNER JOIN (
  SELECT `target_id`, `subject_id`, MAX(`access`) AS `access`
  FROM `fivenet_documents_access`
  WHERE `effect` = 0
  GROUP BY `target_id`, `subject_id`
  HAVING COUNT(*) > 1
) keep_rows
  ON keep_rows.`target_id` = t.`target_id`
  AND keep_rows.`subject_id` = t.`subject_id`
WHERE t.`effect` = 0
  AND t.`access` < keep_rows.`access`;

DELETE t
FROM `fivenet_documents_stamps_access` t
INNER JOIN (
  SELECT `target_id`, `subject_id`, MAX(`access`) AS `access`
  FROM `fivenet_documents_stamps_access`
  WHERE `effect` = 0
  GROUP BY `target_id`, `subject_id`
  HAVING COUNT(*) > 1
) keep_rows
  ON keep_rows.`target_id` = t.`target_id`
  AND keep_rows.`subject_id` = t.`subject_id`
WHERE t.`effect` = 0
  AND t.`access` < keep_rows.`access`;

DELETE t
FROM `fivenet_documents_templates_access` t
INNER JOIN (
  SELECT `target_id`, `subject_id`, MAX(`access`) AS `access`
  FROM `fivenet_documents_templates_access`
  WHERE `effect` = 0
  GROUP BY `target_id`, `subject_id`
  HAVING COUNT(*) > 1
) keep_rows
  ON keep_rows.`target_id` = t.`target_id`
  AND keep_rows.`subject_id` = t.`subject_id`
WHERE t.`effect` = 0
  AND t.`access` < keep_rows.`access`;

DELETE t
FROM `fivenet_mailer_emails_access` t
INNER JOIN (
  SELECT `target_id`, `subject_id`, MAX(`access`) AS `access`
  FROM `fivenet_mailer_emails_access`
  WHERE `effect` = 0
  GROUP BY `target_id`, `subject_id`
  HAVING COUNT(*) > 1
) keep_rows
  ON keep_rows.`target_id` = t.`target_id`
  AND keep_rows.`subject_id` = t.`subject_id`
WHERE t.`effect` = 0
  AND t.`access` < keep_rows.`access`;

DELETE t
FROM `fivenet_qualifications_access` t
INNER JOIN (
  SELECT `target_id`, `subject_id`, MAX(`access`) AS `access`
  FROM `fivenet_qualifications_access`
  WHERE `effect` = 0
  GROUP BY `target_id`, `subject_id`
  HAVING COUNT(*) > 1
) keep_rows
  ON keep_rows.`target_id` = t.`target_id`
  AND keep_rows.`subject_id` = t.`subject_id`
WHERE t.`effect` = 0
  AND t.`access` < keep_rows.`access`;

DELETE t
FROM `fivenet_user_labels_job_job_access` t
INNER JOIN (
  SELECT `target_id`, `subject_id`, MAX(`access`) AS `access`
  FROM `fivenet_user_labels_job_job_access`
  WHERE `effect` = 0
  GROUP BY `target_id`, `subject_id`
  HAVING COUNT(*) > 1
) keep_rows
  ON keep_rows.`target_id` = t.`target_id`
  AND keep_rows.`subject_id` = t.`subject_id`
WHERE t.`effect` = 0
  AND t.`access` < keep_rows.`access`;

DELETE t
FROM `fivenet_wiki_pages_access` t
INNER JOIN (
  SELECT `target_id`, `subject_id`, MAX(`access`) AS `access`
  FROM `fivenet_wiki_pages_access`
  WHERE `effect` = 0
  GROUP BY `target_id`, `subject_id`
  HAVING COUNT(*) > 1
) keep_rows
  ON keep_rows.`target_id` = t.`target_id`
  AND keep_rows.`subject_id` = t.`subject_id`
WHERE t.`effect` = 0
  AND t.`access` < keep_rows.`access`;

COMMIT;
