BEGIN;

ALTER TABLE `fivenet_wiki_pages`
  ADD COLUMN `sort_rank` varchar(32) NOT NULL DEFAULT '' AFTER `parent_id`;

ALTER TABLE `fivenet_wiki_pages`
  ADD KEY `idx_fivenet_wiki_pages_job_parent_sort_rank` (`job`, `parent_id`, `sort_rank`, `id`);

ALTER TABLE `fivenet_wiki_pages`
  ADD KEY `idx_fivenet_wiki_pages_job_startpage_parent_sort_rank` (`job`, `startpage`, `parent_id`, `sort_rank`, `id`);

UPDATE `fivenet_wiki_pages` p
JOIN (
  SELECT
    `id`,
    LPAD(
      CAST(
        ROW_NUMBER() OVER (
          PARTITION BY `parent_id`
          ORDER BY `startpage` DESC, `sort_key` ASC, `draft` ASC, `id` ASC
        ) * 1000 AS CHAR
      ),
      12,
      '0'
    ) AS `new_sort_rank`
  FROM `fivenet_wiki_pages`
) x ON x.id = p.id
SET p.sort_rank = x.new_sort_rank;

COMMIT;
