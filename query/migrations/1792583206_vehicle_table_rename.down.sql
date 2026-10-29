BEGIN;

RENAME TABLE
  `fivenet_vehicle` TO `fivenet_owned_vehicles`,
  `fivenet_vehicle_activity` TO `fivenet_vehicles_activity`;

COMMIT;
