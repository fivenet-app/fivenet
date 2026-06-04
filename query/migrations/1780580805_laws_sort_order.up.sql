BEGIN;

-- Table: `fivenet_lawbooks`
ALTER TABLE `fivenet_lawbooks`
  ADD COLUMN `sort_order` int(11) NOT NULL DEFAULT 0 AFTER `name`;

ALTER TABLE `fivenet_lawbooks`
  ADD KEY `idx_fivenet_lawbooks_sort_order` (`sort_order`, `id`);

UPDATE `fivenet_lawbooks` u
JOIN (
  SELECT
    `id`,
    ROW_NUMBER() OVER (ORDER BY `sort_key` ASC, `id` ASC) - 1 AS `new_sort_order`
  FROM `fivenet_lawbooks`
) x ON x.id = u.id
SET u.sort_order = x.new_sort_order;

-- Table: `fivenet_lawbooks_laws`
ALTER TABLE `fivenet_lawbooks_laws`
  ADD COLUMN `sort_order` int(11) NOT NULL DEFAULT 0 AFTER `name`;

ALTER TABLE `fivenet_lawbooks_laws`
  ADD KEY `idx_fivenet_lawbooks_laws_sort_order` (`sort_order`, `id`);

UPDATE `fivenet_lawbooks_laws` u
JOIN (
  SELECT
    `id`,
    ROW_NUMBER() OVER (PARTITION BY `lawbook_id` ORDER BY `sort_key` ASC, `id` ASC) - 1 AS `new_sort_order`
  FROM `fivenet_lawbooks_laws`
) x ON x.id = u.id
SET u.sort_order = x.new_sort_order;

COMMIT;
