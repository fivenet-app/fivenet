-- Table data: `fivenet_accounts` - Add 3 accounts with password `password`
INSERT INTO `fivenet_accounts`
(`id`, `enabled`, `username`, `password`, `license`, `reg_token`)
VALUES (1, 1, 'user-1', '$2y$10$QHt2PpQ3kYheZZTASOLY5uzpzoi30O9oYijIZabSE78a8yqfp7mjW', '3c7681d6f7ad895eb7b1cc05cf895c7f1d1622c4', NULL);
INSERT INTO `fivenet_accounts`
(`id`, `enabled`, `username`, `password`, `license`, `reg_token`)
VALUES (2, 1, 'user-2', '$2y$10$QHt2PpQ3kYheZZTASOLY5uzpzoi30O9oYijIZabSE78a8yqfp7mjW', 'fcee377a1fda007a8d2cc764a0a272e04d8c5d57', NULL);
INSERT INTO `fivenet_accounts`
(`id`, `enabled`, `username`, `password`, `license`, `reg_token`)
VALUES (3, 1, 'user-3', '$2y$10$QHt2PpQ3kYheZZTASOLY5uzpzoi30O9oYijIZabSE78a8yqfp7mjW', 'db7e039146d5bf1b6781e7bc1bef31f0bb1298ea', NULL);

-- Table data: `fivenet_user` - Add 5 chars into the database
INSERT INTO fivenet_user (id, account_id, identifier, `group`, job, job_grade, firstname, lastname, dateofbirth, sex, height, phone_number, disabled, visum, playtime, created_at, last_seen)
VALUES(1, 1, 'char1:3c7681d6f7ad895eb7b1cc05cf895c7f1d1622c4', 'user', 'ambulance', 17, 'Dr. Amy', 'Clockwork', '08.04.2003', 'f', 182, '3542786', 0, 139, 1493920, '2023-01-26 09:01:51.000', '2023-03-10 22:11:09.000');
INSERT INTO fivenet_user (id, account_id, identifier, `group`, job, job_grade, firstname, lastname, dateofbirth, sex, height, phone_number, disabled, visum, playtime, created_at, last_seen)
VALUES(2, 2, 'char1:fcee377a1fda007a8d2cc764a0a272e04d8c5d57', 'user', 'ambulance', 20, 'Philipp', 'Scott', '01.08.1982', 'm', 185, '1550044', 0, 209, 2244596, '2023-01-26 09:01:51.000', '2023-03-11 21:06:27.000');
INSERT INTO fivenet_user (id, account_id, identifier, `group`, job, job_grade, firstname, lastname, dateofbirth, sex, height, phone_number, disabled, visum, playtime, created_at, last_seen)
VALUES(3, 3, 'char1:db7e039146d5bf1b6781e7bc1bef31f0bb1298ea', 'user', 'doj', 16, 'Jonas', 'Striker', '28.10.1990', 'm', 186, '2488396', 0, 286, 3084976, '2023-01-26 09:01:51.000', '2023-03-11 21:06:27.000');
INSERT INTO fivenet_user (id, account_id, identifier, `group`, job, job_grade, firstname, lastname, dateofbirth, sex, height, phone_number, disabled, visum, playtime, created_at, last_seen)
VALUES(4, 2, 'char2:fcee377a1fda007a8d2cc764a0a272e04d8c5d57', 'user', 'police', 2, 'Hannibal', 'Scott', '15.06.1990', 'm', 180, '1550044', 0, 209, 2244596, '2023-01-26 09:01:51.000', '2023-03-11 21:06:27.000');
INSERT INTO fivenet_user (id, account_id, identifier, `group`, job, job_grade, firstname, lastname, dateofbirth, sex, height, phone_number, disabled, visum, playtime, created_at, last_seen)
VALUES(5, 1, 'char2:3c7681d6f7ad895eb7b1cc05cf895c7f1d1622c4', 'user', 'unemployed', 1, 'Peter', 'Hans', '10.02.1991', 'm', 178, '1550044', 0, 209, 2244596, '2023-01-26 09:01:51.000', '2023-03-11 21:06:27.000');

-- Table data: Add license types into the database
INSERT INTO fivenet_licenses (`type`, label) VALUES('aircraft', 'Aircraft License');
INSERT INTO fivenet_licenses (`type`, label) VALUES('boat', 'Boat License');
INSERT INTO fivenet_licenses (`type`, label) VALUES('dmv', 'DMV License');
INSERT INTO fivenet_licenses (`type`, label) VALUES('drive', 'Driving License');
INSERT INTO fivenet_licenses (`type`, label) VALUES('drive_bike', 'Bike License');
INSERT INTO fivenet_licenses (`type`, label) VALUES('drive_truck', 'Truck License');
INSERT INTO fivenet_licenses (`type`, label) VALUES('weapon', 'Weapon License');

-- Table data: Add some user licenses for more realistic testing that all data is retrieved from the database
INSERT INTO fivenet_user_licenses (`type`, `user_id`) VALUES('aircraft', 1);
INSERT INTO fivenet_user_licenses (`type`, `user_id`) VALUES('boat', 1);
INSERT INTO fivenet_user_licenses (`type`, `user_id`) VALUES('dmv', 1);
INSERT INTO fivenet_user_licenses (`type`, `user_id`) VALUES('dmv', 3);
INSERT INTO fivenet_user_licenses (`type`, `user_id`) VALUES('dmv', 2);
INSERT INTO fivenet_user_licenses (`type`, `user_id`) VALUES('dmv', 5);
INSERT INTO fivenet_user_licenses (`type`, `user_id`) VALUES('drive', 1);
INSERT INTO fivenet_user_licenses (`type`, `user_id`) VALUES('drive', 3);
INSERT INTO fivenet_user_licenses (`type`, `user_id`) VALUES('drive', 2);
INSERT INTO fivenet_user_licenses (`type`, `user_id`) VALUES('drive', 5);
INSERT INTO fivenet_user_licenses (`type`, `user_id`) VALUES('drive_bike', 1);
INSERT INTO fivenet_user_licenses (`type`, `user_id`) VALUES('drive_truck', 1);
INSERT INTO fivenet_user_licenses (`type`, `user_id`) VALUES('weapon', 1);
INSERT INTO fivenet_user_licenses (`type`, `user_id`) VALUES('weapon', 5);

-- Table data: jobs
INSERT INTO fivenet_jobs (`name`, label) VALUES('ambulance', 'LSMD');
INSERT INTO fivenet_jobs (`name`, label) VALUES('doj', 'DOJ');
INSERT INTO fivenet_jobs (`name`, label) VALUES('police', 'LSPD');
INSERT INTO fivenet_jobs (`name`, label) VALUES('unemployed', 'Unemployed');

-- Table data: fivenet_jobs_grades
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('ambulance', 1, 'Auszubildender Rettungshelfer');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('ambulance', 2, 'Rettungshelfer');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('ambulance', 3, 'Rettungssanitäter');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('ambulance', 4, 'Rettungsassistent');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('ambulance', 5, 'Notfallsanitäter');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('ambulance', 6, 'Assistenzarzt der Notfallmedizin');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('ambulance', 7, 'Assistenzarzt der Psychiatrie');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('ambulance', 8, 'Assistenzarzt der Chirurgie');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('ambulance', 9, 'Rettungsspezialist');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('ambulance', 10, 'Technischer Rettungsspezialist');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('ambulance', 11, 'Facharzt der Notfallmedizin');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('ambulance', 12, 'Facharzt der Psychiatrie');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('ambulance', 13, 'Facharzt der Chirurgie');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('ambulance', 14, 'Oberarzt der Notfallmedizin');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('ambulance', 15, 'Oberarzt der Psychiatrie');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('ambulance', 16, 'Oberarzt der Chirurgie');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('ambulance', 17, 'Leitender Oberarzt');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('ambulance', 18, 'Chefarzt');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('ambulance', 19, 'Stellvertr. Ärztlicher Direktor');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('ambulance', 20, 'Ärztlicher Direktor');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('doj', 1, 'Amtsassistent');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('doj', 2, 'Deputy Marshal');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('doj', 3, 'Senior Deputy Marshal');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('doj', 4, 'Supervisory Deputy Marshal');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('doj', 5, 'Assistant Chief Deputy Marshal');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('doj', 6, 'Chief Deputy Marshal');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('doj', 7, 'Marshal');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('doj', 8, 'Probationary Judge');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('doj', 9, 'Associate Judge');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('doj', 10, 'Judge');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('doj', 11, 'Senior Judge');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('doj', 12, 'Prosecutor');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('doj', 13, 'Senior Prosecutor');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('doj', 14, 'District Attorney');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('doj', 15, 'State Attorney');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('doj', 16, 'Senior State Attorney');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('doj', 17, 'Assistant Marshal Director');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('doj', 18, 'Deputy Chief Judge');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('doj', 19, 'Deputy Attorney General');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('doj', 20, 'Marshal Director');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('doj', 21, 'Chief Judge');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('doj', 22, 'Attorney General');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('police', 1, 'Rekrut');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('police', 2, 'Junior Officer');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('police', 3, 'Officer I');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('police', 4, 'Officer II');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('police', 5, 'Officer III');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('police', 6, 'Senior Officer');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('police', 7, 'Sergeant I');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('police', 8, 'Sergeant II');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('police', 9, 'Staff Sergeant');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('police', 10, 'Lieutenant');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('police', 11, 'NOOSE - Officer');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('police', 12, 'NOOSE - Sergeant');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('police', 13, 'NOOSE - Lieutenant');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('police', 14, 'SAHP - Officer');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('police', 15, 'SAHP - Sergeant');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('police', 16, 'SAHP - Lieutenant');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('police', 17, 'Inspector');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('police', 18, 'SAHP - Captain');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('police', 19, 'SAHP - Commander');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('police', 20, 'SWAT - Captain');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('police', 21, 'NOOSE - Commander');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('police', 22, 'Captain');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('police', 23, 'Commander');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('police', 24, 'Chief Inspector');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('police', 25, 'Deputy Chief');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('police', 26, 'Assistant Chief');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('police', 27, 'Chief of Police');
INSERT INTO fivenet_jobs_grades (job_name, grade, label) VALUES('unemployed', 1, 'Arbeitslos');

-- Table data: fivenet_owned_vehicles - Add some normal `aircraft`, `boat` and `car` vehicles to the database
-- aircraft
INSERT INTO fivenet_owned_vehicles (user_id, plate, model, `type`) VALUES(3, 'ABC 381', 'buzzard2', 'aircraft');
INSERT INTO fivenet_owned_vehicles (user_id, plate, model, `type`) VALUES(2, 'DEZ 725', 'supervolito', 'aircraft');
-- boat
INSERT INTO fivenet_owned_vehicles (user_id, plate, model, `type`) VALUES(5, 'BCJ 282', 'seashark', 'boat');
INSERT INTO fivenet_owned_vehicles (user_id, plate, model, `type`) VALUES(3, 'BEG 837', 'yaluxe', 'boat');
-- car
INSERT INTO fivenet_owned_vehicles (user_id, plate, model, `type`) VALUES(1, 'AMJV 079', 'xxxxx', 'car');
INSERT INTO fivenet_owned_vehicles (user_id, plate, model, `type`) VALUES(3, 'GFFT 070', 'veto2', 'car');
INSERT INTO fivenet_owned_vehicles (user_id, plate, model, `type`) VALUES(5, 'GWU 358', 'rmodsuprapandem', 'car');
-- car job specific
INSERT INTO fivenet_owned_vehicles (user_id, plate, model, `type`) VALUES(1, 'OENG 747', '16ramambo', 'car_ambulance');
INSERT INTO fivenet_owned_vehicles (user_id, plate, model, `type`) VALUES(2, 'XJMV 726', 'flyhoe', 'car_ambulance');
INSERT INTO fivenet_owned_vehicles (user_id, plate, model, `type`) VALUES(4, 'JMI 560', '18chgr', 'car_police');
INSERT INTO fivenet_owned_vehicles (user_id, plate, model, `type`) VALUES(4, 'NBH 257', 'poltahoe21', 'car_police');
INSERT INTO fivenet_owned_vehicles (user_id, plate, model, `type`) VALUES(4, 'QNT 765', 'poldurango', 'car_police');
