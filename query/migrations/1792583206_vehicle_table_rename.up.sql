BEGIN;

RENAME TABLE
  `fivenet_owned_vehicles` TO `fivenet_vehicle`,
  `fivenet_vehicles_activity` TO `fivenet_vehicle_activity`;

COMMIT;
