BEGIN;

-- The legacy FK is intentionally not restored. The application uses
-- `fivenet_user` as its canonical user table.

COMMIT;
