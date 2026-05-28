BEGIN;

ALTER TABLE `fivenet_centrum_units`
  DROP KEY `idx_fivenet_centrum_units_job_sort_order`;

ALTER TABLE `fivenet_centrum_units`
  DROP COLUMN `sort_order`;

COMMIT;
