/*M!999999\- enable the sandbox mode */ 
-- MariaDB dump 10.19-11.7.2-MariaDB, for Win64 (AMD64)
--
-- Host: localhost    Database: vireo_gin_admin
-- ------------------------------------------------------
-- Server version	8.0.12

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*M!100616 SET @OLD_NOTE_VERBOSITY=@@NOTE_VERBOSITY, NOTE_VERBOSITY=0 */;

--
-- Table structure for table `config`
--

DROP TABLE IF EXISTS `config`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `config` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `config_name` varchar(100) CHARACTER SET utf8 COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '配置名称',
  `config_key` varchar(100) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL COMMENT '配置键',
  `config_value` varchar(500) CHARACTER SET utf8 COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '配置值',
  `remark` varchar(500) CHARACTER SET utf8 COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '描述备注',
  `creator_id` int(11) DEFAULT NULL COMMENT '创建人ID',
  `dept_id` int(11) DEFAULT NULL COMMENT '部门ID',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_config_key` (`config_key`),
  KEY `idx_creator_id` (`creator_id`),
  KEY `idx_dept_id` (`dept_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=38 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='系统配置表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `config`
--

LOCK TABLES `config` WRITE;
/*!40000 ALTER TABLE `config` DISABLE KEYS */;
INSERT INTO `config` VALUES
(1,'网站名称','WEB_NAME','Vireo后台管理系统','网站名称',1,1,'2025-06-12 23:44:30','2025-06-13 10:17:19',NULL);
/*!40000 ALTER TABLE `config` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `depts`
--

DROP TABLE IF EXISTS `depts`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `depts` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '部门ID',
  `parent_id` bigint(20) unsigned DEFAULT '0' COMMENT '父部门ID',
  `name` varchar(50) NOT NULL COMMENT '部门名称',
  `code` varchar(50) NOT NULL COMMENT '部门编码',
  `sort` int(11) NOT NULL DEFAULT '0' COMMENT '显示排序',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态(0:禁用 1:启用)',
  `leader` varchar(20) DEFAULT NULL COMMENT '负责人',
  `phone` varchar(20) DEFAULT NULL COMMENT '联系电话',
  `email` varchar(50) DEFAULT NULL COMMENT '邮箱',
  `dept_id` int(11) DEFAULT NULL,
  `creator_id` int(11) DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='部门表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `depts`
--

LOCK TABLES `depts` WRITE;
/*!40000 ALTER TABLE `depts` DISABLE KEYS */;
INSERT INTO `depts` VALUES
(1,0,'总部','ZONGBU',0,1,NULL,NULL,NULL,10,1,'2025-06-06 08:06:07','2025-06-12 06:43:42',NULL),
(10,1,'市场部','SHICHANGBU',0,1,NULL,NULL,NULL,10,1,'2025-06-06 09:46:09','2025-06-12 08:30:42',NULL),
(11,0,'石家庄分公司','SHIJIAZHUANG',1,1,NULL,NULL,NULL,10,1,'2025-06-06 12:19:48','2025-06-12 08:35:24',NULL);
/*!40000 ALTER TABLE `depts` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `dict_item`
--

DROP TABLE IF EXISTS `dict_item`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dict_item` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `dict_code` varchar(50) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL,
  `value` varchar(50) COLLATE utf8_unicode_ci NOT NULL,
  `label` varchar(200) COLLATE utf8_unicode_ci DEFAULT NULL,
  `tag_type` varchar(200) COLLATE utf8_unicode_ci DEFAULT NULL,
  `status` tinyint(4) DEFAULT '1',
  `sort` int(11) DEFAULT '0',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=14 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `dict_item`
--

LOCK TABLES `dict_item` WRITE;
/*!40000 ALTER TABLE `dict_item` DISABLE KEYS */;
INSERT INTO `dict_item` VALUES
(2,'gender','1','男','primary',1,1,'2025-06-09 22:24:13','2025-06-09 22:24:13',NULL),
(3,'gender','2','女','primary',1,2,'2025-06-09 22:34:20','2025-06-09 22:34:20',NULL),
(4,'gender','0','未知','primary',1,3,'2025-06-09 22:34:34','2025-06-09 22:34:34',NULL),
(6,'notice_type','1','系统升级','primary',1,1,'2025-06-13 15:01:45','2025-06-13 15:01:45',NULL),
(7,'notice_type','2','系统维护','danger',1,2,'2025-06-13 15:02:01','2025-06-13 15:02:01',NULL),
(8,'notice_type','3','安全警告','warning',1,3,'2025-06-13 15:02:24','2025-06-13 15:02:24',NULL),
(9,'notice_type','4','假期通知','success',1,4,'2025-06-13 15:02:41','2025-06-13 15:02:41',NULL),
(10,'notice_type','5','公司新闻','success',1,5,'2025-06-13 15:02:58','2025-06-13 15:02:58',NULL),
(11,'notice_level','L','低','success',1,1,'2025-06-13 15:04:17','2025-06-13 15:04:17',NULL),
(12,'notice_level','M','中','warning',1,2,'2025-06-13 15:04:29','2025-06-13 15:04:29',NULL),
(13,'notice_level','H','高','danger',1,3,'2025-06-13 15:04:41','2025-06-13 15:04:41',NULL);
/*!40000 ALTER TABLE `dict_item` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `dict_type`
--

DROP TABLE IF EXISTS `dict_type`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dict_type` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `dict_code` varchar(50) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL,
  `name` varchar(50) COLLATE utf8_unicode_ci NOT NULL,
  `remark` varchar(200) COLLATE utf8_unicode_ci DEFAULT NULL,
  `status` tinyint(4) DEFAULT '1',
  `sort` int(11) DEFAULT '0',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `code` (`dict_code`)
) ENGINE=MyISAM AUTO_INCREMENT=6 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `dict_type`
--

LOCK TABLES `dict_type` WRITE;
/*!40000 ALTER TABLE `dict_type` DISABLE KEYS */;
INSERT INTO `dict_type` VALUES
(3,'1','1','性别',1,0,'2025-06-09 19:33:11','2025-06-09 21:02:37','2025-06-09 21:44:57'),
(2,'gender','性别','性别',1,0,'2025-06-09 19:33:11','2025-06-10 13:49:14',NULL),
(4,'notice_type','通知类型','',1,0,'2025-06-13 15:00:50','2025-06-13 15:00:50',NULL),
(5,'notice_level','通知等级','',1,0,'2025-06-13 15:03:51','2025-06-13 15:03:51',NULL);
/*!40000 ALTER TABLE `dict_type` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `menu`
--

DROP TABLE IF EXISTS `menu`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `menu` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `parent_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '父菜单ID (0表示根菜单)',
  `title` varchar(32) NOT NULL COMMENT '菜单名称',
  `name` varchar(32) DEFAULT NULL COMMENT '前端路由名称 (需唯一)',
  `path` varchar(128) DEFAULT NULL COMMENT '前端路由路径',
  `component` varchar(128) DEFAULT NULL COMMENT '前端组件路径',
  `icon` varchar(64) DEFAULT NULL COMMENT '图标标识',
  `sort` int(11) NOT NULL DEFAULT '0' COMMENT '排序值 (越小越靠前)',
  `visible` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否隐藏 (1=显示, 0=隐藏)',
  `always_show` tinyint(1) NOT NULL DEFAULT '0',
  `perm` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '权限标识 (type=2时必填)，根据权限控制按钮的显示',
  `params` varchar(255) DEFAULT NULL COMMENT '路由参数',
  `redirect` varchar(128) DEFAULT NULL COMMENT '重定向路径',
  `keep_alive` tinyint(1) DEFAULT '0' COMMENT '是否缓存页面 (0=否, 1=是)',
  `type` tinyint(1) NOT NULL COMMENT '菜单类型 (1=菜单, 2=按钮/权限)',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_parent_id` (`parent_id`)
) ENGINE=InnoDB AUTO_INCREMENT=25 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='系统菜单表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `menu`
--

LOCK TABLES `menu` WRITE;
/*!40000 ALTER TABLE `menu` DISABLE KEYS */;
INSERT INTO `menu` VALUES
(1,0,'系统管理','System','/system','Layout','el-icon-setting',1,1,0,'','','/system/user',1,1,'2025-06-03 17:03:51','2025-06-10 22:45:19',NULL),
(2,1,'用户管理','User','user','system/user/index','el-icon-user',10,1,0,NULL,NULL,NULL,1,1,'2025-06-03 17:03:51','2025-06-09 21:42:36',NULL),
(4,1,'菜单管理','Menu','menu','system/menu/index','el-icon-menu',30,1,0,NULL,NULL,'',1,1,'2025-06-03 17:03:51','2025-06-09 21:42:41',NULL),
(5,1,'角色管理','Role','role','system/role/index','role',40,1,0,'',NULL,'',1,1,'2025-06-04 15:39:24','2025-06-05 22:22:18',NULL),
(6,5,'角色查询','','','','',1,1,0,'sys:role:query','','',1,4,'2025-06-04 17:12:32','2025-06-10 16:10:20','2025-06-05 15:48:43'),
(20,1,'部门管理','Dept','dept','system/dept/index','menu',50,1,0,'','','',0,1,'2025-06-06 12:53:05','2025-06-06 12:53:33',NULL),
(21,1,'字典管理','Dict','dict','system/dict/index','dict',60,1,0,'','','',1,1,'2025-06-07 13:23:22','2025-06-09 22:35:39',NULL),
(22,1,'字典数据','DictItem','dict-item','system/dict/dict-item','document',61,0,0,'','','',1,1,'2025-06-09 21:38:07','2025-06-10 09:30:37',NULL),
(23,1,'系统配置','Config','config','system/config/index','system',70,1,0,'','','',1,1,'2025-06-12 23:06:25','2025-06-12 23:06:25',NULL),
(24,1,'通知公告','Notice','notice','system/notice/index','',80,1,0,'','','',1,1,'2025-06-13 14:49:54','2025-06-13 15:05:09',NULL);
/*!40000 ALTER TABLE `menu` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `migration_logs`
--

DROP TABLE IF EXISTS `migration_logs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `migration_logs` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `name` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=MyISAM AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `migration_logs`
--

LOCK TABLES `migration_logs` WRITE;
/*!40000 ALTER TABLE `migration_logs` DISABLE KEYS */;
INSERT INTO `migration_logs` VALUES
(1,'2025-05-25 13:00:32','20231010_rbac_init');
/*!40000 ALTER TABLE `migration_logs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `notice_receiver`
--

DROP TABLE IF EXISTS `notice_receiver`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `notice_receiver` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `notice_id` bigint(20) NOT NULL COMMENT '通知ID',
  `user_id` bigint(20) NOT NULL COMMENT '接收人ID',
  `user_name` varchar(32) DEFAULT NULL COMMENT '接收人姓名',
  `is_read` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否已读(0-未读 1-已读)',
  `read_time` datetime DEFAULT NULL COMMENT '阅读时间',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_notice_user` (`notice_id`,`user_id`) COMMENT '通知用户唯一索引',
  KEY `idx_user_read` (`user_id`,`is_read`) COMMENT '用户阅读状态索引'
) ENGINE=InnoDB AUTO_INCREMENT=39 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='通知接收人表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `notice_receiver`
--

LOCK TABLES `notice_receiver` WRITE;
/*!40000 ALTER TABLE `notice_receiver` DISABLE KEYS */;
INSERT INTO `notice_receiver` VALUES
(16,3,4,NULL,0,NULL,'2025-06-16 07:29:18','2025-06-16 13:51:26',NULL),
(17,4,4,NULL,0,NULL,'2025-06-16 08:12:44',NULL,NULL),
(18,7,4,NULL,0,NULL,'2025-06-16 08:12:55',NULL,NULL),
(19,7,5,NULL,1,'2025-06-16 14:28:46','2025-06-16 08:12:55',NULL,NULL),
(20,7,1,NULL,1,'2025-06-16 12:49:09','2025-06-16 08:12:55',NULL,NULL),
(21,6,1,NULL,1,'2025-06-16 12:48:46','2025-06-16 08:18:35',NULL,NULL),
(22,4,5,NULL,1,'2025-06-16 14:28:35','2025-06-16 11:15:43',NULL,NULL),
(23,4,1,NULL,1,'2025-06-16 12:45:24','2025-06-16 11:15:43',NULL,NULL),
(25,3,1,NULL,0,NULL,'2025-06-16 11:15:54','2025-06-16 13:51:26',NULL),
(26,3,5,NULL,1,'2025-06-16 14:28:50','2025-06-16 11:15:54','2025-06-16 13:51:26',NULL),
(28,31,5,NULL,1,'2025-06-16 14:28:50','2025-06-16 11:16:03',NULL,NULL),
(29,31,1,NULL,1,'2025-06-16 13:07:07','2025-06-16 11:16:03',NULL,NULL),
(30,31,4,NULL,0,NULL,'2025-06-16 11:16:03',NULL,NULL),
(31,5,2,NULL,0,NULL,'2025-06-16 12:50:07',NULL,NULL),
(32,5,3,NULL,0,NULL,'2025-06-16 12:50:07',NULL,NULL);
/*!40000 ALTER TABLE `notice_receiver` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `notices`
--

DROP TABLE IF EXISTS `notices`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `notices` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '通知ID',
  `title` varchar(50) NOT NULL COMMENT '通知标题',
  `content` text NOT NULL COMMENT '通知内容',
  `type` varchar(2) DEFAULT NULL COMMENT '通知类型(1-公告 2-通知 3-提醒)',
  `level` char(1) DEFAULT 'M' COMMENT '优先级(L-低 M-中 H-高)',
  `target_type` tinyint(4) DEFAULT '1' COMMENT '目标类型(1-全体 2-指定部门 3-指定角色 4-指定用户)',
  `target_ids` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `status` tinyint(4) DEFAULT '0' COMMENT '状态(0-草稿 1-发布 2-撤回)',
  `creator_id` bigint(20) DEFAULT NULL COMMENT '创建人ID',
  `dept_id` bigint(20) DEFAULT NULL COMMENT '创建人ID',
  `is_read` tinyint(1) DEFAULT NULL,
  `published_at` datetime DEFAULT NULL COMMENT '发布时间',
  `revoked_at` datetime DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_type_status` (`type`,`status`) COMMENT '类型状态联合索引',
  KEY `idx_creator` (`creator_id`) COMMENT '创建人索引'
) ENGINE=InnoDB AUTO_INCREMENT=33 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='通知公告表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `notices`
--

LOCK TABLES `notices` WRITE;
/*!40000 ALTER TABLE `notices` DISABLE KEYS */;
INSERT INTO `notices` VALUES
(3,'测试','<p>少时诵诗书少时诵诗书是撒是撒是撒是撒+</p>','1','L',1,NULL,2,1,1,0,'2025-06-16 13:07:21','2025-06-16 13:51:44','2025-06-13 18:29:51','2025-06-16 13:51:44',NULL),
(4,'22','<p>2222222222222</p>','1','M',1,NULL,1,1,1,0,'2025-06-16 11:15:43','2025-06-16 11:15:41','2025-06-13 19:25:25','2025-06-16 11:15:43',NULL),
(5,'我问问','<p>请求</p>','1','L',2,'[2,3]',1,1,1,0,'2025-06-16 12:50:07','2025-06-16 12:50:06','2025-06-13 23:36:04','2025-06-16 12:50:07',NULL),
(6,'我委屈','<p>呜呜呜呜</p>','1','L',2,'[1]',1,1,1,NULL,'2025-06-16 08:18:35',NULL,'2025-06-13 23:46:43','2025-06-16 08:18:35',NULL),
(7,'发发发','<p>反反复复</p>','1','M',2,'[1,4,5]',1,1,1,NULL,'2025-06-16 08:12:55','2025-06-15 17:15:38','2025-06-13 23:47:17','2025-06-16 08:12:55',NULL),
(27,'222','<p>222222</p>','1','L',2,'[1,4]',2,0,0,0,'2025-06-15 16:56:46','2025-06-15 20:38:29','2025-06-15 16:54:45','2025-06-15 20:38:29',NULL),
(28,'对对对','<p>222222</p>','1','L',2,'[1,4]',1,0,0,0,'2025-06-15 22:20:50',NULL,'2025-06-15 16:56:13','2025-06-15 22:20:50',NULL),
(31,'2222','<p>2222</p>','1','L',1,NULL,1,0,0,0,'2025-06-16 11:16:03','2025-06-16 11:16:02','2025-06-15 17:33:15','2025-06-16 11:16:03',NULL),
(32,'22','<p>2222</p>','1','L',2,'[1,5]',0,0,0,0,NULL,NULL,'2025-06-15 17:36:51','2025-06-15 17:36:51',NULL);
/*!40000 ALTER TABLE `notices` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `permissions`
--

DROP TABLE IF EXISTS `permissions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `permissions` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `code` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '权限标识',
  `name` varchar(255) NOT NULL COMMENT '权限名称',
  `module` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `description` varchar(500) DEFAULT NULL COMMENT '描述',
  `type` enum('menu','button','api') DEFAULT 'api' COMMENT '权限类型',
  `icon` varchar(100) DEFAULT NULL COMMENT '菜单图标',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_permissions_code` (`code`),
  KEY `idx_permissions_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=70 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='系统权限表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `permissions`
--

LOCK TABLES `permissions` WRITE;
/*!40000 ALTER TABLE `permissions` DISABLE KEYS */;
INSERT INTO `permissions` VALUES
(1,'sys:dict-item:edit','编辑字典项','字典项管理','查看Dict选项','api',NULL,NULL,NULL,NULL),
(2,'sys:role:options','角色下拉列表','角色管理','获取角色下拉列表','api',NULL,NULL,NULL,NULL),
(3,'sys:user:info','用户信息表单','用户管理','获取用户信息','api',NULL,NULL,NULL,NULL),
(4,'sys:dict-item:form','字典项表单','字典项管理','查看Dict选项','api',NULL,NULL,NULL,NULL),
(5,'sys:dept:view','查看部门','部门管理','获取部门详情','api',NULL,NULL,NULL,NULL),
(6,'sys:notice:query','Notices列表','Notices管理','查看Notices列表','api',NULL,NULL,NULL,NULL),
(7,'sys:dict-item:query','字典项查询','字典项管理','查看Dict选项','api',NULL,NULL,NULL,NULL),
(8,'sys:noticereceiver:details','NoticeReceiver详情','NoticeReceiver管理','查看NoticeReceiver详情','api',NULL,NULL,NULL,NULL),
(9,'sys:config:add','新建Config','Config管理','创建Config','api',NULL,NULL,NULL,NULL),
(10,'sys:config:delete','删除Config','Config管理','删除Config','api',NULL,NULL,NULL,NULL),
(11,'sys:role:edit','编辑角色','角色管理','更新角色','api',NULL,NULL,NULL,NULL),
(12,'sys:role:menu','角色菜单列表','角色管理','获取角色菜单','api',NULL,NULL,NULL,NULL),
(13,'sys:config:update','更新Config','Config管理','更新Config','api',NULL,NULL,NULL,NULL),
(14,'sys:dict-item:add','新增字典项','字典项管理','查看Dict选项','api',NULL,NULL,NULL,NULL),
(15,'sys:menu:options','菜单下拉列表','菜单管理','获取菜单下拉选项','api',NULL,NULL,NULL,NULL),
(16,'sys:menu:edit','编辑菜单','菜单管理','编辑菜单','api',NULL,NULL,NULL,NULL),
(17,'sys:noticereceiver:add','新建NoticeReceiver','NoticeReceiver管理','创建NoticeReceiver','api',NULL,NULL,NULL,NULL),
(18,'sys:user:reset-password','重置密码','用户管理','重置用户密码','api',NULL,NULL,NULL,NULL),
(19,'sys:noticereceiver:delete','删除NoticeReceiver','NoticeReceiver管理','删除NoticeReceiver','api',NULL,NULL,NULL,NULL),
(20,'sys:notice:add','新建Notices','Notices管理','创建Notices','api',NULL,NULL,NULL,NULL),
(21,'sys:notice:delete','删除Notices','Notices管理','删除Notices','api',NULL,NULL,NULL,NULL),
(22,'sys:dept:options','部门下拉列表','部门管理','获取部门下拉列表','api',NULL,NULL,NULL,NULL),
(23,'sys:dict-item:details','查看字典项','字典项管理','查看Dict选项','api',NULL,NULL,NULL,NULL),
(24,'sys:role:perm:update','更新角色权限','角色管理','更新角色权限','api',NULL,NULL,NULL,NULL),
(25,'sys:config:view','Config详情','Config管理','查看Config详情','api',NULL,NULL,NULL,NULL),
(26,'sys:config:query','Config列表','Config管理','查看Config列表','api',NULL,NULL,NULL,NULL),
(27,'sys:menu:details','菜单详情','菜单管理','查看菜单详情','api',NULL,NULL,NULL,NULL),
(28,'sys:menu:add','新增菜单','菜单管理','新增菜单','api',NULL,NULL,NULL,NULL),
(29,'sys:notice:my-detail','我的Notices详情','Notices管理','查看我的Notices详情','api',NULL,NULL,NULL,NULL),
(30,'sys:role:query','角色列表','角色管理','获取角色分页列表','api',NULL,NULL,NULL,NULL),
(31,'sys:dict:details','查看字典','字典管理','查看Dict详情','api',NULL,NULL,NULL,NULL),
(32,'sys:menu:routes','路由列表','菜单管理','查看菜单路由','api',NULL,NULL,NULL,NULL),
(33,'sys:noticereceiver:view','NoticeReceiver详情','NoticeReceiver管理','查看NoticeReceiver详情','api',NULL,NULL,NULL,NULL),
(34,'sys:notice:read-all','标记全部为已读','Notices管理','标记全部为已读','api',NULL,NULL,NULL,NULL),
(35,'sys:dept:edit','编辑部门','部门管理','更新部门','api',NULL,NULL,NULL,NULL),
(36,'sys:noticereceiver:update','更新NoticeReceiver','NoticeReceiver管理','更新NoticeReceiver','api',NULL,NULL,NULL,NULL),
(37,'sys:perm:options','权限下拉列表','角色管理','获取权限下拉选项','api',NULL,NULL,NULL,NULL),
(38,'sys:dict:query','字典查询','字典管理','查看Dict列表','api',NULL,NULL,NULL,NULL),
(39,'sys:dict:add','新增字典','字典管理','创建Dict','api',NULL,NULL,NULL,NULL),
(40,'sys:dict-item:delete','删除字典项','字典管理','删除字典项','api',NULL,NULL,NULL,NULL),
(41,'sys:notice:publish','发布公告','Notices管理','发布公告','api',NULL,NULL,NULL,NULL),
(42,'sys:user:add','用户新增','用户管理','创建新用户','api',NULL,NULL,NULL,NULL),
(43,'sys:notice:revoke','撤销公告','Notices管理','撤销公告','api',NULL,NULL,NULL,NULL),
(44,'sys:user:page','用户分页列表','用户管理','获取用户分页列表','api',NULL,NULL,NULL,NULL),
(45,'sys:user:profile','个人信息','个人中心','获取当前登录用户的个人中心信息','api',NULL,NULL,NULL,NULL),
(46,'sys:role:menu:update','更新角色菜单','角色管理','更新角色菜单','api',NULL,NULL,NULL,NULL),
(47,'sys:user:options','用户下拉选项','用户管理','获取用户下拉选项','api',NULL,NULL,NULL,NULL),
(48,'sys:notice:update','更新Notices','Notices管理','更新Notices','api',NULL,NULL,NULL,NULL),
(49,'sys:notice:form','Notices表单','Notices管理','查看Notices详情','api',NULL,NULL,NULL,NULL),
(50,'sys:role:delete','删除角色','角色管理','删除角色','api',NULL,NULL,NULL,NULL),
(51,'sys:config:details','Config详情','Config管理','查看Config详情','api',NULL,NULL,NULL,NULL),
(52,'sys:dept:add','创建部门','部门管理','创建部门','api',NULL,NULL,NULL,NULL),
(53,'sys:dept:delete','删除部门','部门管理','删除部门','api',NULL,NULL,NULL,NULL),
(54,'sys:dict:edit','编辑字典','字典管理','创建Dict','api',NULL,NULL,NULL,NULL),
(55,'sys:menu:delete','删除菜单','菜单管理','删除菜单','api',NULL,NULL,NULL,NULL),
(56,'sys:menu:query','菜单查询','菜单管理','查看菜单列表','api',NULL,NULL,NULL,NULL),
(57,'sys:notice:detail','Notices详情','Notices管理','查看Notices详情','api',NULL,NULL,NULL,NULL),
(58,'sys:role:add','新增角色','角色管理','创建角色','api',NULL,NULL,NULL,NULL),
(59,'sys:noticereceiver:query','NoticeReceiver列表','NoticeReceiver管理','查看NoticeReceiver列表','api',NULL,NULL,NULL,NULL),
(60,'sys:user:me','获取当前用户信息','个人中心','获取当前登录用户的详细信息','api',NULL,NULL,NULL,NULL),
(61,'sys:user:delete','删除用户','用户管理','删除用户','api',NULL,NULL,NULL,NULL),
(62,'sys:user:edit','用户编辑','用户管理','更新用户信息','api',NULL,NULL,NULL,NULL),
(63,'sys:user:change-password','修改密码','个人中心','修改当前登录用户的密码','api',NULL,NULL,NULL,NULL),
(64,'sys:dict:delete','删除字典项','字典项管理','删除字典项','api',NULL,NULL,NULL,NULL),
(65,'sys:notice:mynotice','我的公告列表','Notices管理','查看我的列表','api',NULL,NULL,NULL,NULL),
(66,'sys:role:detail','角色详情','角色管理','获取角色详情','api',NULL,NULL,NULL,NULL),
(67,'sys:role:perm','角色权限列表','角色管理','获取角色权限','api',NULL,NULL,NULL,NULL),
(68,'sys:user:update-profile','修改个人信息','个人中心','修改当前登录用户的个人中心信息','api',NULL,NULL,NULL,NULL),
(69,'sys:dept:query','部门列表','部门管理','获取部门列表','api',NULL,NULL,NULL,NULL);
/*!40000 ALTER TABLE `permissions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `role_menu`
--

DROP TABLE IF EXISTS `role_menu`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `role_menu` (
  `role_id` bigint(20) NOT NULL,
  `menu_id` bigint(20) NOT NULL,
  PRIMARY KEY (`role_id`,`menu_id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `role_menu`
--

LOCK TABLES `role_menu` WRITE;
/*!40000 ALTER TABLE `role_menu` DISABLE KEYS */;
INSERT INTO `role_menu` VALUES
(1,1),
(1,2),
(1,4),
(1,5),
(1,6),
(1,20),
(1,21),
(1,22),
(1,23),
(1,24),
(2,1),
(2,4),
(2,5),
(2,6),
(11,1),
(11,2),
(11,4),
(11,5),
(11,6),
(11,20),
(11,21),
(11,22),
(11,23),
(11,24);
/*!40000 ALTER TABLE `role_menu` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `role_permissions`
--

DROP TABLE IF EXISTS `role_permissions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `role_permissions` (
  `role_id` bigint(20) NOT NULL,
  `permission_code` varchar(50) COLLATE utf8_unicode_ci NOT NULL,
  PRIMARY KEY (`role_id`,`permission_code`),
  KEY `permission_id` (`permission_code`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `role_permissions`
--

LOCK TABLES `role_permissions` WRITE;
/*!40000 ALTER TABLE `role_permissions` DISABLE KEYS */;
INSERT INTO `role_permissions` VALUES
(1,'sys:config:add'),
(1,'sys:config:delete'),
(1,'sys:config:details'),
(1,'sys:config:query'),
(1,'sys:config:update'),
(1,'sys:config:view'),
(1,'sys:dept:add'),
(1,'sys:dept:delete'),
(1,'sys:dept:edit'),
(1,'sys:dept:options'),
(1,'sys:dept:query'),
(1,'sys:dept:view'),
(1,'sys:dict-item:add'),
(1,'sys:dict-item:delete'),
(1,'sys:dict-item:details'),
(1,'sys:dict-item:edit'),
(1,'sys:dict-item:form'),
(1,'sys:dict-item:query'),
(1,'sys:dict:add'),
(1,'sys:dict:delete'),
(1,'sys:dict:details'),
(1,'sys:dict:edit'),
(1,'sys:dict:query'),
(1,'sys:menu:add'),
(1,'sys:menu:delete'),
(1,'sys:menu:details'),
(1,'sys:menu:edit'),
(1,'sys:menu:options'),
(1,'sys:menu:query'),
(1,'sys:menu:routes'),
(1,'sys:notice:add'),
(1,'sys:notice:delete'),
(1,'sys:notice:detail'),
(1,'sys:notice:form'),
(1,'sys:notice:my-detail'),
(1,'sys:notice:mynotice'),
(1,'sys:notice:publish'),
(1,'sys:notice:query'),
(1,'sys:notice:read-all'),
(1,'sys:notice:revoke'),
(1,'sys:notice:update'),
(1,'sys:perm:options'),
(1,'sys:role:add'),
(1,'sys:role:delete'),
(1,'sys:role:detail'),
(1,'sys:role:edit'),
(1,'sys:role:menu'),
(1,'sys:role:menu:update'),
(1,'sys:role:options'),
(1,'sys:role:perm'),
(1,'sys:role:perm:update'),
(1,'sys:role:query'),
(1,'sys:user:add'),
(1,'sys:user:change-password'),
(1,'sys:user:delete'),
(1,'sys:user:edit'),
(1,'sys:user:info'),
(1,'sys:user:me'),
(1,'sys:user:options'),
(1,'sys:user:page'),
(1,'sys:user:profile'),
(1,'sys:user:reset-password'),
(1,'sys:user:update-profile'),
(2,'sys:menu:add'),
(2,'sys:menu:delete'),
(2,'sys:menu:details'),
(2,'sys:menu:edit'),
(2,'sys:menu:options'),
(2,'sys:menu:query'),
(2,'sys:menu:routes'),
(2,'sys:role:add'),
(2,'sys:role:detail'),
(2,'sys:role:menu'),
(2,'sys:role:options'),
(2,'sys:role:query'),
(2,'sys:user:change-password'),
(2,'sys:user:create'),
(2,'sys:user:delete'),
(2,'sys:user:edit'),
(2,'sys:user:info'),
(2,'sys:user:me'),
(2,'sys:user:page'),
(2,'sys:user:profile'),
(2,'sys:user:reset-password'),
(2,'sys:user:update-profile'),
(11,'sys:dict-item:delete'),
(11,'sys:dict:add'),
(11,'sys:dict:details'),
(11,'sys:dict:edit'),
(11,'sys:dict:query'),
(11,'sys:menu:routes'),
(11,'sys:user:change-password'),
(11,'sys:user:create'),
(11,'sys:user:delete'),
(11,'sys:user:edit'),
(11,'sys:user:info'),
(11,'sys:user:me'),
(11,'sys:user:page'),
(11,'sys:user:profile'),
(11,'sys:user:reset-password'),
(11,'sys:user:update-profile');
/*!40000 ALTER TABLE `role_permissions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `roles`
--

DROP TABLE IF EXISTS `roles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `roles` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(191) COLLATE utf8_unicode_ci NOT NULL,
  `description` longtext COLLATE utf8_unicode_ci,
  `code` varchar(255) CHARACTER SET utf8 COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '角色编码',
  `status` tinyint(1) DEFAULT '1' COMMENT '0-禁用，1-可用',
  `sort` int(11) DEFAULT '1',
  `data_scope` tinyint(1) DEFAULT '1' COMMENT '1--全部数据，2--部门机子部门数据 ，3--本部门数据，4--本人数据,5--自定义部门数据',
  `dept_ids` varchar(100) CHARACTER SET utf8 COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '自定义部门权限，data_scope=5时生效，  1,2,3',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_roles_name` (`name`),
  UNIQUE KEY `roles_unique` (`code`),
  KEY `idx_roles_deleted_at` (`deleted_at`)
) ENGINE=MyISAM AUTO_INCREMENT=12 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `roles`
--

LOCK TABLES `roles` WRITE;
/*!40000 ALTER TABLE `roles` DISABLE KEYS */;
INSERT INTO `roles` VALUES
(1,'super_admin','管理员','SUPER_ADMIN',1,3,1,'1,10,11','2025-05-25 17:52:19.000','2025-06-12 12:41:07.439',NULL),
(2,'admin','普通管理员','ADMIN',1,2,1,NULL,'2025-05-25 20:40:25.548','2025-06-12 14:37:59.511',NULL),
(11,'11',NULL,'11',1,1,4,NULL,'2025-06-11 19:48:47.376','2025-06-12 13:05:12.919',NULL);
/*!40000 ALTER TABLE `roles` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_dept`
--

DROP TABLE IF EXISTS `user_dept`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_dept` (
  `user_id` bigint(20) unsigned DEFAULT NULL,
  `dept_id` bigint(20) unsigned DEFAULT NULL
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_dept`
--

LOCK TABLES `user_dept` WRITE;
/*!40000 ALTER TABLE `user_dept` DISABLE KEYS */;
/*!40000 ALTER TABLE `user_dept` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_mfas`
--

DROP TABLE IF EXISTS `user_mfas`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_mfas` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) unsigned NOT NULL,
  `secret_key` longtext COLLATE utf8_unicode_ci NOT NULL,
  `is_enabled` tinyint(1) DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_mfas_user_id` (`user_id`)
) ENGINE=MyISAM AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_mfas`
--

LOCK TABLES `user_mfas` WRITE;
/*!40000 ALTER TABLE `user_mfas` DISABLE KEYS */;
INSERT INTO `user_mfas` VALUES
(1,1,'PMUKB3QRCFN2RXJ6OX54F5HVFPHPTPK6',1,'2025-05-25 16:24:14.163','2025-05-25 16:24:14.163');
/*!40000 ALTER TABLE `user_mfas` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_roles`
--

DROP TABLE IF EXISTS `user_roles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_roles` (
  `user_id` bigint(20) NOT NULL,
  `role_id` bigint(20) NOT NULL,
  PRIMARY KEY (`user_id`,`role_id`),
  KEY `role_id` (`role_id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_roles`
--

LOCK TABLES `user_roles` WRITE;
/*!40000 ALTER TABLE `user_roles` DISABLE KEYS */;
INSERT INTO `user_roles` VALUES
(1,1),
(2,1),
(2,2),
(3,1),
(3,2),
(4,2),
(4,11),
(5,1),
(5,2),
(6,11);
/*!40000 ALTER TABLE `user_roles` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `username` varchar(50) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL COMMENT '用户名',
  `password` varchar(255) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL COMMENT '密码',
  `nickname` varchar(100) CHARACTER SET utf8 COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '昵称',
  `avatar` varchar(255) CHARACTER SET utf8 COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '头像',
  `mobile` varchar(20) CHARACTER SET utf8 COLLATE utf8_unicode_ci DEFAULT NULL,
  `gender` varchar(10) COLLATE utf8_unicode_ci NOT NULL DEFAULT '0' COMMENT '0-未知，1-男，2-女',
  `email` varchar(100) CHARACTER SET utf8 COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '邮箱',
  `dept_id` int(11) DEFAULT NULL,
  `open_id` varchar(100) CHARACTER SET utf8 COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '微信openid',
  `creator_id` int(11) DEFAULT NULL COMMENT '创建者id',
  `salt` varchar(50) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL COMMENT '随机字符串（Salt）',
  `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '1-启用 0-禁用',
  `last_login_time` datetime DEFAULT NULL,
  `last_login_ip` varchar(100) CHARACTER SET utf8 COLLATE utf8_unicode_ci DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `last_login_device` varchar(255) CHARACTER SET utf8 COLLATE utf8_unicode_ci DEFAULT NULL,
  `last_login_os` varchar(255) CHARACTER SET utf8 COLLATE utf8_unicode_ci DEFAULT NULL,
  `last_login_browser` varchar(255) CHARACTER SET utf8 COLLATE utf8_unicode_ci DEFAULT NULL,
  `last_login_user_agent` varchar(255) CHARACTER SET utf8 COLLATE utf8_unicode_ci DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`)
) ENGINE=MyISAM AUTO_INCREMENT=7 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES
(1,'admin','$2a$10$ReJ0zQcIdTLFa9QWW0qDmOLmNiB/vOuIdqlzDb0w7xTYYBzCl19zS','admin','','','1','',1,NULL,1,'jkwcOzb1V7mTANy0IkcZAA==',1,'2025-06-16 09:40:23','::1','2025-05-25 09:52:19',NULL,NULL,NULL,NULL,'2025-06-16 01:40:23',NULL),
(3,'gary_deleted_20250610103834','$2a$10$lOC9YP8Fuf8aBRrR7k9YsuG1g50ViCfIAFC83xQOfwzpcMims5MaW','强哥','','17731185460','1','94204000@qq.com',11,NULL,NULL,'aHI6kyOgXkfiZRriouDFvQ==',1,NULL,NULL,'2025-06-10 02:38:29',NULL,NULL,NULL,NULL,'2025-06-10 02:38:34','2025-06-10 02:38:34'),
(4,'gary','$2a$10$zj6mF.XIxFC78JkTAGQK2OSqrbVhzhWwCQVUc8HHloRZSL.vM/L82','强哥','','17731185460','1','94204000@qq.com',1,NULL,NULL,'IVA1Yat0vT+c2VOIUOPSMQ==',1,'2025-06-15 22:06:32','::1','2025-06-10 02:54:34','Desktop','Windows 10','Chrome','Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36','2025-06-15 14:06:32',NULL),
(5,'miracle','$2a$10$anoSxc0hqH/mK8DlqDFv..0dDSB6CNNWAGf4JpunGVBFqvkl0lepW','夜雨','','','1','',11,NULL,4,'e5tsXQr0d7cMDFvL5Bn/fQ==',1,'2025-06-16 09:55:06','::1','2025-06-10 04:16:00','Desktop','Windows 10','Chrome','Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36','2025-06-16 01:55:06',NULL),
(6,'bobo','$2a$10$lC./w7zwngVV/QK.MTvoUuo78jTXlCr75ChVufGnjkR8PAm5aI.B6','bobo','','','1','',1,'',NULL,'FUIK6Hx2CPIwXcns2JNG/Q==',1,NULL,NULL,'2025-06-16 01:40:47',NULL,NULL,NULL,NULL,'2025-06-16 01:40:47',NULL);
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping routines for database 'vireo_gin_admin'
--
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*M!100616 SET NOTE_VERBOSITY=@OLD_NOTE_VERBOSITY */;

-- Dump completed on 2025-06-16 15:08:26
