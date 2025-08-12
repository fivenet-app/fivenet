-- MariaDB dump 10.19  Distrib 10.11.2-MariaDB, for Linux (x86_64)
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
-- Dumping data for table `fivenet_documents`
--

LOCK TABLES `fivenet_documents` WRITE;
/*!40000 ALTER TABLE `fivenet_documents` DISABLE KEYS */;
INSERT INTO `fivenet_documents` (`id`, `created_at`, `updated_at`, `deleted_at`, `category_id`, `title`, `content_type`, `content`, `data`, `creator_id`, `creator_job`, `state`, `closed`, `public`) VALUES (1,'2023-03-17 19:57:09.898','2023-03-17 18:57:16.587',NULL,NULL,'Public Document without category',0,'I\'m a public Document without a category.',NULL,1,'ambulance','Open',NULL,1),
(2,'2023-03-17 19:57:13.244','2023-03-17 18:57:16.596',NULL,4,'Public Document with category (Closed State)',0,'I\'m a public Document with a category that is closed.',NULL,1,'ambulance','Closed',1,1),
(3,'2023-03-17 18:54:44.115','2023-03-17 18:57:04.438',NULL,1,'Patientenakte Thomas G.',0,'Only for Ambulance.',NULL,2,'ambulance','Open',0,0),
(4,'2023-03-17 18:57:04.391',NULL,NULL,1,'Bloodresults for DOJ Case',0,'Only for DOJ, Ambulance and the patient.',NULL,2,'ambulance','Open',0,0),
(5,'2023-03-17 18:57:04.413',NULL,NULL,1,'Drugtest for DOJ Case',0,'Only for PD, DOJ and Ambulance.',NULL,2,'ambulance','Open',0,0),
(6,'2023-03-17 18:57:55.203',NULL,NULL,2,'Police document about a criminal investigation',0,'Only for PD.',NULL,3,'police','Open',0,0),
(7,'2023-03-17 18:58:53.956','2023-03-17 18:59:51.616',NULL,NULL,'DOJ Request for medical bloodtests',0,'Only for DOJ and Ambulance.',NULL,3,'doj','Closed',0,0),
(8,'2023-03-17 19:28:38.145',NULL,NULL,4,'Internal Ambulance Doc',0,'Internal Ambulance Doc',NULL,2,'ambulance','Open',0,0),
(9,'2023-03-17 19:28:38.155',NULL,NULL,4,'Internal Ambulance Doc only grade 17 and higher',0,'Internal Ambulance Doc only grade 17 and higher',NULL,1,'ambulance','Open',0,0);
/*!40000 ALTER TABLE `fivenet_documents` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping data for table `fivenet_documents_categories`
--

LOCK TABLES `fivenet_documents_categories` WRITE;
/*!40000 ALTER TABLE `fivenet_documents_categories` DISABLE KEYS */;
INSERT INTO `fivenet_documents_categories` (`id`, `name`, `description`, `job`) VALUES (1,'Patient / File','Patient files (e.g., reports, results)','ambulance'),
(2,'Criminal Record','Criminal record of a citizen','police'),
(3,'Non-Existant','Document Category for a non-existent job, no person should see it.','non-existent'),
(4,'Patient / Unused','Unused category for testing','ambulance');
/*!40000 ALTER TABLE `fivenet_documents_categories` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping data for table `fivenet_documents_comments`
--

LOCK TABLES `fivenet_documents_comments` WRITE;
/*!40000 ALTER TABLE `fivenet_documents_comments` DISABLE KEYS */;
INSERT INTO `fivenet_documents_comments` (`id`, `created_at`, `updated_at`, `deleted_at`, `document_id`, `comment`, `creator_id`) VALUES (1,'2023-03-17 19:34:30.052',NULL,NULL,1,'Hello World!',2);
/*!40000 ALTER TABLE `fivenet_documents_comments` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping data for table `fivenet_documents_access`
--

LOCK TABLES `fivenet_documents_access` WRITE;
/*!40000 ALTER TABLE `fivenet_documents_access` DISABLE KEYS */;
INSERT INTO `fivenet_documents_access` (`id`, `target_id`, `job`, `minimum_grade`, `access`) VALUES (1,3,'ambulance',0,3),
(2,4,'ambulance',0,3),
(3,4,'doj',0,3),
(4,5,'ambulance',0,3),
(5,5,'doj',0,3),
(6,6,'police',0,3),
(7,7,'doj',0,3),
(8,7,'ambulance',0,2),
(9,8,'ambulance',0,3),
(10,9,'ambulance',17,3);
INSERT INTO `fivenet_documents_access` (`id`, `target_id`, `user_id`, `access`) VALUES (11,4,5,1);
/*!40000 ALTER TABLE `fivenet_documents_access` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping data for table `fivenet_documents_references`
--

LOCK TABLES `fivenet_documents_references` WRITE;
/*!40000 ALTER TABLE `fivenet_documents_references` DISABLE KEYS */;
INSERT INTO `fivenet_documents_references` (`id`, `created_at`, `deleted_at`, `source_document_id`, `reference`, `target_document_id`, `creator_id`) VALUES (1,'2023-03-17 18:59:57.652',NULL,5,1,7,1),
(2,'2023-03-18 17:34:35.564',NULL,5,0,5,1);
/*!40000 ALTER TABLE `fivenet_documents_references` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping data for table `fivenet_documents_relations`
--

LOCK TABLES `fivenet_documents_relations` WRITE;
/*!40000 ALTER TABLE `fivenet_documents_relations` DISABLE KEYS */;
INSERT INTO `fivenet_documents_relations` (`id`, `created_at`, `deleted_at`, `document_id`, `source_user_id`, `relation`, `target_user_id`) VALUES (1,'2023-03-17 19:02:48.621',NULL,5,1,0,5),
(2,'2023-03-17 19:03:12.428',NULL,6,4,1,5);
/*!40000 ALTER TABLE `fivenet_documents_relations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping data for table `fivenet_documents_templates`
--

LOCK TABLES `fivenet_documents_templates` WRITE;
/*!40000 ALTER TABLE `fivenet_documents_templates` DISABLE KEYS */;
INSERT INTO `fivenet_documents_templates` (`id`, `created_at`, `updated_at`, `deleted_at`, `weight`, `category_id`, `title`, `description`, `content_title`, `content`, `access`, `schema`, `creator_job`) VALUES (1,'2023-03-17 19:31:28.661','2023-03-17 19:31:35.698',NULL,10,1,'Patientenakte','LSMD Patientenakten Template','Patientenakte NAME','Patientenakte für Name','{}',NULL,'ambulance');
/*!40000 ALTER TABLE `fivenet_documents_templates` ENABLE KEYS */;
UNLOCK TABLES;

/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2023-06-03 17:00:36
