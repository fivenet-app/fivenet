-- MariaDB dump 10.19  Distrib 10.10.3-MariaDB, for Linux (x86_64)
--
-- Host: 127.0.0.1    Database: fivenet
-- ------------------------------------------------------
-- Server version	10.11.3-MariaDB-1:10.11.3+maria~ubu2204

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
-- Dumping data for table `fivenet_permissions`
--

LOCK TABLES `fivenet_permissions` WRITE;
/*!40000 ALTER TABLE `fivenet_permissions` DISABLE KEYS */;
INSERT INTO `fivenet_permissions` VALUES
(1,'2023-03-18 16:37:43.273',NULL,'AuthService.GetCharacters','authservice-getcharacters',''),
(3,'2023-03-18 16:37:43.290',NULL,'CitizenStoreService.ListCitizens','citizenstoreservice-listusers',''),
(4,'2023-03-18 16:37:43.299',NULL,'CitizenStoreService.ListCitizens.Fields.Licenses','citizenstoreservice-listusers-licenses',''),
(5,'2023-03-18 16:37:43.307',NULL,'CitizenStoreService.ListCitizens.Fields.UserProps','citizenstoreservice-listusers-userprops',''),
(7,'2023-03-18 16:37:43.316',NULL,'CitizenStoreService.ListUserActivity','citizenstoreservice-listuseractivity',''),
(8,'2023-03-18 16:37:43.325',NULL,'CitizenStoreService.ListUserActivity.SourceUser','citizenstoreservice-listuseractivity-sourceuser',''),
(9,'2023-03-18 16:37:43.333',NULL,'CitizenStoreService.ListUserDocuments','citizenstoreservice-getuserdocuments',''),
(10,'2023-03-18 16:37:43.342',NULL,'CitizenStoreService.SetUserProps','citizenstoreservice-setuserprops',''),
(11,'2023-03-18 16:37:43.350',NULL,'CitizenStoreService.SetUserProps.Wanted','citizenstoreservice-setuserprops-wanted',''),
(12,'2023-03-18 16:37:43.359',NULL,'CompletorService.CompleteDocumentCategories','completorservice-completedocumentcategory',''),
(13,'2023-03-18 16:37:43.367',NULL,'CompletorService.CompleteDocumentCategories.ambulance','completorservice-completedocumentcategory-ambulance',''),
(14,'2023-03-18 16:37:43.376',NULL,'CompletorService.CompleteDocumentCategories.doj','completorservice-completedocumentcategory-doj',''),
(15,'2023-03-18 16:37:43.384',NULL,'CompletorService.CompleteDocumentCategories.fib','completorservice-completedocumentcategory-fib',''),
(16,'2023-03-18 16:37:43.393',NULL,'CompletorService.CompleteDocumentCategories.police','completorservice-completedocumentcategory-police',''),
(17,'2023-03-18 16:37:43.401',NULL,'CompletorService.CompleteJobGrades','completorservice-completejobgrades',''),
(19,'2023-03-18 16:37:43.418',NULL,'DocStoreService.AddDocumentReference','docstoreservice-adddocumentreference',''),
(20,'2023-03-18 16:37:43.426',NULL,'DocStoreService.AddDocumentRelation','docstoreservice-adddocumentrelation',''),
(21,'2023-03-18 16:37:43.435',NULL,'DocStoreService.CreateDocument','docstoreservice-createorupdatedocument',''),
(22,'2023-03-18 16:37:43.443',NULL,'DocStoreService.PostDocumentComment','docstoreservice-postdocumentcomment',''),
(23,'2023-03-18 16:37:43.460',NULL,'DocStoreService.ListDocuments','docstoreservice-listdocuments',''),
(24,'2023-03-18 16:37:43.478',NULL,'DocStoreService.GetDocument','docstoreservice-getdocument',''),
(25,'2023-03-18 16:37:43.497',NULL,'DocStoreService.GetDocumentAccess','docstoreservice-getdocumentaccess',''),
(26,'2023-03-18 16:37:43.514',NULL,'DocStoreService.GetDocumentComments','docstoreservice-getdocumentcomments',''),
(29,'2023-03-18 16:37:43.523',NULL,'DocStoreService.ListTemplates','docstoreservice-listtemplates',''),
(34,'2023-03-18 16:37:43.533',NULL,'DocStoreService.SetDocumentAccess','docstoreservice-setdocumentaccess',''),
(35,'2023-03-18 16:37:43.435',NULL,'DocStoreService.UpdateDocument','docstoreservice-updatedocument',''),
(36,'2023-03-18 16:37:43.542',NULL,'LivemapperService.Stream','livemapperservice-stream',''),
(37,'2023-03-18 16:37:43.550',NULL,'LivemapperService.Stream.ambulance','livemapperservice-stream-ambulance',''),
(38,'2023-03-18 16:37:43.559',NULL,'LivemapperService.Stream.doj','livemapperservice-stream-doj',''),
(39,'2023-03-18 16:37:43.568',NULL,'LivemapperService.Stream.fib','livemapperservice-stream-fib',''),
(40,'2023-03-18 16:37:43.576',NULL,'LivemapperService.Stream.police','livemapperservice-stream-police',''),
(41,'2023-03-18 16:37:43.585',NULL,'NotificatorService.GetNotifications','notificatorservice-getnotifications',''),
(44,'2023-03-18 16:37:43.594',NULL,'Overview.View','overview-view',''),
(45,'2023-03-18 18:00:58.866',NULL,'CompletorService.CompleteCitizens','completorservice-completecharnames','');
/*!40000 ALTER TABLE `fivenet_permissions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping data for table `fivenet_role_permissions`
--

LOCK TABLES `fivenet_role_permissions` WRITE;
/*!40000 ALTER TABLE `fivenet_role_permissions` DISABLE KEYS */;
INSERT INTO `fivenet_role_permissions` VALUES
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
(1,19),
(1,20),
(1,21),
(1,22),
(1,23),
(1,24),
(1,25),
(1,26),
(1,29),
(1,34),
(1,36),
(1,37),
(1,38),
(1,39),
(1,40),
(1,41),
(1,44),
(2,1),
(2,3),
(2,9),
(2,12),
(2,13),
(2,17),
(2,19),
(2,20),
(2,21),
(2,22),
(2,23),
(2,24),
(2,25),
(2,26),
(2,29),
(2,34),
(2,36),
(2,37),
(2,41),
(2,44),
(2,45),
(3,1),
(3,3),
(3,4),
(3,5),
(3,7),
(3,8),
(3,9),
(3,10),
(3,11),
(3,12),
(3,14),
(3,17),
(3,19),
(3,20),
(3,21),
(3,22),
(3,23),
(3,24),
(3,25),
(3,26),
(3,29),
(3,34),
(3,36),
(3,37),
(3,41),
(3,44),
(3,45),
(4,1),
(4,3),
(4,4),
(4,5),
(4,7),
(4,8),
(4,9),
(4,10),
(4,11),
(4,12),
(4,14),
(4,17),
(4,19),
(4,20),
(4,21),
(4,22),
(4,23),
(4,24),
(4,25),
(4,26),
(4,29),
(4,34),
(4,36),
(4,37),
(4,41),
(4,44),
(4,45),
(5,1),
(5,22),
(5,24),
(5,26),
(5,41),
(5,44);
/*!40000 ALTER TABLE `fivenet_role_permissions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping data for table `fivenet_roles`
--

LOCK TABLES `fivenet_roles` WRITE;
/*!40000 ALTER TABLE `fivenet_roles` DISABLE KEYS */;
INSERT INTO `fivenet_roles` VALUES
(2,'2023-03-18 14:45:33.015',NULL,'ambulance-1','ambulance-1','Role for ambulance (Rank: 1)'),
(3,'2023-03-18 14:45:33.185',NULL,'doj-1','doj-1','Role for doj (Rank: 1)'),
(4,'2023-03-18 14:45:33.368',NULL,'police-1','police-1','Role for police (Rank: 1)'),
(5,'2023-03-18 14:45:33.590',NULL,'unemployed-1','unemployed-1','Role for unemployed (Rank: 1)');
/*!40000 ALTER TABLE `fivenet_roles` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping data for table `fivenet_user_permissions`
--

LOCK TABLES `fivenet_user_permissions` WRITE;
/*!40000 ALTER TABLE `fivenet_user_permissions` DISABLE KEYS */;
/*!40000 ALTER TABLE `fivenet_user_permissions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping data for table `fivenet_user_roles`
--

LOCK TABLES `fivenet_user_roles` WRITE;
/*!40000 ALTER TABLE `fivenet_user_roles` DISABLE KEYS */;
INSERT INTO `fivenet_user_roles` VALUES
(1,2),
(2,2);
/*!40000 ALTER TABLE `fivenet_user_roles` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2023-03-18 19:02:19
