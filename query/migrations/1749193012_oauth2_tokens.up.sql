BEGIN;

-- Table: `fivenet_accounts_oauth2` - Add columns to store tokens
ALTER TABLE `fivenet_accounts_oauth2`
  ADD COLUMN `access_token` varchar(512),
  ADD COLUMN `refresh_token` varchar(512),
  ADD COLUMN `token_type` varchar(16),
  ADD COLUMN `scope` varchar(128),
  ADD COLUMN `expires_in` INT,
  ADD COLUMN `obtained_at` DATETIME;

COMMIT;
