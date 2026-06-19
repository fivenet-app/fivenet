-- gksphone_job_message definition
CREATE TABLE IF NOT EXISTS `gksphone_job_message` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` longtext,
  `number` varchar(50) DEFAULT NULL,
  `message` longtext,
  `photo` longtext,
  `gps` varchar(255) DEFAULT NULL,
  `owner` int NOT NULL DEFAULT '0',
  `jobm` varchar(255) DEFAULT NULL,
  `anon` varchar(50) DEFAULT NULL,
  `time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB;

-- gksphone_settings definition
CREATE TABLE IF NOT EXISTS `gksphone_settings` (
  `id` int NOT NULL AUTO_INCREMENT,
  `identifier` longtext,
  `crypto` varchar(535) DEFAULT NULL,
  `phone_number` varchar(50) DEFAULT NULL,
  `avatar_url` longtext,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB;

-- phone_phones definition
CREATE TABLE IF NOT EXISTS `phone_phones` (
  `id` varchar(100) NOT NULL,
  `owner_id` varchar(100) NOT NULL,
  `phone_number` varchar(15) NOT NULL,
  `name` varchar(50) DEFAULT NULL,
  `pin` varchar(4) DEFAULT NULL,
  `face_id` varchar(100) DEFAULT NULL,
  `settings` longtext,
  `is_setup` tinyint(1) DEFAULT '0',
  `assigned` tinyint(1) DEFAULT '0',
  `battery` int NOT NULL DEFAULT '100',
  `last_seen` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `phone_number` (`phone_number`)
) ENGINE=InnoDB;

-- phone_services_channels definition
CREATE TABLE IF NOT EXISTS `phone_services_channels` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `phone_number` varchar(15) NOT NULL,
  `company` varchar(50) NOT NULL,
  `last_message` varchar(100) DEFAULT NULL,
  `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB;

-- phone_services_messages definition
CREATE TABLE IF NOT EXISTS `phone_services_messages` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `channel_id` int unsigned NOT NULL,
  `sender` varchar(15) NOT NULL,
  `message` varchar(1000) NOT NULL,
  `x_pos` int DEFAULT NULL,
  `y_pos` int DEFAULT NULL,
  `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `channel_id` (`channel_id`),
  CONSTRAINT `phone_services_messages_ibfk_1` FOREIGN KEY (`channel_id`) REFERENCES `phone_services_channels` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB;
