-- MySQL dump 10.13  Distrib 5.5.42, for Linux (x86_64)
--
-- Host: localhost    Database: osinstall
-- ------------------------------------------------------
-- Server version	5.5.42

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `devices`
--

DROP TABLE IF EXISTS `devices`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `devices` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `batch_number` varchar(255) NOT NULL,
  `sn` varchar(255) NOT NULL,
  `hostname` varchar(255) NOT NULL,
  `ip` varchar(255) NOT NULL,
  `network_id` int(10) unsigned NOT NULL,
  `os_id` int(10) unsigned NOT NULL,
  `hardware_id` int(10) unsigned DEFAULT NULL,
  `system_id` int(10) unsigned NOT NULL,
  `location` varchar(255) NOT NULL,
  `location_id` int(11) NOT NULL,
  `asset_number` varchar(255) DEFAULT NULL,
  `status` varchar(255) NOT NULL,
  `install_progress` decimal(11,4) DEFAULT '0.0000',
  `install_log` text,
  PRIMARY KEY (`id`),
  UNIQUE KEY `sn` (`sn`),
  UNIQUE KEY `ip` (`ip`),
  KEY `batch_number` (`batch_number`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `devices`
--

LOCK TABLES `devices` WRITE;
/*!40000 ALTER TABLE `devices` DISABLE KEYS */;
INSERT INTO `devices` VALUES (1,'2015-11-27 08:41:08','2015-11-27 08:41:15',NULL,'20151127001','11111111111','idocs-db-01','192.168.1.1',1,1,1,1,'',1,NULL,'pre_install',0.0000,NULL),(2,'2015-11-28 02:41:48','2015-11-28 02:43:46',NULL,'20151128001','2','3','192.168.1.3',2,1,1,8,'',1,'','pre_install',0.0000,''),(4,'2015-11-28 02:48:05','2015-11-28 02:48:05',NULL,'20151128002','4','4','192.168.1.4',2,3,1,1,'',0,'','pre_install',0.0000,''),(5,'2015-11-28 02:48:58','2015-11-28 02:48:58',NULL,'20151128003','5','5','192.168.1.5',2,1,2,1,'',0,'','pre_install',0.0000,''),(6,'2015-11-28 02:53:41','2015-11-28 03:24:43',NULL,'20151128005','6','6','192.168.1.6',2,1,1,1,'',0,'6','pre_install',0.0000,'');
/*!40000 ALTER TABLE `devices` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `hardwares`
--

DROP TABLE IF EXISTS `hardwares`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `hardwares` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `company` varchar(255) NOT NULL,
  `product` varchar(255) NOT NULL,
  `model_name` varchar(255) NOT NULL,
  `raid` varchar(255) DEFAULT NULL,
  `oob` varchar(255) DEFAULT NULL,
  `bios` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `hardwares`
--

LOCK TABLES `hardwares` WRITE;
/*!40000 ALTER TABLE `hardwares` DISABLE KEYS */;
INSERT INTO `hardwares` VALUES (1,'2015-11-20 11:41:50','2015-11-20 11:41:52',NULL,'Dell','PowerEdge','R420','/usr/yunji/osinstall/dell/raid.sh -c=on -d=off -e=on','/usr/yunji/osinstall/dell/oob.sh -a=on -b=off -g=on','/usr/yunji/osinstall/dell/bios.sh -a=on -f=off -m=on'),(2,'2015-11-27 03:11:04','2015-11-27 03:11:06',NULL,'Dell','PowerEdge','R100','raid','oob','bios');
/*!40000 ALTER TABLE `hardwares` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ips`
--

DROP TABLE IF EXISTS `ips`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ips` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `network_id` int(10) unsigned NOT NULL,
  `ip` varchar(255) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=347 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ips`
--

LOCK TABLES `ips` WRITE;
/*!40000 ALTER TABLE `ips` DISABLE KEYS */;
INSERT INTO `ips` VALUES (255,'2015-11-26 08:17:36','2015-11-26 08:17:36',NULL,1,'192.168.100.1'),(256,'2015-11-26 08:17:36','2015-11-26 08:17:36',NULL,1,'192.168.100.2'),(257,'2015-11-26 08:17:36','2015-11-26 08:17:36',NULL,1,'192.168.100.3'),(258,'2015-11-26 08:17:36','2015-11-26 08:17:36',NULL,1,'192.168.100.4'),(259,'2015-11-26 08:17:36','2015-11-26 08:17:36',NULL,1,'192.168.100.5'),(260,'2015-11-26 08:17:36','2015-11-26 08:17:36',NULL,1,'192.168.100.6'),(261,'2015-11-26 08:17:36','2015-11-26 08:17:36',NULL,1,'192.168.100.7'),(262,'2015-11-26 08:17:36','2015-11-26 08:17:36',NULL,1,'192.168.100.8'),(263,'2015-11-26 08:17:36','2015-11-26 08:17:36',NULL,1,'192.168.100.9'),(264,'2015-11-26 08:17:36','2015-11-26 08:17:36',NULL,1,'192.168.100.10'),(265,'2015-11-26 08:17:36','2015-11-26 08:17:36',NULL,1,'192.168.100.11'),(266,'2015-11-26 08:17:36','2015-11-26 08:17:36',NULL,1,'192.168.100.12'),(267,'2015-11-26 08:17:36','2015-11-26 08:17:36',NULL,1,'192.168.100.13'),(268,'2015-11-26 08:17:36','2015-11-26 08:17:36',NULL,1,'192.168.100.14'),(269,'2015-11-26 08:17:37','2015-11-26 08:17:37',NULL,1,'192.168.100.15'),(270,'2015-11-26 08:17:37','2015-11-26 08:17:37',NULL,1,'192.168.100.16'),(271,'2015-11-26 08:17:37','2015-11-26 08:17:37',NULL,1,'192.168.100.17'),(272,'2015-11-26 08:17:37','2015-11-26 08:17:37',NULL,1,'192.168.100.18'),(273,'2015-11-26 08:17:37','2015-11-26 08:17:37',NULL,1,'192.168.100.19'),(274,'2015-11-26 08:17:37','2015-11-26 08:17:37',NULL,1,'192.168.100.20'),(275,'2015-11-26 08:17:37','2015-11-26 08:17:37',NULL,1,'192.168.100.21'),(276,'2015-11-26 08:17:37','2015-11-26 08:17:37',NULL,1,'192.168.100.22'),(277,'2015-11-26 08:17:37','2015-11-26 08:17:37',NULL,1,'192.168.100.23'),(278,'2015-11-26 08:17:37','2015-11-26 08:17:37',NULL,1,'192.168.100.24'),(279,'2015-11-26 08:17:37','2015-11-26 08:17:37',NULL,1,'192.168.100.25'),(280,'2015-11-26 08:17:37','2015-11-26 08:17:37',NULL,1,'192.168.100.26'),(281,'2015-11-26 08:17:37','2015-11-26 08:17:37',NULL,1,'192.168.100.27'),(282,'2015-11-26 08:17:37','2015-11-26 08:17:37',NULL,1,'192.168.100.28'),(283,'2015-11-26 08:17:37','2015-11-26 08:17:37',NULL,1,'192.168.100.29'),(284,'2015-11-26 08:17:37','2015-11-26 08:17:37',NULL,1,'192.168.100.30'),(285,'2015-11-26 08:17:37','2015-11-26 08:17:37',NULL,1,'192.168.100.31'),(286,'2015-11-26 08:17:37','2015-11-26 08:17:37',NULL,1,'192.168.100.32'),(287,'2015-11-26 08:17:37','2015-11-26 08:17:37',NULL,1,'192.168.100.33'),(288,'2015-11-26 08:17:37','2015-11-26 08:17:37',NULL,1,'192.168.100.34'),(289,'2015-11-26 08:17:37','2015-11-26 08:17:37',NULL,1,'192.168.100.35'),(290,'2015-11-26 08:17:37','2015-11-26 08:17:37',NULL,1,'192.168.100.36'),(291,'2015-11-26 08:17:37','2015-11-26 08:17:37',NULL,1,'192.168.100.37'),(292,'2015-11-26 08:17:37','2015-11-26 08:17:37',NULL,1,'192.168.100.38'),(293,'2015-11-26 08:17:37','2015-11-26 08:17:37',NULL,1,'192.168.100.39'),(294,'2015-11-26 08:17:37','2015-11-26 08:17:37',NULL,1,'192.168.100.40'),(295,'2015-11-26 08:17:37','2015-11-26 08:17:37',NULL,1,'192.168.100.41'),(296,'2015-11-26 08:17:37','2015-11-26 08:17:37',NULL,1,'192.168.100.42'),(297,'2015-11-26 08:17:37','2015-11-26 08:17:37',NULL,1,'192.168.100.43'),(298,'2015-11-26 08:17:37','2015-11-26 08:17:37',NULL,1,'192.168.100.44'),(299,'2015-11-26 08:17:37','2015-11-26 08:17:37',NULL,1,'192.168.100.45'),(300,'2015-11-26 08:17:37','2015-11-26 08:17:37',NULL,1,'192.168.100.46'),(301,'2015-11-26 08:17:37','2015-11-26 08:17:37',NULL,1,'192.168.100.47'),(302,'2015-11-26 08:17:37','2015-11-26 08:17:37',NULL,1,'192.168.100.48'),(303,'2015-11-26 08:17:37','2015-11-26 08:17:37',NULL,1,'192.168.100.49'),(304,'2015-11-26 08:17:37','2015-11-26 08:17:37',NULL,1,'192.168.100.50'),(305,'2015-11-26 08:17:37','2015-11-26 08:17:37',NULL,1,'192.168.100.51'),(306,'2015-11-26 08:17:37','2015-11-26 08:17:37',NULL,1,'192.168.100.52'),(307,'2015-11-26 08:17:37','2015-11-26 08:17:37',NULL,1,'192.168.100.53'),(308,'2015-11-26 08:17:37','2015-11-26 08:17:37',NULL,1,'192.168.100.54'),(309,'2015-11-26 08:17:37','2015-11-26 08:17:37',NULL,1,'192.168.100.55'),(310,'2015-11-26 08:17:37','2015-11-26 08:17:37',NULL,1,'192.168.100.56'),(311,'2015-11-26 08:17:37','2015-11-26 08:17:37',NULL,1,'192.168.100.57'),(312,'2015-11-26 08:17:37','2015-11-26 08:17:37',NULL,1,'192.168.100.58'),(313,'2015-11-26 08:17:37','2015-11-26 08:17:37',NULL,1,'192.168.100.59'),(314,'2015-11-26 08:17:37','2015-11-26 08:17:37',NULL,1,'192.168.100.60'),(315,'2015-11-26 08:17:37','2015-11-26 08:17:37',NULL,1,'192.168.100.61'),(316,'2015-11-26 08:17:37','2015-11-26 08:17:37',NULL,1,'192.168.100.62'),(317,'2015-11-26 08:58:54','2015-11-26 08:58:54',NULL,2,'192.168.1.1'),(318,'2015-11-26 08:58:54','2015-11-26 08:58:54',NULL,2,'192.168.1.2'),(319,'2015-11-26 08:58:54','2015-11-26 08:58:54',NULL,2,'192.168.1.3'),(320,'2015-11-26 08:58:54','2015-11-26 08:58:54',NULL,2,'192.168.1.4'),(321,'2015-11-26 08:58:54','2015-11-26 08:58:54',NULL,2,'192.168.1.5'),(322,'2015-11-26 08:58:54','2015-11-26 08:58:54',NULL,2,'192.168.1.6'),(323,'2015-11-26 08:58:54','2015-11-26 08:58:54',NULL,2,'192.168.1.7'),(324,'2015-11-26 08:58:54','2015-11-26 08:58:54',NULL,2,'192.168.1.8'),(325,'2015-11-26 08:58:54','2015-11-26 08:58:54',NULL,2,'192.168.1.9'),(326,'2015-11-26 08:58:54','2015-11-26 08:58:54',NULL,2,'192.168.1.10'),(327,'2015-11-26 08:58:54','2015-11-26 08:58:54',NULL,2,'192.168.1.11'),(328,'2015-11-26 08:58:54','2015-11-26 08:58:54',NULL,2,'192.168.1.12'),(329,'2015-11-26 08:58:54','2015-11-26 08:58:54',NULL,2,'192.168.1.13'),(330,'2015-11-26 08:58:54','2015-11-26 08:58:54',NULL,2,'192.168.1.14'),(331,'2015-11-26 08:58:54','2015-11-26 08:58:54',NULL,2,'192.168.1.15'),(332,'2015-11-26 08:58:54','2015-11-26 08:58:54',NULL,2,'192.168.1.16'),(333,'2015-11-26 08:58:54','2015-11-26 08:58:54',NULL,2,'192.168.1.17'),(334,'2015-11-26 08:58:54','2015-11-26 08:58:54',NULL,2,'192.168.1.18'),(335,'2015-11-26 08:58:54','2015-11-26 08:58:54',NULL,2,'192.168.1.19'),(336,'2015-11-26 08:58:54','2015-11-26 08:58:54',NULL,2,'192.168.1.20'),(337,'2015-11-26 08:58:54','2015-11-26 08:58:54',NULL,2,'192.168.1.21'),(338,'2015-11-26 08:58:54','2015-11-26 08:58:54',NULL,2,'192.168.1.22'),(339,'2015-11-26 08:58:54','2015-11-26 08:58:54',NULL,2,'192.168.1.23'),(340,'2015-11-26 08:58:54','2015-11-26 08:58:54',NULL,2,'192.168.1.24'),(341,'2015-11-26 08:58:54','2015-11-26 08:58:54',NULL,2,'192.168.1.25'),(342,'2015-11-26 08:58:54','2015-11-26 08:58:54',NULL,2,'192.168.1.26'),(343,'2015-11-26 08:58:54','2015-11-26 08:58:54',NULL,2,'192.168.1.27'),(344,'2015-11-26 08:58:54','2015-11-26 08:58:54',NULL,2,'192.168.1.28'),(345,'2015-11-26 08:58:54','2015-11-26 08:58:54',NULL,2,'192.168.1.29'),(346,'2015-11-26 08:58:54','2015-11-26 08:58:54',NULL,2,'192.168.1.30');
/*!40000 ALTER TABLE `ips` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `locations`
--

DROP TABLE IF EXISTS `locations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `locations` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `pid` int(10) unsigned NOT NULL,
  `name` varchar(255) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `locations`
--

LOCK TABLES `locations` WRITE;
/*!40000 ALTER TABLE `locations` DISABLE KEYS */;
INSERT INTO `locations` VALUES (2,'2015-11-25 06:04:56','2015-11-25 06:04:56',NULL,0,'留下机房'),(3,'2015-11-25 06:05:03','2015-11-25 06:05:03',NULL,0,'萧山机房'),(4,'2015-11-25 06:05:08','2015-11-25 06:05:08',NULL,0,'千岛湖机房'),(5,'2015-11-25 08:11:01','2015-11-25 08:11:01',NULL,2,'A1'),(6,'2015-11-25 08:11:06','2015-11-25 08:11:06',NULL,2,'A2'),(7,'2015-11-25 08:11:19','2015-11-25 09:23:37',NULL,5,'1012'),(8,'2015-11-25 08:11:23','2015-11-25 08:11:23',NULL,5,'102'),(10,'2015-11-25 09:18:18','2015-11-25 09:18:18',NULL,3,'301'),(11,'2015-11-26 01:59:35','2015-11-26 02:30:31',NULL,0,'test'),(12,'2015-11-26 02:38:57','2015-11-26 02:38:57',NULL,0,'test1'),(14,'2015-11-28 03:09:08','2015-11-28 03:09:08',NULL,2,'A3');
/*!40000 ALTER TABLE `locations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `networks`
--

DROP TABLE IF EXISTS `networks`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `networks` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `network` varchar(255) NOT NULL,
  `netmask` varchar(255) NOT NULL,
  `gateway` varchar(255) NOT NULL,
  `vlan` varchar(255) DEFAULT NULL,
  `trunk` varchar(255) DEFAULT NULL,
  `bonding` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `network` (`network`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `networks`
--

LOCK TABLES `networks` WRITE;
/*!40000 ALTER TABLE `networks` DISABLE KEYS */;
INSERT INTO `networks` VALUES (1,'2015-11-26 08:16:20','2015-11-26 08:17:36',NULL,'192.168.100.1/26','255.255.255.0','192.168.1.100','vlan','trunk','bonding'),(2,'2015-11-26 08:58:54','2015-11-26 08:58:54',NULL,'192.168.1.1/27','255.255.255.224','192.168.1.100','','','');
/*!40000 ALTER TABLE `networks` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `os_configs`
--

DROP TABLE IF EXISTS `os_configs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `os_configs` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `name` varchar(255) NOT NULL,
  `pxe` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=25 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `os_configs`
--

LOCK TABLES `os_configs` WRITE;
/*!40000 ALTER TABLE `os_configs` DISABLE KEYS */;
INSERT INTO `os_configs` VALUES (1,'2015-11-24 07:54:05','2015-11-24 07:54:05',NULL,'centos5u7-x86_64','PXE配置\nPXE配置\nPXE配置'),(2,'2015-11-24 08:00:44','2015-11-24 08:00:44',NULL,'centos6u7-x86_64','PXE配置\nPXE配置\nPXE配置'),(3,'2015-11-24 08:00:58','2015-11-24 08:00:58',NULL,'centos7u1-x86_64','PXE配置\nPXE配置\nPXE配置'),(4,'2015-11-24 08:03:59','2015-11-26 03:04:46',NULL,'centos6u1-x86_64','PXE配置\nPXE配置'),(5,'2015-11-24 08:04:53','2015-11-24 10:47:52','2015-11-24 10:58:17','centos6u2-x86_641','PXE1');
/*!40000 ALTER TABLE `os_configs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `system_configs`
--

DROP TABLE IF EXISTS `system_configs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `system_configs` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `name` varchar(255) NOT NULL,
  `content` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `system_configs`
--

LOCK TABLES `system_configs` WRITE;
/*!40000 ALTER TABLE `system_configs` DISABLE KEYS */;
INSERT INTO `system_configs` VALUES (1,'2015-11-25 02:24:10','2015-11-25 02:31:03',NULL,'系统模板1','内容1\n内容1'),(8,'2015-11-26 08:58:22','2015-11-26 08:58:22',NULL,'系统模板2','内容2\n内容2');
/*!40000 ALTER TABLE `system_configs` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2015-11-28 11:27:58
