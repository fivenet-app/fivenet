BEGIN;

ALTER TABLE `fivenet_documents_templates`
  ADD COLUMN `sort_rank` varchar(32) NOT NULL DEFAULT '' AFTER `deleted_at`;

ALTER TABLE `fivenet_documents_templates`
  ADD KEY `idx_fivenet_documents_templates_creator_job_sort_rank` (`creator_job`, `sort_rank`, `id`);

UPDATE `fivenet_documents_templates` t
JOIN (
  SELECT
    `id`,
    LPAD(
      CAST(
        ROW_NUMBER() OVER (
          PARTITION BY `creator_job`
          ORDER BY `weight` DESC, `id` ASC
        ) * 1000 AS CHAR
      ),
      12,
      '0'
    ) AS `new_sort_rank`
  FROM `fivenet_documents_templates`
) x ON x.id = t.id
SET t.sort_rank = x.new_sort_rank;

ALTER TABLE `fivenet_documents_templates`
  DROP KEY `idx_fivenet_documents_templates_weight`;

ALTER TABLE `fivenet_documents_templates`
  DROP COLUMN `weight`;

COMMIT;
