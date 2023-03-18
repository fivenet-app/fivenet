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
(1,3),
(1,4),
(1,5),
(1,7),
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
(1,25),
(1,30),
(1,32),
(1,33),
(1,34),
(1,35),
(1,36),
(1,37),
(1,40);
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
(3,'2023-03-18 14:45:33.185',NULL,'job-doj-1','job-doj-1','Role for doj Job (Rank: 1)'),
(4,'2023-03-18 14:45:33.368',NULL,'job-police-1','job-police-1','Role for police Job (Rank: 1)'),
(5,'2023-03-18 14:45:33.590',NULL,'job-unemployed-1','job-unemployed-1','Role for unemployed Job (Rank: 1)');
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
/*!40000 ALTER TABLE `arpanet_user_roles` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2023-03-18 17:32:50
