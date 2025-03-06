CREATE DATABASE  IF NOT EXISTS `g09-mysql-backend` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;
USE `g09-mysql-backend`;
-- MySQL dump 10.13  Distrib 8.0.36, for Win64 (x86_64)
--
-- Host: 127.0.0.1    Database: g09-mysql-backend
-- ------------------------------------------------------
-- Server version	8.0.37

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `todo_items`
--

DROP TABLE IF EXISTS `todo_items`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `todo_items` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL,
  `title` varchar(150) DEFAULT NULL,
  `description` longtext,
  `image` json DEFAULT NULL,
  `status` enum('Doing','Done','Deleted') DEFAULT 'Doing',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `status_index` (`status`)
) ENGINE=InnoDB AUTO_INCREMENT=19 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `todo_items`
--

LOCK TABLES `todo_items` WRITE;
/*!40000 ALTER TABLE `todo_items` DISABLE KEYS */;
INSERT INTO `todo_items` VALUES (1,0,'This is task 1','Task 1 is always the most challenging one',NULL,'Doing','2025-02-20 22:21:41','2025-02-20 22:21:41'),(3,0,'Task 3','NO GOOD',NULL,'Deleted','2025-02-20 22:25:54','2025-02-25 12:27:55'),(4,0,'Task 4','4th',NULL,'Deleted','2025-02-20 22:25:54','2025-03-01 22:28:31'),(5,1,'Task after begin delete','Way to many',NULL,'Done','2025-02-20 22:25:54','2025-03-04 20:14:32'),(6,2,'Task 6','I cant take it anymore',NULL,'Doing','2025-02-20 22:25:54','2025-03-04 20:14:32'),(7,2,'Task 7','Im going nuts',NULL,'Done','2025-02-20 22:25:54','2025-03-04 20:14:32'),(8,7,'Task 8','It\'s done',NULL,'Done','2025-02-20 22:25:54','2025-03-06 08:16:48'),(9,0,'New task after second updated','After updating, everything has a meaning now !!!!',NULL,'Doing','2025-02-23 15:05:15','2025-02-23 22:46:02'),(10,7,'Clean Architecture task ','Guess what',NULL,'Doing','2025-03-01 10:38:36','2025-03-06 08:16:48'),(11,0,'Clean Architecture task two again','Guess who',NULL,'Doing','2025-03-01 10:43:10','2025-03-01 10:43:10'),(12,4,'Clean Architecture task three again','Cursor AI is sht',NULL,'Doing','2025-03-01 10:43:28','2025-03-04 20:17:11'),(13,0,'Clean Architecture task three again','Cursor AI is sht','{\"id\": 0, \"url\": \"http://localhost:8080/static/1740924885003452200.aws.png\", \"width\": 100, \"height\": 100, \"cloud_name\": \"local\"}','Doing','2025-03-02 21:15:25','2025-03-02 21:15:25'),(14,7,'Clean Architecture task three again','Cursor AI is sht','{\"id\": 0, \"url\": \"http://localhost:8080/static/1740924885003452200.aws.png\", \"width\": 100, \"height\": 100, \"cloud_name\": \"local\"}','Doing','2025-03-03 10:05:56','2025-03-06 08:16:48'),(15,4,'Clean Architecture task three again','Cursor AI is sht','{\"id\": 0, \"url\": \"http://localhost:8080/static/1740924885003452200.aws.png\", \"width\": 100, \"height\": 100, \"cloud_name\": \"local\"}','Doing','2025-03-03 10:06:00','2025-03-04 20:17:11'),(16,7,'Testing out new Auth middleware','This is absurd',NULL,'Deleted','2025-03-04 19:26:22','2025-03-06 08:16:48'),(17,1,'Admin succeed','Wrong table',NULL,'Doing','2025-03-04 19:28:24','2025-03-04 19:59:56'),(18,3,'Testing out new Auth middleware','why is id = 0',NULL,'Doing','2025-03-04 19:39:50','2025-03-04 19:39:50');
/*!40000 ALTER TABLE `todo_items` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `id` int NOT NULL AUTO_INCREMENT,
  `salt` varchar(100) DEFAULT NULL,
  `email` varchar(100) DEFAULT NULL,
  `password` varchar(100) DEFAULT NULL,
  `first_name` varchar(100) DEFAULT NULL,
  `last_name` varchar(100) DEFAULT NULL,
  `phone` varchar(45) DEFAULT NULL,
  `role` enum('user','admin','shipper','mod') DEFAULT 'user',
  `status` int DEFAULT '1',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_email` (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (1,'gWzlbrFwLRPKGIUNsAWIcKNbUwQRutSzPBsijNXZUyHMlfWoUt','hmp@gmail.com','92fc54c9381ab8e4e16195292a26b878','Hoang','Phuc',NULL,'user',1,'2025-03-03 14:32:12','2025-03-04 11:36:31'),(2,'XuJihBIjhElTqIhKHNKnoTXlCnEHWTepcIBiRZvJmlKEOhvlJU','','1e5e86716939ab8a5606e2ac83124998','Dao','Hung',NULL,'user',1,'2025-03-04 11:36:40','2025-03-04 11:36:40'),(3,'jycirnTZhXGthGbhVNoafGJCOsznBLQheWRiorJjpxIgmbUezf','hmmsalp@gmail.com','5411089a279e4c663403aca3b3184067','Dadddo','Hung',NULL,'user',1,'2025-03-04 11:54:26','2025-03-04 11:54:26'),(4,'wSRRIokRTQKBbuVPQEGVqQeflfCkndVWTPyBGBGogsibNChsbC','admin@gmail.com','f4b9c0f55abea76117ba3b4b151874f7','Hoang','Phuc Minh',NULL,'admin',1,'2025-03-04 12:52:53','2025-03-04 12:59:39'),(5,'QyPvZynvTarIUNusMZkehExhChCJVBeqIFJhQxhJpcdWYGColZ','iamid@gmail.com','$2a$12$chSwNYor59h/BcDK96hZ2uB8q7SwPoGMGvQ55cInINKxMIdlnSXh2','Hoang','MVLMFL',NULL,'user',1,'2025-03-05 09:57:21','2025-03-05 09:57:21'),(6,'xpSKLzdWLaRodulCoFxzJWCFigRqDafVdSVzitqbBniTemuOez','mmml@gmail.com','$2a$12$z6sGWMNp9cPxhaV8UWY8L.REMDINHtA18NwQ4pWzvmkguroKGtt.y','Hoang','FL',NULL,'user',1,'2025-03-05 10:04:20','2025-03-05 10:04:20'),(7,'zYMWOhRBvQoqYEStfVSnOgmaNwHOvGCwoZIzRZxcvvKJkAllYW','user1@gmail.com','$2a$12$t5.IhLhBGaBAEfKWttiw4OOObRfLiuEJGebLtvZqdo9jvSiYWoOAC','Hoang','FL',NULL,'user',1,'2025-03-05 10:07:39','2025-03-05 10:07:39'),(8,'VXRZlKQburXIJOWAkbIhONiyUlmmggJuQOPSXcUWhGigKPldfF','user2@gmail.com','$2a$12$SK534dQ42Rj8/SpxcXEXmuEtgyCVBN9eOLFsN211wIoWmRAnTgT/S','Hoang','FL',NULL,'user',1,'2025-03-05 12:18:28','2025-03-05 12:18:28');
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2025-03-06 15:32:51
