-- Table: users
CREATE TABLE IF NOT EXISTS `users` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `identifier` varchar(64) NOT NULL,
  `license` varchar(50) DEFAULT NULL,
  `group` varchar(50) DEFAULT NULL,
  `skin` longtext DEFAULT NULL,
  `job` varchar(20) DEFAULT 'unemployed',
  `job_grade` int(11) DEFAULT 1,
  `loadout` longtext DEFAULT NULL,
  `position` text DEFAULT NULL,
  `firstname` varchar(50) DEFAULT NULL,
  `lastname` varchar(50) DEFAULT NULL,
  `dateofbirth` varchar(25) DEFAULT NULL,
  `sex` varchar(10) DEFAULT NULL,
  `height` varchar(5) DEFAULT NULL,
  `is_dead` tinyint(1) DEFAULT 0,
  `last_property` varchar(255) DEFAULT NULL,
  `jail` int(11) NOT NULL DEFAULT 0,
  `inventory` longtext DEFAULT NULL,
  `phone_number` varchar(20) DEFAULT NULL,
  `accounts` longtext DEFAULT NULL,
  `tattoos` longtext DEFAULT NULL,
  `disabled` tinyint(1) DEFAULT 0,
  `visum` int(11) DEFAULT NULL,
  `playtime` int(11) DEFAULT NULL,
  `levelData` text DEFAULT NULL,
  `onDuty` tinyint(4) DEFAULT 0,
  `health` int(11) DEFAULT 200,
  `armor` int(11) DEFAULT 0,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP(),
  `last_seen` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(),
  `meta` longtext DEFAULT NULL,
  `metadata` longtext DEFAULT NULL,
  PRIMARY KEY (`identifier`),
  UNIQUE KEY `id` (`id`),
  KEY `idx_users_job` (`job`),
  KEY `idx_users_dateofbirth` (`dateofbirth`),
  FULLTEXT KEY `idx_users_firstname_lastname_fulltext` (`firstname`,`lastname`)
) ENGINE = InnoDB AUTO_INCREMENT = 1;
-- Table data: Add 5 chars into the database
INSERT INTO users (id, identifier, license, `group`, skin, job, job_grade, loadout, `position`, firstname, lastname, dateofbirth, sex, height, is_dead, last_property, jail, inventory, phone_number, accounts, tattoos, disabled, visum, playtime, levelData, onDuty, health, armor, created_at, last_seen, meta)
VALUES(1, 'char1:3c7681d6f7ad895eb7b1cc05cf895c7f1d1622c4', NULL, 'user', '{}', 'ambulance', 17, '{}', '{"z":49.45,"heading":141.73,"x":-1820.74,"y":-353.96}', 'Dr. Amy', 'Clockwork', '08.04.2003', 'f', '182', 0, NULL, 0, '{}', '3542786', '{}', '[]', 0, 139, 1493920, '{}', 0, 138, 0, '2023-01-26 09:01:51.000', '2023-03-10 22:11:09.000', NULL);
INSERT INTO users (id, identifier, license, `group`, skin, job, job_grade, loadout, `position`, firstname, lastname, dateofbirth, sex, height, is_dead, last_property, jail, inventory, phone_number, accounts, tattoos, disabled, visum, playtime, levelData, onDuty, health, armor, created_at, last_seen, meta)
VALUES(2, 'char1:fcee377a1fda007a8d2cc764a0a272e04d8c5d57', NULL, 'user', '{}', 'ambulance', 20, '{}', '{"x":981.96,"y":45.73,"z":80.99,"heading":342.99}', 'Philipp', 'Scott', '01.08.1982', 'm', '185', 0, NULL, 0, '{}', '1550044', '', '', 0, 209, 2244596, '', 0, 200, 0, '2023-01-26 09:01:51.000', '2023-03-11 21:06:27.000', NULL);
INSERT INTO users (id, identifier, license, `group`, skin, job, job_grade, loadout, `position`, firstname, lastname, dateofbirth, sex, height, is_dead, last_property, jail, inventory, phone_number, accounts, tattoos, disabled, visum, playtime, levelData, onDuty, health, armor, created_at, last_seen, meta)
VALUES(3, 'char1:db7e039146d5bf1b6781e7bc1bef31f0bb1298ea', NULL, 'user', '{}', 'doj', 16, '{}', '{"x":-561.13,"y":-176.81,"z":39.0,"heading":328.82}', 'Jonas', 'Striker', '28.10.1990', 'm', '186', 0, NULL, 0, '{}', '2488396', '', NULL, 0, 286, 3084976, '', 1, 200, 77, '2023-01-26 09:01:51.000', '2023-03-11 21:06:27.000', NULL);
INSERT INTO users (id, identifier, license, `group`, skin, job, job_grade, loadout, `position`, firstname, lastname, dateofbirth, sex, height, is_dead, last_property, jail, inventory, phone_number, accounts, tattoos, disabled, visum, playtime, levelData, onDuty, health, armor, created_at, last_seen, meta)
VALUES(4, 'char2:fcee377a1fda007a8d2cc764a0a272e04d8c5d57', NULL, 'user', '{}', 'police', 2, '{}', '{"x":981.96,"y":45.73,"z":80.99,"heading":342.99}', 'Hannibal', 'Scott', '15.06.1990', 'm', '180', 0, NULL, 0, '{}', '1550044', '', '', 0, 209, 2244596, '', 0, 200, 0, '2023-01-26 09:01:51.000', '2023-03-11 21:06:27.000', NULL);
INSERT INTO users (id, identifier, license, `group`, skin, job, job_grade, loadout, `position`, firstname, lastname, dateofbirth, sex, height, is_dead, last_property, jail, inventory, phone_number, accounts, tattoos, disabled, visum, playtime, levelData, onDuty, health, armor, created_at, last_seen, meta)
VALUES(5, 'char2:3c7681d6f7ad895eb7b1cc05cf895c7f1d1622c4', NULL, 'user', '{}', 'unemployed', 1, '{}', '{"x":981.96,"y":45.73,"z":80.99,"heading":342.99}', 'Peter', 'Hans', '10.02.1991', 'm', '178', 0, NULL, 0, '{}', '1550044', NULL, NULL, 0, 209, 2244596, NULL, 0, 200, 0, '2023-01-26 09:01:51.000', '2023-03-11 21:06:27.000', '2023-03-11 21:06:27.000');
-- Table: user_licenses
CREATE TABLE IF NOT EXISTS `user_licenses` (
  `type` varchar(60) NOT NULL,
  `owner` varchar(64) NOT NULL,
  PRIMARY KEY (`type`,`owner`),
  KEY `user_licenses_owner_IDX` (`owner`)
) ENGINE=InnoDB;
-- Table data: Add some user licenses for more realistic testing that all data is retrieved from the database
INSERT INTO user_licenses (`type`, owner) VALUES('aircraft', 'char1:3c7681d6f7ad895eb7b1cc05cf895c7f1d1622c4');
INSERT INTO user_licenses (`type`, owner) VALUES('boat', 'char1:3c7681d6f7ad895eb7b1cc05cf895c7f1d1622c4');
INSERT INTO user_licenses (`type`, owner) VALUES('dmv', 'char1:3c7681d6f7ad895eb7b1cc05cf895c7f1d1622c4');
INSERT INTO user_licenses (`type`, owner) VALUES('dmv', 'char1:db7e039146d5bf1b6781e7bc1bef31f0bb1298ea');
INSERT INTO user_licenses (`type`, owner) VALUES('dmv', 'char1:fcee377a1fda007a8d2cc764a0a272e04d8c5d57');
INSERT INTO user_licenses (`type`, owner) VALUES('dmv', 'char2:3c7681d6f7ad895eb7b1cc05cf895c7f1d1622c4');
INSERT INTO user_licenses (`type`, owner) VALUES('drive', 'char1:3c7681d6f7ad895eb7b1cc05cf895c7f1d1622c4');
INSERT INTO user_licenses (`type`, owner) VALUES('drive', 'char1:db7e039146d5bf1b6781e7bc1bef31f0bb1298ea');
INSERT INTO user_licenses (`type`, owner) VALUES('drive', 'char1:fcee377a1fda007a8d2cc764a0a272e04d8c5d57');
INSERT INTO user_licenses (`type`, owner) VALUES('drive', 'char2:3c7681d6f7ad895eb7b1cc05cf895c7f1d1622c4');
INSERT INTO user_licenses (`type`, owner) VALUES('drive_bike', 'char1:3c7681d6f7ad895eb7b1cc05cf895c7f1d1622c4');
INSERT INTO user_licenses (`type`, owner) VALUES('drive_truck', 'char1:3c7681d6f7ad895eb7b1cc05cf895c7f1d1622c4');
INSERT INTO user_licenses (`type`, owner) VALUES('weapon', 'char1:3c7681d6f7ad895eb7b1cc05cf895c7f1d1622c4');
INSERT INTO user_licenses (`type`, owner) VALUES('weapon', 'char2:3c7681d6f7ad895eb7b1cc05cf895c7f1d1622c4');
-- Table: jobs
CREATE TABLE IF NOT EXISTS `jobs` (
  `name` varchar(50) NOT NULL,
  `label` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`name`)
) ENGINE = InnoDB;
-- Table data: jobs
INSERT INTO jobs (name, label) VALUES('ambulance', 'LSMD');
INSERT INTO jobs (name, label) VALUES('doj', 'DOJ');
INSERT INTO jobs (name, label) VALUES('police', 'LSPD');
INSERT INTO jobs (name, label) VALUES('unemployed', 'Unemployed');
-- Table: job_grades
CREATE TABLE IF NOT EXISTS `job_grades` (
  `job_name` varchar(50) NOT NULL,
  `grade` int(11) NOT NULL,
  `name` varchar(50) NOT NULL,
  `label` varchar(50) NOT NULL,
  `salary` int(11) NOT NULL,
  `skin_male` longtext NOT NULL,
  `skin_female` longtext NOT NULL,
  PRIMARY KEY (`job_name`, `grade`)
) ENGINE = InnoDB;
-- Table data: job_grades
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('ambulance', 1, 'auszubildender_rettungshelfer', 'Auszubildender Rettungshelfer', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('ambulance', 2, 'rettungshelfer', 'Rettungshelfer', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('ambulance', 3, 'rettungssanitäter', 'Rettungssanitäter', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('ambulance', 4, 'rettungsassistent', 'Rettungsassistent', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('ambulance', 5, 'notfallsanitäter', 'Notfallsanitäter', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('ambulance', 6, 'assistenzarzt_der_notfallmedizin', 'Assistenzarzt der Notfallmedizin', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('ambulance', 7, 'assistenzarzt_der_psychiatrie', 'Assistenzarzt der Psychiatrie', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('ambulance', 8, 'assistenzarzt_der_chirurgie', 'Assistenzarzt der Chirurgie', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('ambulance', 9, 'rettungsspezialist', 'Rettungsspezialist', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('ambulance', 10, 'technischer_rettungsspezialist', 'Technischer Rettungsspezialist', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('ambulance', 11, 'facharzt_der_notfallmedizin', 'Facharzt der Notfallmedizin', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('ambulance', 12, 'facharzt_der_psychiatrie', 'Facharzt der Psychiatrie', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('ambulance', 13, 'facharzt_der_chirurgie', 'Facharzt der Chirurgie', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('ambulance', 14, 'oberarzt_der_notfallmedizin', 'Oberarzt der Notfallmedizin', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('ambulance', 15, 'oberarzt_der_psychiatrie', 'Oberarzt der Psychiatrie', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('ambulance', 16, 'oberarzt_der_chirurgie', 'Oberarzt der Chirurgie', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('ambulance', 17, 'leitender_oberarzt', 'Leitender Oberarzt', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('ambulance', 18, 'chefarzt', 'Chefarzt', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('ambulance', 19, 'stellvertr_ärztlicher_direktor', 'Stellvertr. Ärztlicher Direktor', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('ambulance', 20, 'ärztlicher_direktor', 'Ärztlicher Direktor', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('doj', 1, 'amtsassistent', 'Amtsassistent', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('doj', 2, 'deputy_marshal', 'Deputy Marshal', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('doj', 3, 'senior_deputy_marshal', 'Senior Deputy Marshal', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('doj', 4, 'supervisory_deputy_marshal', 'Supervisory Deputy Marshal', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('doj', 5, 'assistant_chief_deputy_marshal', 'Assistant Chief Deputy Marshal', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('doj', 6, 'chief_deputy_marshal', 'Chief Deputy Marshal', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('doj', 7, 'marshal', 'Marshal', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('doj', 8, 'probationary_judge', 'Probationary Judge', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('doj', 9, 'associate_judge', 'Associate Judge', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('doj', 10, 'judge', 'Judge', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('doj', 11, 'senior_judge', 'Senior Judge', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('doj', 12, 'prosecutor', 'Prosecutor', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('doj', 13, 'senior_prosecutor', 'Senior Prosecutor', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('doj', 14, 'district_attorney', 'District Attorney', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('doj', 15, 'state_attorney', 'State Attorney', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('doj', 16, 'senior_state_attorney', 'Senior State Attorney', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('doj', 17, 'assistant_marshal_director', 'Assistant Marshal Director', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('doj', 18, 'deputy_chief_judge', 'Deputy Chief Judge', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('doj', 19, 'deputy_attorney_general', 'Deputy Attorney General', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('doj', 20, 'marshal_director', 'Marshal Director', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('doj', 21, 'chief_judge', 'Chief Judge', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('doj', 22, 'attorney_general', 'Attorney General', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('police', 1, 'rekrut', 'Rekrut', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('police', 2, 'junior_officer', 'Junior Officer', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('police', 3, 'officer_i', 'Officer I', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('police', 4, 'officer_ii', 'Officer II', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('police', 5, 'officer_iii', 'Officer III', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('police', 6, 'senior_officer', 'Senior Officer', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('police', 7, 'sergeant_i', 'Sergeant I', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('police', 8, 'sergeant_ii', 'Sergeant II', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('police', 9, 'staff_sergeant', 'Staff Sergeant', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('police', 10, 'lieutenant', 'Lieutenant', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('police', 11, 'noose_officer', 'NOOSE - Officer', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('police', 12, 'noose_sergeant', 'NOOSE - Sergeant', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('police', 13, 'noose_lieutenant', 'NOOSE - Lieutenant', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('police', 14, 'sahp_officer', 'SAHP - Officer', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('police', 15, 'sahp_sergeant', 'SAHP - Sergeant', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('police', 16, 'sahp_lieutenant', 'SAHP - Lieutenant', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('police', 17, 'inspector', 'Inspector', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('police', 18, 'sahp_captain', 'SAHP - Captain', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('police', 19, 'sahp_commander', 'SAHP - Commander', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('police', 20, 'noose_captain', 'SWAT - Captain', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('police', 21, 'noose_commander', 'NOOSE - Commander', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('police', 22, 'captain', 'Captain', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('police', 23, 'commander', 'Commander', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('police', 24, 'chief_inspector', 'Chief Inspector', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('police', 25, 'deputy_chief', 'Deputy Chief', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('police', 26, 'assistant_chief', 'Assistant Chief', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('police', 27, 'chief_of_police', 'Chief of Police', 0, '{}', '{}');
INSERT INTO job_grades (job_name, grade, name, label, salary, skin_male, skin_female) VALUES('unemployed', 1, 'arbeitslos', 'Arbeitslos', 0, '{}', '{}');
-- Table: owned_vehicles
CREATE TABLE IF NOT EXISTS `owned_vehicles` (
  `owner` varchar(64) DEFAULT NULL,
  `plate` varchar(12) NOT NULL,
  `model` varchar(60) NOT NULL,
  `vehicle` longtext DEFAULT NULL,
  `type` varchar(20) NOT NULL,
  `stored` tinyint(1) NOT NULL DEFAULT 0,
  `carseller` int(11) DEFAULT 0,
  `owners` longtext DEFAULT NULL,
  `trunk` longtext DEFAULT NULL,
  PRIMARY KEY (`plate`),
  UNIQUE KEY `IDX_OWNED_VEHICLES_OWNERPLATE` (`owner`,`plate`),
  KEY `IDX_OWNED_VEHICLES_OWNER` (`owner`),
  KEY `IDX_OWNED_VEHICLES_OWNERTYPE` (`owner`,`type`),
  KEY `IDX_OWNED_VEHICLES_OWNERRMODELTYPE` (`owner`,`model`,`type`)
) ENGINE=InnoDB;
-- Table data: owned_vehicles - Add some normal `aircraft`, `boat` and `car` vehicles to the database
-- aircraft
INSERT INTO owned_vehicles (owner, plate, model, vehicle, `type`, `stored`, carseller, owners, trunk) VALUES('char1:db7e039146d5bf1b6781e7bc1bef31f0bb1298ea', 'ABC 381', 'buzzard2', '{}', 'aircraft', 1, 0, NULL, NULL);
INSERT INTO owned_vehicles (owner, plate, model, vehicle, `type`, `stored`, carseller, owners, trunk) VALUES('char1:fcee377a1fda007a8d2cc764a0a272e04d8c5d57', 'DEZ 725', 'supervolito', '{}', 'aircraft', 1, 0, NULL, '{}');
-- boat
INSERT INTO owned_vehicles (owner, plate, model, vehicle, `type`, `stored`, carseller, owners, trunk) VALUES('char2:3c7681d6f7ad895eb7b1cc05cf895c7f1d1622c4', 'BCJ 282', 'seashark', '{}', 'boat', 1, 0, NULL, NULL);
INSERT INTO owned_vehicles (owner, plate, model, vehicle, `type`, `stored`, carseller, owners, trunk) VALUES('char1:db7e039146d5bf1b6781e7bc1bef31f0bb1298ea', 'BEG 837', 'yaluxe', '{}', 'boat', 1, 0, NULL, NULL);
-- car
INSERT INTO owned_vehicles (owner, plate, model, vehicle, `type`, `stored`, carseller, owners, trunk) VALUES('char1:3c7681d6f7ad895eb7b1cc05cf895c7f1d1622c4', 'AMJV 079', 'xxxxx', '{}', 'car', 1, 0, NULL, NULL);
INSERT INTO owned_vehicles (owner, plate, model, vehicle, `type`, `stored`, carseller, owners, trunk) VALUES('char1:db7e039146d5bf1b6781e7bc1bef31f0bb1298ea', 'GFFT 070', 'veto2', '{}', 'car', 1, 0, NULL, NULL);
INSERT INTO owned_vehicles (owner, plate, model, vehicle, `type`, `stored`, carseller, owners, trunk) VALUES('char2:3c7681d6f7ad895eb7b1cc05cf895c7f1d1622c4', 'GWU 358', 'rmodsuprapandem', '{}', 'car', 1, 0, NULL, NULL);
-- car job specific
INSERT INTO owned_vehicles (owner, plate, model, vehicle, `type`, `stored`, carseller, owners, trunk) VALUES('char1:3c7681d6f7ad895eb7b1cc05cf895c7f1d1622c4', 'OENG 747', '16ramambo', '{}', 'car_ambulance', 1, 0, NULL, NULL);
INSERT INTO owned_vehicles (owner, plate, model, vehicle, `type`, `stored`, carseller, owners, trunk) VALUES('char1:fcee377a1fda007a8d2cc764a0a272e04d8c5d57', 'XJMV 726', 'flyhoe', '{}', 'car_ambulance', 1, 0, NULL, '{}');
INSERT INTO owned_vehicles (owner, plate, model, vehicle, `type`, `stored`, carseller, owners, trunk) VALUES('char2:fcee377a1fda007a8d2cc764a0a272e04d8c5d57', 'JMI 560', '18chgr', '{}', 'car_police', 1, 0, NULL, NULL);
INSERT INTO owned_vehicles (owner, plate, model, vehicle, `type`, `stored`, carseller, owners, trunk) VALUES('char2:fcee377a1fda007a8d2cc764a0a272e04d8c5d57', 'NBH 257', 'poltahoe21', '{}', 'car_police', 1, 0, NULL, NULL);
INSERT INTO owned_vehicles (owner, plate, model, vehicle, `type`, `stored`, carseller, owners, trunk) VALUES('char2:fcee377a1fda007a8d2cc764a0a272e04d8c5d57', 'QNT 765', 'poldurango', '{}', 'car_police', 1, 0, NULL, NULL);
