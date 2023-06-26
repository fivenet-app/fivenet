BEGIN;

-- Table: job_grades - Should already exist
-- CREATE TABLE IF NOT EXISTS `job_grades` (
--   `job_name` varchar(50) NOT NULL,
--   `grade` int(11) NOT NULL,
--   `name` varchar(50) NOT NULL,
--   `label` varchar(50) NOT NULL,
--   `salary` int(11) NOT NULL,
--   `skin_male` longtext NOT NULL,
--   `skin_female` longtext NOT NULL,
--   PRIMARY KEY (`job_name`,`grade`)
-- ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Table: jobs - Should already exist
-- CREATE TABLE IF NOT EXISTS `jobs` (
--   `name` varchar(50) NOT NULL,
--   `label` varchar(50) DEFAULT NULL,
--   PRIMARY KEY (`name`)
-- ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Table: owned_vehicles -- Should already exist
-- CREATE TABLE IF NOT EXISTS `owned_vehicles` (
--   `owner` varchar(64) DEFAULT NULL,
--   `plate` varchar(12) NOT NULL,
--   `model` varchar(60) NOT NULL,
--   `vehicle` longtext DEFAULT NULL,
--   `type` varchar(20) NOT NULL,
--   `stored` tinyint(1) NOT NULL DEFAULT 0,
--   `carseller` int(11) DEFAULT 0,
--   `owners` longtext DEFAULT NULL,
--   `trunk` longtext DEFAULT NULL,
--   PRIMARY KEY (`plate`),
--   UNIQUE KEY `IDX_OWNED_VEHICLES_OWNERPLATE` (`owner`,`plate`) USING BTREE,
--   KEY `IDX_OWNED_VEHICLES_OWNER` (`owner`),
--   KEY `IDX_OWNED_VEHICLES_OWNERTYPE` (`owner`,`type`),
--   KEY `IDX_OWNED_VEHICLES_OWNERRMODELTYPE` (`owner`,`model`,`type`)
-- ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
-- Add indexes for better sorting performance
set @x := (select count(*) from information_schema.statistics where table_name = 'owned_vehicles' and index_name = 'idx_owned_vehicles_model' and table_schema = database());
set @sql := if( @x > 0, 'select ''Vehicles model index exists.''', 'ALTER TABLE owned_vehicles ADD KEY `idx_owned_vehicles_model` (`model`);');
PREPARE stmt FROM @sql;
EXECUTE stmt;

set @x := (select count(*) from information_schema.statistics where table_name = 'owned_vehicles' and index_name = 'idx_owned_vehicles_type' and table_schema = database());
set @sql := if( @x > 0, 'select ''Vehicles type index exists.''', 'ALTER TABLE owned_vehicles ADD KEY `idx_owned_vehicles_type` (`type`);');
PREPARE stmt FROM @sql;
EXECUTE stmt;

-- Table: user_licenses - Should already exist
-- CREATE TABLE IF NOT EXISTS `user_licenses` (
--   `type` varchar(60) NOT NULL,
--   `owner` varchar(64) NOT NULL,
--   PRIMARY KEY (`type`,`owner`),
--   KEY `idx_user_licenses_owner` (`owner`) USING BTREE
-- ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Table: users - Should already exist
-- Add firstname + lastname fulltext index
set @x := (select count(*) from information_schema.statistics where table_name = 'users' and index_name = 'idx_users_firstname_lastname_fulltext' and table_schema = database());
set @sql := if( @x > 0, 'select ''users fulltext index exists.''', 'ALTER TABLE users ADD FULLTEXT KEY `idx_users_firstname_lastname_fulltext` (`firstname`,`lastname`);');
PREPARE stmt FROM @sql;
EXECUTE stmt;

-- Add dateofbirth index
set @x := (select count(*) from information_schema.statistics where table_name = 'users' and index_name = 'idx_users_dateofbirth' and table_schema = database());
set @sql := if( @x > 0, 'select ''users fulltext index exists.''', 'ALTER TABLE users ADD KEY `idx_users_dateofbirth` (`dateofbirth`);');
PREPARE stmt FROM @sql;
EXECUTE stmt;

-- Add job index
set @x := (select count(*) from information_schema.statistics where table_name = 'users' and index_name = 'idx_users_job' and table_schema = database());
set @sql := if( @x > 0, 'select ''users fulltext index exists.''', 'ALTER TABLE users ADD KEY `idx_users_job` (`job`);');
PREPARE stmt FROM @sql;
EXECUTE stmt;

COMMIT;
