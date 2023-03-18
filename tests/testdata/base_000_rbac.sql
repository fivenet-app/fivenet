-- MariaDB dump 10.19  Distrib 10.10.3-MariaDB, for Linux (x86_64)
--
-- Host: 127.0.0.1    Database: arpanet
-- ------------------------------------------------------
-- Server version	10.11.2-MariaDB-1:10.11.2+maria~ubu2204

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Dumping data for table `arpanet_role_permissions`
--

LOCK TABLES `arpanet_role_permissions` WRITE;
/*!40000 ALTER TABLE `arpanet_role_permissions` DISABLE KEYS */;
INSERT INTO `arpanet_role_permissions` VALUES
(1,1),
(1,2),
(1,4),
(1,5),
(1,6),
(1,8),
(1,9),
(1,10),
(1,11),
(1,12),
(1,13),
(1,14),
(1,15),
(1,16),
(1,17),
(1,18),
(1,19),
(1,20),
(1,21),
(1,22),
(1,23),
(1,24),
(1,25),
(1,26),
(1,27),
(1,28),
(1,29),
(1,30),
(1,31),
(1,32),
(1,33),
(1,34),
(1,35),
(1,36),
(1,39),
(1,44),
(1,46),
(1,47),
(1,48),
(1,49),
(1,50),
(1,51),
(1,52),
(1,53),
(1,54),
(18,1),
(18,2),
(18,4),
(18,10),
(18,13),
(18,14),
(18,18),
(18,19),
(18,23),
(18,24),
(18,29),
(18,30),
(18,31),
(18,32),
(18,33),
(18,34),
(18,35),
(18,36),
(18,39),
(18,44),
(18,47),
(18,48),
(18,52),
(18,53),
(18,54),
(21,1),
(21,2),
(37,1),
(37,2),
(45,1),
(45,2),
(71,1),
(71,2);
/*!40000 ALTER TABLE `arpanet_role_permissions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping data for table `arpanet_roles`
--

LOCK TABLES `arpanet_roles` WRITE;
/*!40000 ALTER TABLE `arpanet_roles` DISABLE KEYS */;
INSERT INTO `arpanet_roles` VALUES
(1,'2023-03-18 14:45:32.988',NULL,'masterofdisaster','masterofdisaster',''),
(2,'2023-03-18 14:45:33.015',NULL,'job-ambulance-1','job-ambulance-1','Role for ambulance Job (Rank: 1)'),
(3,'2023-03-18 14:45:33.031',NULL,'job-ambulance-2','job-ambulance-2','Role for ambulance Job (Rank: 2)'),
(4,'2023-03-18 14:45:33.039',NULL,'job-ambulance-3','job-ambulance-3','Role for ambulance Job (Rank: 3)'),
(5,'2023-03-18 14:45:33.047',NULL,'job-ambulance-4','job-ambulance-4','Role for ambulance Job (Rank: 4)'),
(6,'2023-03-18 14:45:33.056',NULL,'job-ambulance-5','job-ambulance-5','Role for ambulance Job (Rank: 5)'),
(7,'2023-03-18 14:45:33.064',NULL,'job-ambulance-6','job-ambulance-6','Role for ambulance Job (Rank: 6)'),
(8,'2023-03-18 14:45:33.072',NULL,'job-ambulance-7','job-ambulance-7','Role for ambulance Job (Rank: 7)'),
(9,'2023-03-18 14:45:33.080',NULL,'job-ambulance-8','job-ambulance-8','Role for ambulance Job (Rank: 8)'),
(10,'2023-03-18 14:45:33.088',NULL,'job-ambulance-9','job-ambulance-9','Role for ambulance Job (Rank: 9)'),
(11,'2023-03-18 14:45:33.096',NULL,'job-ambulance-10','job-ambulance-10','Role for ambulance Job (Rank: 10)'),
(12,'2023-03-18 14:45:33.104',NULL,'job-ambulance-11','job-ambulance-11','Role for ambulance Job (Rank: 11)'),
(13,'2023-03-18 14:45:33.112',NULL,'job-ambulance-12','job-ambulance-12','Role for ambulance Job (Rank: 12)'),
(14,'2023-03-18 14:45:33.120',NULL,'job-ambulance-13','job-ambulance-13','Role for ambulance Job (Rank: 13)'),
(15,'2023-03-18 14:45:33.128',NULL,'job-ambulance-14','job-ambulance-14','Role for ambulance Job (Rank: 14)'),
(16,'2023-03-18 14:45:33.136',NULL,'job-ambulance-15','job-ambulance-15','Role for ambulance Job (Rank: 15)'),
(17,'2023-03-18 14:45:33.145',NULL,'job-ambulance-16','job-ambulance-16','Role for ambulance Job (Rank: 16)'),
(18,'2023-03-18 14:45:33.153',NULL,'job-ambulance-17','job-ambulance-17','Role for ambulance Job (Rank: 17)'),
(19,'2023-03-18 14:45:33.161',NULL,'job-ambulance-18','job-ambulance-18','Role for ambulance Job (Rank: 18)'),
(20,'2023-03-18 14:45:33.169',NULL,'job-ambulance-19','job-ambulance-19','Role for ambulance Job (Rank: 19)'),
(21,'2023-03-18 14:45:33.177',NULL,'job-ambulance-20','job-ambulance-20','Role for ambulance Job (Rank: 20)'),
(22,'2023-03-18 14:45:33.185',NULL,'job-doj-1','job-doj-1','Role for doj Job (Rank: 1)'),
(23,'2023-03-18 14:45:33.193',NULL,'job-doj-2','job-doj-2','Role for doj Job (Rank: 2)'),
(24,'2023-03-18 14:45:33.201',NULL,'job-doj-3','job-doj-3','Role for doj Job (Rank: 3)'),
(25,'2023-03-18 14:45:33.209',NULL,'job-doj-4','job-doj-4','Role for doj Job (Rank: 4)'),
(26,'2023-03-18 14:45:33.217',NULL,'job-doj-5','job-doj-5','Role for doj Job (Rank: 5)'),
(27,'2023-03-18 14:45:33.225',NULL,'job-doj-6','job-doj-6','Role for doj Job (Rank: 6)'),
(28,'2023-03-18 14:45:33.234',NULL,'job-doj-7','job-doj-7','Role for doj Job (Rank: 7)'),
(29,'2023-03-18 14:45:33.246',NULL,'job-doj-8','job-doj-8','Role for doj Job (Rank: 8)'),
(30,'2023-03-18 14:45:33.254',NULL,'job-doj-9','job-doj-9','Role for doj Job (Rank: 9)'),
(31,'2023-03-18 14:45:33.262',NULL,'job-doj-10','job-doj-10','Role for doj Job (Rank: 10)'),
(32,'2023-03-18 14:45:33.270',NULL,'job-doj-11','job-doj-11','Role for doj Job (Rank: 11)'),
(33,'2023-03-18 14:45:33.278',NULL,'job-doj-12','job-doj-12','Role for doj Job (Rank: 12)'),
(34,'2023-03-18 14:45:33.286',NULL,'job-doj-13','job-doj-13','Role for doj Job (Rank: 13)'),
(35,'2023-03-18 14:45:33.294',NULL,'job-doj-14','job-doj-14','Role for doj Job (Rank: 14)'),
(36,'2023-03-18 14:45:33.302',NULL,'job-doj-15','job-doj-15','Role for doj Job (Rank: 15)'),
(37,'2023-03-18 14:45:33.310',NULL,'job-doj-16','job-doj-16','Role for doj Job (Rank: 16)'),
(38,'2023-03-18 14:45:33.319',NULL,'job-doj-17','job-doj-17','Role for doj Job (Rank: 17)'),
(39,'2023-03-18 14:45:33.327',NULL,'job-doj-18','job-doj-18','Role for doj Job (Rank: 18)'),
(40,'2023-03-18 14:45:33.335',NULL,'job-doj-19','job-doj-19','Role for doj Job (Rank: 19)'),
(41,'2023-03-18 14:45:33.343',NULL,'job-doj-20','job-doj-20','Role for doj Job (Rank: 20)'),
(42,'2023-03-18 14:45:33.351',NULL,'job-doj-21','job-doj-21','Role for doj Job (Rank: 21)'),
(43,'2023-03-18 14:45:33.360',NULL,'job-doj-22','job-doj-22','Role for doj Job (Rank: 22)'),
(44,'2023-03-18 14:45:33.368',NULL,'job-police-1','job-police-1','Role for police Job (Rank: 1)'),
(45,'2023-03-18 14:45:33.376',NULL,'job-police-2','job-police-2','Role for police Job (Rank: 2)'),
(46,'2023-03-18 14:45:33.385',NULL,'job-police-3','job-police-3','Role for police Job (Rank: 3)'),
(47,'2023-03-18 14:45:33.393',NULL,'job-police-4','job-police-4','Role for police Job (Rank: 4)'),
(48,'2023-03-18 14:45:33.401',NULL,'job-police-5','job-police-5','Role for police Job (Rank: 5)'),
(49,'2023-03-18 14:45:33.409',NULL,'job-police-6','job-police-6','Role for police Job (Rank: 6)'),
(50,'2023-03-18 14:45:33.417',NULL,'job-police-7','job-police-7','Role for police Job (Rank: 7)'),
(51,'2023-03-18 14:45:33.426',NULL,'job-police-8','job-police-8','Role for police Job (Rank: 8)'),
(52,'2023-03-18 14:45:33.434',NULL,'job-police-9','job-police-9','Role for police Job (Rank: 9)'),
(53,'2023-03-18 14:45:33.442',NULL,'job-police-10','job-police-10','Role for police Job (Rank: 10)'),
(54,'2023-03-18 14:45:33.450',NULL,'job-police-11','job-police-11','Role for police Job (Rank: 11)'),
(55,'2023-03-18 14:45:33.459',NULL,'job-police-12','job-police-12','Role for police Job (Rank: 12)'),
(56,'2023-03-18 14:45:33.467',NULL,'job-police-13','job-police-13','Role for police Job (Rank: 13)'),
(57,'2023-03-18 14:45:33.476',NULL,'job-police-14','job-police-14','Role for police Job (Rank: 14)'),
(58,'2023-03-18 14:45:33.484',NULL,'job-police-15','job-police-15','Role for police Job (Rank: 15)'),
(59,'2023-03-18 14:45:33.493',NULL,'job-police-16','job-police-16','Role for police Job (Rank: 16)'),
(60,'2023-03-18 14:45:33.501',NULL,'job-police-17','job-police-17','Role for police Job (Rank: 17)'),
(61,'2023-03-18 14:45:33.509',NULL,'job-police-18','job-police-18','Role for police Job (Rank: 18)'),
(62,'2023-03-18 14:45:33.517',NULL,'job-police-19','job-police-19','Role for police Job (Rank: 19)'),
(63,'2023-03-18 14:45:33.525',NULL,'job-police-20','job-police-20','Role for police Job (Rank: 20)'),
(64,'2023-03-18 14:45:33.533',NULL,'job-police-21','job-police-21','Role for police Job (Rank: 21)'),
(65,'2023-03-18 14:45:33.541',NULL,'job-police-22','job-police-22','Role for police Job (Rank: 22)'),
(66,'2023-03-18 14:45:33.549',NULL,'job-police-23','job-police-23','Role for police Job (Rank: 23)'),
(67,'2023-03-18 14:45:33.557',NULL,'job-police-24','job-police-24','Role for police Job (Rank: 24)'),
(68,'2023-03-18 14:45:33.566',NULL,'job-police-25','job-police-25','Role for police Job (Rank: 25)'),
(69,'2023-03-18 14:45:33.574',NULL,'job-police-26','job-police-26','Role for police Job (Rank: 26)'),
(70,'2023-03-18 14:45:33.582',NULL,'job-police-27','job-police-27','Role for police Job (Rank: 27)'),
(71,'2023-03-18 14:45:33.590',NULL,'job-unemployed-1','job-unemployed-1','Role for unemployed Job (Rank: 1)');
/*!40000 ALTER TABLE `arpanet_roles` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping data for table `arpanet_user_permissions`
--

LOCK TABLES `arpanet_user_permissions` WRITE;
/*!40000 ALTER TABLE `arpanet_user_permissions` DISABLE KEYS */;
/*!40000 ALTER TABLE `arpanet_user_permissions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping data for table `arpanet_user_roles`
--

LOCK TABLES `arpanet_user_roles` WRITE;
/*!40000 ALTER TABLE `arpanet_user_roles` DISABLE KEYS */;
INSERT INTO `arpanet_user_roles` VALUES
(1,18),
(2,21),
(3,37),
(4,45),
(5,71);
/*!40000 ALTER TABLE `arpanet_user_roles` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2023-03-18 15:59:40
