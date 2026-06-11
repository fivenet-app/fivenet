BEGIN;

ALTER TABLE `fivenet_wiki_pages`
  DROP KEY `idx_fivenet_wiki_pages_job_startpage_parent_sort_rank`;

ALTER TABLE `fivenet_wiki_pages`
  DROP KEY `idx_fivenet_wiki_pages_job_parent_sort_rank`;

ALTER TABLE `fivenet_wiki_pages`
  DROP COLUMN `sort_rank`;

COMMIT;
