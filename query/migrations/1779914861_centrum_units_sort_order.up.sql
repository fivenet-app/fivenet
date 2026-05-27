BEGIN;

ALTER TABLE `fivenet_centrum_units`
  ADD COLUMN `sort_order` int(11) NOT NULL DEFAULT 0 AFTER `job`;

ALTER TABLE `fivenet_centrum_units`
  ADD KEY `idx_fivenet_centrum_units_job_sort_order` (`job`, `sort_order`, `id`);

UPDATE `fivenet_centrum_units` u
JOIN (
  SELECT
    `id`,
    ROW_NUMBER() OVER (PARTITION BY `job` ORDER BY `name` ASC, `id` ASC) - 1 AS `new_sort_order`
  FROM `fivenet_centrum_units`
) x ON x.id = u.id
SET u.sort_order = x.new_sort_order;

COMMIT;
