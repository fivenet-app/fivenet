-- MariaDB dump 10.19  Distrib 10.11.3-MariaDB, for Linux (x86_64)
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
-- Dumping data for table `fivenet_attrs`
--

LOCK TABLES `fivenet_attrs` WRITE;
/*!40000 ALTER TABLE `fivenet_attrs` DISABLE KEYS */;
INSERT INTO `fivenet_attrs` (`id`, `created_at`, `permission_id`, `key`, `type`, `valid_values`) VALUES (1,'2023-05-12 18:00:13.327',7595,'Jobs','JobGradeList',NULL),
(2,'2023-05-12 18:00:13.345',3,'Fields','StringList','{"stringList":{"strings":["PhoneNumber","Licenses","UserProps.Wanted","UserProps.Job","UserProps.TrafficInfractionPoints","UserProps.OpenFines","UserProps.BloodType"]}}'),
(3,'2023-05-12 18:00:13.356',6,'Fields','StringList','{"stringList":{"strings":["SourceUser","Own"]}}'),
(4,'2023-05-12 18:00:13.366',8,'Fields','StringList','{"stringList":{"strings":["Wanted","Job","TrafficInfractionPoints"]}}'),
(5,'2023-05-12 18:00:13.378',11,'Jobs','JobList','{"jobList":{}}'),
(6,'2023-05-12 18:00:13.410',31,'Dispatches','JobList','{"jobList":{}}'),
(7,'2023-05-12 18:00:13.419',31,'Players','JobList','{"jobList":{}}'),
(8,'2023-05-12 18:00:13.434',77,'Jobs','JobList','{"jobList":{}}'),
(9,'2023-05-24 17:11:04.280',5497,'Access','StringList','{"stringList":{"strings":["Own","Lower_Rank","Same_Rank","Any"]}}'),
(10,'2023-05-24 17:11:04.291',21,'Access','StringList','{"stringList":{"strings":["Own","Lower_Rank","Same_Rank","Any"]}}'),
(11,'2023-05-24 17:11:04.304',30,'Access','StringList','{"stringList":{"strings":["Own","Lower_Rank","Same_Rank","Any"]}}');
/*!40000 ALTER TABLE `fivenet_attrs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping data for table `fivenet_permissions`
--

LOCK TABLES `fivenet_permissions` WRITE;
/*!40000 ALTER TABLE `fivenet_permissions` DISABLE KEYS */;
INSERT INTO `fivenet_permissions` (`id`, `created_at`, `category`, `name`, `guard_name`) VALUES (1,'2023-03-31 15:35:21.865','AuthService','ChooseCharacter','authservice-choosecharacter'),
(3,'2023-03-31 15:35:21.920','CitizenStoreService','ListCitizens','citizenstoreservice-listcitizens'),
(6,'2023-03-31 15:35:22.001','CitizenStoreService','ListUserActivity','citizenstoreservice-listuseractivity'),
(8,'2023-03-31 15:35:22.057','CitizenStoreService','SetUserProps','citizenstoreservice-setuserprops'),
(10,'2023-03-31 15:35:22.109','CompletorService','CompleteCitizens','completorservice-completecitizens'),
(11,'2023-03-31 15:35:22.140','CompletorService','CompleteCategories','completorservice-completedocumentcategories'),
(16,'2023-03-31 15:35:22.283','CompletorService','CompleteJobs','completorservice-completejobs'),
(17,'2023-03-31 15:35:22.310','DMVService','ListVehicles','dmvservice-listvehicles'),
(18,'2023-03-31 15:35:22.338','DocStoreService','AddDocumentReference','docstoreservice-adddocumentreference'),
(19,'2023-03-31 15:35:22.367','DocStoreService','AddDocumentRelation','docstoreservice-adddocumentrelation'),
(20,'2023-03-31 15:35:22.387','DocStoreService','CreateDocument','docstoreservice-createdocument'),
(21,'2023-03-31 15:35:22.414','DocStoreService','DeleteComment','docstoreservice-deletecomment'),
(22,'2023-03-31 15:35:22.440','DocStoreService','ListDocuments','docstoreservice-listdocuments'),
(24,'2023-03-31 15:35:22.498','DocStoreService','GetDocumentAccess','docstoreservice-getdocumentaccess'),
(26,'2023-03-31 15:35:22.554','DocStoreService','ListUserDocuments','docstoreservice-listuserdocuments'),
(27,'2023-03-31 15:35:22.577','DocStoreService','ListTemplates','docstoreservice-listtemplates'),
(28,'2023-03-31 15:35:22.606','DocStoreService','PostComment','docstoreservice-postcomment'),
(29,'2023-03-31 15:35:22.626','DocStoreService','SetDocumentAccess','docstoreservice-setdocumentaccess'),
(30,'2023-03-31 15:35:22.651','DocStoreService','UpdateDocument','docstoreservice-updatedocument'),
(31,'2023-03-31 15:35:22.680','LivemapperService','Stream','livemapperservice-stream'),
(74,'2023-04-03 12:17:15.106','RectorService','UpdateRolePerms','rectorservice-updateroleperms'),
(75,'2023-04-03 12:17:15.119','RectorService','CreateRole','rectorservice-createrole'),
(76,'2023-04-03 12:17:15.132','RectorService','DeleteRole','rectorservice-deleterole'),
(77,'2023-04-03 12:17:15.144','RectorService','GetPermissions','rectorservice-getpermissions'),
(82,'2023-04-03 12:17:15.208','RectorService','GetRoles','rectorservice-getroles'),
(2198,'2023-04-05 17:52:13.377','DocStoreService','CreateOrUpdateCategory','docstoreservice-createorupdatecategory'),
(2199,'2023-04-05 17:52:13.438','DocStoreService','DeleteCategory','docstoreservice-deletecategory'),
(3056,'2023-04-06 20:26:04.619','DocStoreService','ListCategories','docstoreservice-listcategories'),
(3392,'2023-04-08 16:51:59.783','DocStoreService','CreateTemplate','docstoreservice-createtemplate'),
(3395,'2023-04-08 16:51:59.983','DocStoreService','DeleteTemplate','docstoreservice-deletetemplate'),
(3671,'2023-04-09 02:04:57.236','RectorService','ViewAuditLog','rectorservice-viewauditlog'),
(5497,'2023-04-11 17:38:42.619','DocStoreService','DeleteDocument','docstoreservice-deletedocument'),
(7595,'2023-04-18 15:40:04.382','CitizenStoreService','GetUser','citizenstoreservice-getuser'),
(8848,'2023-04-23 20:44:58.283','RectorService','GetJobProps','rectorservice-getjobprops'),
(8855,'2023-04-23 20:44:58.851','RectorService','SetJobProps','rectorservice-setjobprops');
/*!40000 ALTER TABLE `fivenet_permissions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping data for table `fivenet_roles`
--

LOCK TABLES `fivenet_roles` WRITE;
/*!40000 ALTER TABLE `fivenet_roles` DISABLE KEYS */;
INSERT INTO `fivenet_roles` (`id`, `created_at`, `job`, `grade`) VALUES (1,'2023-05-16 14:37:57.379','__default__',1),
(3,'2023-04-03 12:20:11.642','ambulance',1),
(4,'2023-04-03 12:20:28.416','doj',1),
(5,'2023-04-03 12:20:28.433','police',1),
(6,'2023-04-03 12:24:14.457','ambulance',20),
(19,'2023-04-03 16:58:24.661','police',22),
(33,'2023-04-04 16:04:37.578','fib',1),
(54,'2023-04-10 16:37:23.277','ambulance',19),
(93,'2023-04-17 21:20:46.738','fib',19),
(97,'2023-04-19 23:24:25.057','doj',21);
/*!40000 ALTER TABLE `fivenet_roles` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping data for table `fivenet_role_attrs`
--

LOCK TABLES `fivenet_role_attrs` WRITE;
/*!40000 ALTER TABLE `fivenet_role_attrs` DISABLE KEYS */;
INSERT INTO `fivenet_role_attrs` (`role_id`, `created_at`, `updated_at`, `attr_id`, `value`) VALUES (3,'2023-05-16 18:05:22.073','2023-05-16 20:06:25.640',4,'{"stringList":{"strings":["UserProps.Wanted","UserProps.Job"]}}'),
(3,'2023-05-16 18:09:00.879','2023-05-16 20:06:25.649',6,'{"jobList":{"strings":["ambulance","police"]}}'),
(3,'2023-05-16 18:09:02.555','2023-05-16 20:15:20.546',7,'{"jobList":{"strings":["ambulance","police"]}}'),
(6,'2023-05-13 17:06:22.053',NULL,4,'{"stringList":{"strings":["UserProps.Wanted","UserProps.Job"]}}');
/*!40000 ALTER TABLE `fivenet_role_attrs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping data for table `fivenet_role_permissions`
--

LOCK TABLES `fivenet_role_permissions` WRITE;
/*!40000 ALTER TABLE `fivenet_role_permissions` DISABLE KEYS */;
INSERT INTO `fivenet_role_permissions` (`role_id`, `permission_id`, `val`) VALUES (1,1,1),
(1,16,1),
(1,22,1),
(1,23,1),
(1,25,1),
(1,28,1),
(3,1,1),
(3,3,1),
(3,6,0),
(3,8,0),
(3,10,1),
(3,11,1),
(3,16,1),
(3,17,1),
(3,18,1),
(3,19,1),
(3,22,1),
(3,31,1),
(3,7595,0),
(4,1,1),
(4,3,1),
(4,6,1),
(4,8,1),
(4,10,1),
(4,11,1),
(4,16,1),
(4,17,1),
(4,27,1),
(4,31,1),
(4,3056,1),
(5,1,1),
(5,3,1),
(5,6,1),
(5,8,1),
(5,10,1),
(5,11,1),
(5,16,1),
(5,17,1),
(5,27,1),
(5,31,1),
(5,3056,1),
(6,1,1),
(6,3,1),
(6,6,1),
(6,8,1),
(6,10,1),
(6,11,1),
(6,16,1),
(6,17,1),
(6,18,1),
(6,19,1),
(6,20,1),
(6,21,1),
(6,22,1),
(6,23,1),
(6,24,1),
(6,25,1),
(6,26,1),
(6,27,1),
(6,28,1),
(6,29,1),
(6,30,1),
(6,31,1),
(6,74,1),
(6,75,1),
(6,76,1),
(6,77,1),
(6,82,1),
(6,2198,1),
(6,2199,1),
(6,3056,1),
(6,3392,1),
(6,3395,1),
(6,3671,1),
(6,7595,1),
(19,1,1),
(19,3,1),
(19,6,1),
(19,8,1),
(19,10,1),
(19,11,1),
(19,16,1),
(19,17,1),
(19,18,1),
(19,19,1),
(19,20,1),
(19,21,1),
(19,22,1),
(19,23,1),
(19,24,1),
(19,25,1),
(19,26,1),
(19,27,1),
(19,28,1),
(19,29,1),
(19,30,1),
(19,31,1),
(19,74,1),
(19,75,1),
(19,76,1),
(19,77,1),
(19,82,1),
(19,2198,1),
(19,2199,1),
(19,3056,1),
(19,3392,1),
(19,3395,1),
(19,3671,1),
(33,1,1),
(33,3,1),
(33,6,1),
(33,8,1),
(33,10,1),
(33,11,1),
(33,16,1),
(33,17,1),
(33,27,1),
(33,31,1),
(33,3056,1),
(54,1,1),
(54,31,1),
(93,1,1),
(93,3,1),
(93,6,1),
(93,8,1),
(93,10,1),
(93,11,1),
(93,16,1),
(93,17,1),
(93,27,1),
(93,31,1),
(93,3056,1),
(97,1,1),
(97,3,1),
(97,6,1),
(97,10,1),
(97,11,1),
(97,16,1),
(97,18,1),
(97,19,1),
(97,20,1),
(97,22,1),
(97,23,1),
(97,24,1),
(97,25,1),
(97,26,1),
(97,27,1),
(97,31,1),
(97,74,1),
(97,75,1),
(97,76,1),
(97,77,1),
(97,82,1),
(97,2198,1),
(97,2199,1),
(97,3392,1),
(97,3395,1),
(97,3671,1),
(97,5497,1);
/*!40000 ALTER TABLE `fivenet_role_permissions` ENABLE KEYS */;
UNLOCK TABLES;

LOCK TABLES `fivenet_job_permissions` WRITE;
/*!40000 ALTER TABLE `fivenet_job_permissions` DISABLE KEYS */;
INSERT INTO `fivenet_job_permissions` (`job`, `permission_id`, `val`) VALUES ('ambulance',1,1),
('ambulance',3,1),
('ambulance',6,0),
('ambulance',8,0),
('ambulance',10,1),
('ambulance',11,1),
('ambulance',16,1),
('ambulance',17,1),
('ambulance',18,1),
('ambulance',19,1),
('ambulance',22,1),
('ambulance',23,1),
('ambulance',24,1),
('ambulance',25,1),
('ambulance',26,1),
('ambulance',27,1),
('ambulance',28,1),
('ambulance',29,1),
('ambulance',30,1),
('ambulance',31,1),
('ambulance',74,1),
('ambulance',75,1),
('ambulance',76,1),
('ambulance',77,1),
('ambulance',82,1),
('ambulance',7595,0);
/*!40000 ALTER TABLE `fivenet_job_permissions` ENABLE KEYS */;
UNLOCK TABLES;

/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2023-06-16 18:56:18
