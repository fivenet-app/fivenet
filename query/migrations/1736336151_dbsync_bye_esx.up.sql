BEGIN;

-- Table: fivenet_jobs
CREATE TABLE IF NOT EXISTS `fivenet_jobs` (
  `name` varchar(50) NOT NULL,
  `label` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`name`)
) ENGINE=InnoDB;

-- Table: fivenet_job_grades
CREATE TABLE IF NOT EXISTS `fivenet_job_grades` (
  `job_name` varchar(50) NOT NULL,
  `grade` int NOT NULL,
  `label` varchar(50) NOT NULL,
  PRIMARY KEY (`job_name`, `grade`),
  CONSTRAINT `fk_fivenet_job_grades_job_name` FOREIGN KEY (`job_name`) REFERENCES `fivenet_jobs` (`name`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Table: fivenet_licenses
CREATE TABLE IF NOT EXISTS `fivenet_licenses` (
  `type` varchar(60) NOT NULL,
  `label` varchar(60) NOT NULL,
  PRIMARY KEY (`type`)
) ENGINE=InnoDB;

-- Table: fivenet_users
CREATE TABLE IF NOT EXISTS `fivenet_users` (
  `id` int NOT NULL AUTO_INCREMENT,
  `identifier` varchar(64) NOT NULL,
  `group` varchar(50) DEFAULT NULL,
  `job` varchar(20) DEFAULT 'unemployed',
  `job_grade` int DEFAULT '0',
  `firstname` varchar(50) DEFAULT NULL,
  `lastname` varchar(50) DEFAULT NULL,
  `dateofbirth` varchar(25) DEFAULT NULL,
  `sex` varchar(10) DEFAULT NULL,
  `height` varchar(5) DEFAULT NULL,
  `phone_number` varchar(20) DEFAULT NULL,
  `disabled` tinyint(1) DEFAULT '0',
  `visum` int DEFAULT NULL,
  `playtime` int DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`identifier`),
  UNIQUE KEY `id` (`id`),
  KEY `idx_fivenet_users_job` (`job`),
  KEY `idx_fivenet_users_dateofbirth` (`dateofbirth`),
  FULLTEXT KEY `idx_fivenet_users_firstname_lastname_fulltext` (`firstname`, `lastname`)
) ENGINE=InnoDB;

-- Table: fivenet_user_licenses
CREATE TABLE IF NOT EXISTS `fivenet_user_licenses` (
  `type` varchar(60) NOT NULL,
  `owner` varchar(64) NOT NULL,
  PRIMARY KEY (`type`,`owner`),
  KEY `fivenet_user_licenses_owner_IDX` (`owner`),
  CONSTRAINT `fk_fivenet_user_licenses_type` FOREIGN KEY (`type`) REFERENCES `fivenet_licenses` (`type`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_user_licenses_owner` FOREIGN KEY (`owner`) REFERENCES `{{.UsersTableName}}` (`identifier`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Table: fivenet_owned_vehicles
CREATE TABLE IF NOT EXISTS `fivenet_owned_vehicles` (
  `owner` varchar(64) DEFAULT NULL,
  `plate` varchar(12) NOT NULL,
  `model` varchar(60) DEFAULT NULL,
  `type` varchar(20) NOT NULL,
  PRIMARY KEY (`plate`),
  UNIQUE KEY `idx_fivenet_owned_vehicles_ownerplate` (`owner`, `plate`),
  KEY `idx_fivenet_owned_vehicles_owner` (`owner`),
  KEY `idx_fivenet_owned_vehicles_owner_type` (`owner`, `type`),
  KEY `idx_fivenet_owned_vehicles_owner_model_type` (`owner`, `model`, `type`),
  KEY `idx_fivenet_owned_vehicles_model` (`model`),
  KEY `idx_fivenet_owned_vehicles_type` (`type`)
) ENGINE=InnoDB;

COMMIT;
