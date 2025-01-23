BEGIN;

-- Table: `job_grades` - Should already exist! This is the bare minimum structure required by FiveNet (column order doesn't matter).
--
-- CREATE TABLE IF NOT EXISTS `job_grades` (
--   `job_name` varchar(50) NOT NULL,
--   `grade` int(11) NOT NULL,
--   `name` varchar(50) NOT NULL,
--   `label` varchar(50) NOT NULL,
--   PRIMARY KEY (`job_name`,`grade`)
-- ) ENGINE=InnoDB;

-- Table: `jobs` - Should already exist
-- CREATE TABLE IF NOT EXISTS `jobs` (
--   `name` varchar(50) NOT NULL,
--   `label` varchar(50) DEFAULT NULL,
--   PRIMARY KEY (`name`)
-- ) ENGINE=InnoDB;

-- Table: `owned_vehicles` -- Should already exist! This is the bare minimum structure required by FiveNet (column order doesn't matter).
--
-- CREATE TABLE IF NOT EXISTS `owned_vehicles` (
--   `owner` varchar(64) DEFAULT NULL,
--   `plate` varchar(12) NOT NULL,
--   `model` varchar(60) DEFAULT NULL,
--   `type` varchar(20) NOT NULL,
--   PRIMARY KEY (`plate`),
--   UNIQUE KEY `idx_fivenet_owned_vehicles_ownerplate` (`owner`, `plate`),
--   KEY `idx_fivenet_owned_vehicles_owner` (`owner`),
--   KEY `idx_fivenet_owned_vehicles_owner_type` (`owner`, `type`),
--   KEY `idx_fivenet_owned_vehicles_owner_model_type` (`owner`, `model`, `type`),
--   KEY `idx_fivenet_owned_vehicles_model` (`model`),
--   KEY `idx_fivenet_owned_vehicles_type` (`type`)
-- ) ENGINE=InnoDB;

-- Table: owned_vehicles - Add `model` index for better sorting performance if the table exists
set @x := (
select
	1
from
	information_schema.statistics stats
inner join information_schema.columns on
	(columns.table_name = 'owned_vehicles'
		and columns.column_name = 'model'
		and columns.table_schema = database())
where
	stats.table_name = 'owned_vehicles'
	and stats.index_name = 'idx_owned_vehicles_model'
	and stats.table_schema = database()
);
set @sql := if( @x is null or @x > 0, 'select ''owned_vehicles model colum doesnt exist or index already exists.''', 'ALTER TABLE owned_vehicles ADD KEY `idx_owned_vehicles_model` (`model`);');
PREPARE stmt FROM @sql;
EXECUTE stmt;

-- Table: owned_vehicles - Add `type` index for better sorting performance if the table exists
set @x := (select 1 from information_schema.statistics where table_name = 'owned_vehicles' and index_name = 'idx_owned_vehicles_type' and table_schema = database());
set @sql := if( @x is null or @x > 0, 'select ''owned_vehicles type index exists.''', 'ALTER TABLE owned_vehicles ADD KEY `idx_owned_vehicles_type` (`type`);');
PREPARE stmt FROM @sql;
EXECUTE stmt;

-- Table: `user_licenses` - Should already exist! This is the bare minimum structure required by FiveNet (column order doesn't matter).
--
-- CREATE TABLE IF NOT EXISTS `user_licenses` (
--   `type` varchar(60) NOT NULL,
--   `owner` varchar(64) NOT NULL,
--   PRIMARY KEY (`type`,`owner`),
--   KEY `idx_user_licenses_owner` (`owner`)
-- ) ENGINE=InnoDB;

-- Table: `users` - Should already exist! This is the bare minimum structure required by FiveNet (column order doesn't matter).
--
-- CREATE TABLE IF NOT EXISTS `users` (
--   `id` int(11) NOT NULL AUTO_INCREMENT,
--   `identifier` varchar(64) NOT NULL,
--   `group` varchar(50) DEFAULT NULL,
--   `job` varchar(20) DEFAULT 'unemployed',
--   `job_grade` int(11) DEFAULT 1,
--   `firstname` varchar(50) DEFAULT NULL,
--   `lastname` varchar(50) DEFAULT NULL,
--   `dateofbirth` varchar(25) DEFAULT NULL,
--   `sex` varchar(10) DEFAULT NULL,
--   `height` varchar(5) DEFAULT NULL,
--   `phone_number` varchar(20) DEFAULT NULL,
--   `disabled` tinyint(1) DEFAULT '0',
--   `visum` int(11) DEFAULT NULL, -- Optional
--   `playtime` int(11) DEFAULT NULL, -- Optional
--   `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
--   PRIMARY KEY (`identifier`),
--   UNIQUE KEY `id` (`id`),
--   KEY `idx_users_job` (`job`),
--   KEY `idx_users_dateofbirth` (`dateofbirth`),
--   FULLTEXT KEY `idx_users_firstname_lastname_fulltext` (`firstname`,`lastname`)
-- ) ENGINE = InnoDB AUTO_INCREMENT = 1;

-- Table: `users` - Add `firstname` + `lastname` fulltext index if the table exists
set @x := (select 1 from information_schema.statistics where table_name = 'users' and index_name = 'idx_users_firstname_lastname_fulltext' and table_schema = database());
set @sql := if( @x is null or @x > 0, 'select ''users name fulltext index exists.''', 'ALTER TABLE users ADD FULLTEXT KEY `idx_users_firstname_lastname_fulltext` (`firstname`,`lastname`);');
PREPARE stmt FROM @sql;
EXECUTE stmt;

-- Table: `users` - Add `dateofbirth` column index if the table exists
set @x := (select 1 from information_schema.statistics where table_name = 'users' and index_name = 'idx_users_dateofbirth' and table_schema = database());
set @sql := if( @x is null or @x > 0, 'select ''users dateofbirth index exists.''', 'ALTER TABLE users ADD KEY `idx_users_dateofbirth` (`dateofbirth`);');
PREPARE stmt FROM @sql;
EXECUTE stmt;

-- Table: `users` - Add `job` column index if the table exists
set @x := (select 1 from information_schema.statistics where table_name = 'users' and index_name = 'idx_users_job' and table_schema = database());
set @sql := if( @x is null or @x > 0, 'select ''users job index exists.''', 'ALTER TABLE users ADD KEY `idx_users_job` (`job`);');
PREPARE stmt FROM @sql;
EXECUTE stmt;

COMMIT;
