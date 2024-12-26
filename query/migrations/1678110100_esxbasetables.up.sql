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
-- ) ENGINE=InnoDB;

-- Table: jobs - Should already exist
-- CREATE TABLE IF NOT EXISTS `jobs` (
--   `name` varchar(50) NOT NULL,
--   `label` varchar(50) DEFAULT NULL,
--   PRIMARY KEY (`name`)
-- ) ENGINE=InnoDB;

-- Table: owned_vehicles -- Should already exist
-- CREATE TABLE IF NOT EXISTS `owned_vehicles` (
--   `owner` varchar(64) DEFAULT NULL,
--   `plate` varchar(12) NOT NULL,
--   `model` varchar(60) NOT NULL,
--   `vehicle` longtext DEFAULT NULL,
--   `type` varchar(20) NOT NULL, -- (Optional)
--   `stored` tinyint(1) NOT NULL DEFAULT 0,
--   `carseller` int(11) DEFAULT 0,
--   `owners` longtext DEFAULT NULL,
--   `trunk` longtext DEFAULT NULL,
--   PRIMARY KEY (`plate`),
--   UNIQUE KEY `IDX_OWNED_VEHICLES_OWNERPLATE` (`owner`,`plate`),
--   KEY `IDX_OWNED_VEHICLES_OWNER` (`owner`),
--   KEY `IDX_OWNED_VEHICLES_OWNERTYPE` (`owner`,`type`),
--   KEY `IDX_OWNED_VEHICLES_OWNERRMODELTYPE` (`owner`,`model`,`type`)
-- ) ENGINE=InnoDB;
-- Add indexes for better sorting performance
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

set @x := (select count(*) from information_schema.statistics where table_name = 'owned_vehicles' and index_name = 'idx_owned_vehicles_type' and table_schema = database());
set @sql := if( @x > 0, 'select ''owned_vehicles type index exists.''', 'ALTER TABLE owned_vehicles ADD KEY `idx_owned_vehicles_type` (`type`);');
PREPARE stmt FROM @sql;
EXECUTE stmt;

-- Table: user_licenses - Should already exist
-- CREATE TABLE IF NOT EXISTS `user_licenses` (
--   `type` varchar(60) NOT NULL,
--   `owner` varchar(64) NOT NULL,
--   PRIMARY KEY (`type`,`owner`),
--   KEY `idx_user_licenses_owner` (`owner`)
-- ) ENGINE=InnoDB;

-- Table: users - Should already exist
-- Add firstname + lastname fulltext index
set @x := (select count(*) from information_schema.statistics where table_name = 'users' and index_name = 'idx_users_firstname_lastname_fulltext' and table_schema = database());
set @sql := if( @x > 0, 'select ''users name fulltext index exists.''', 'ALTER TABLE users ADD FULLTEXT KEY `idx_users_firstname_lastname_fulltext` (`firstname`,`lastname`);');
PREPARE stmt FROM @sql;
EXECUTE stmt;

-- Add dateofbirth index
set @x := (select count(*) from information_schema.statistics where table_name = 'users' and index_name = 'idx_users_dateofbirth' and table_schema = database());
set @sql := if( @x > 0, 'select ''users dateofbirth index exists.''', 'ALTER TABLE users ADD KEY `idx_users_dateofbirth` (`dateofbirth`);');
PREPARE stmt FROM @sql;
EXECUTE stmt;

-- Add job index
set @x := (select count(*) from information_schema.statistics where table_name = 'users' and index_name = 'idx_users_job' and table_schema = database());
set @sql := if( @x > 0, 'select ''users job index exists.''', 'ALTER TABLE users ADD KEY `idx_users_job` (`job`);');
PREPARE stmt FROM @sql;
EXECUTE stmt;

COMMIT;
