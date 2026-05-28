BEGIN;

ALTER TABLE `fivenet_user_labels_job`
  ADD COLUMN `sort_order` int(11) NOT NULL DEFAULT 0 AFTER `job`;

ALTER TABLE `fivenet_user_labels_job`
  ADD KEY `idx_fivenet_user_labels_job_job_sort_order` (`job`, `sort_order`, `id`);

UPDATE `fivenet_user_labels_job` u
JOIN (
  SELECT
    `id`,
    ROW_NUMBER() OVER (PARTITION BY `job` ORDER BY `name` ASC, `id` ASC) - 1 AS `new_sort_order`
  FROM `fivenet_user_labels_job`
) x ON x.id = u.id
SET u.sort_order = x.new_sort_order;

ALTER TABLE `fivenet_job_labels` CHANGE `order` `sort_order` mediumint(9) DEFAULT 0 NULL;
ALTER TABLE `fivenet_job_labels` MODIFY COLUMN `sort_order` int(11) DEFAULT 0 NULL;

COMMIT;
