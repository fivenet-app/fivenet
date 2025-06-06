BEGIN;

-- Table: `fivenet_accounts_oauth2` - Add columns to store tokens
ALTER TABLE `fivenet_accounts_oauth2`
  ADD COLUMN `access_token` VARCHAR(512),
  ADD COLUMN `refresh_token` VARCHAR(512),
  ADD COLUMN `token_type` VARCHAR(16),
  ADD COLUMN `scope` VARCHAR(128),
  ADD COLUMN `expires_in` INT,
  ADD COLUMN `obtained_at` DATETIME;

COMMIT;
