BEGIN;

ALTER TABLE `fivenet_documents_templates`
  ADD COLUMN `weight` int(11) unsigned DEFAULT 0 AFTER `deleted_at`;

UPDATE `fivenet_documents_templates` t
JOIN (
  SELECT
    `id`,
    (COUNT(*) OVER (PARTITION BY `creator_job`) - ROW_NUMBER() OVER (
      PARTITION BY `creator_job`
      ORDER BY `sort_rank` ASC, `id` ASC
    ) + 1) * 1000 AS `new_weight`
  FROM `fivenet_documents_templates`
) x ON x.id = t.id
SET t.weight = x.new_weight;

ALTER TABLE `fivenet_documents_templates`
  ADD KEY `idx_fivenet_documents_templates_weight` (`weight`);

ALTER TABLE `fivenet_documents_templates`
  DROP KEY `idx_fivenet_documents_templates_creator_job_sort_rank`;

ALTER TABLE `fivenet_documents_templates`
  DROP COLUMN `sort_rank`;

COMMIT;
