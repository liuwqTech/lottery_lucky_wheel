/*
 Navicat Premium Data Transfer

 Source Server         : mysql
 Source Server Type    : MySQL
 Source Server Version : 50731
 Source Host           : localhost:3306
 Source Schema         : lucky_wheel

 Target Server Type    : MySQL
 Target Server Version : 50731
 File Encoding         : 65001

 Date: 08/06/2023 17:06:49
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for index_def
-- ----------------------------
DROP TABLE IF EXISTS `index_def`;
CREATE TABLE `index_def` (
  `id` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of index_def
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for lw_blackip
-- ----------------------------
DROP TABLE IF EXISTS `lw_blackip`;
CREATE TABLE `lw_blackip` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `ip` varchar(50) DEFAULT NULL COMMENT 'IP地址',
  `blacktime` int(11) DEFAULT NULL COMMENT '黑名单限制到期时间',
  `sys_created` int(11) DEFAULT NULL COMMENT '创建时间',
  `sys_updated` int(11) DEFAULT NULL COMMENT '修改时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of lw_blackip
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for lw_code
-- ----------------------------
DROP TABLE IF EXISTS `lw_code`;
CREATE TABLE `lw_code` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `gift_id` int(11) DEFAULT NULL COMMENT '奖品id，关联lw_gift表',
  `code` varchar(255) DEFAULT NULL COMMENT '虚拟券编码',
  `sys_created` int(11) DEFAULT NULL COMMENT '创建时间',
  `sys_updated` int(11) DEFAULT NULL COMMENT '更新时间',
  `sys_status` smallint(6) DEFAULT NULL COMMENT '状态，0正常，1作废，2已发放',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of lw_code
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for lw_gift
-- ----------------------------
DROP TABLE IF EXISTS `lw_gift`;
CREATE TABLE `lw_gift` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `title` varchar(255) DEFAULT NULL COMMENT '奖品名称',
  `prize_num` int(11) DEFAULT NULL COMMENT '奖品数量，0：无限量，>0：限量，<0：无奖品',
  `left_num` int(11) DEFAULT NULL COMMENT '剩余奖品数量',
  `prize_code` varchar(50) DEFAULT NULL COMMENT '0-9999表示100%，0-0表示万分之一的中奖概率',
  `prize_time` int(11) DEFAULT NULL COMMENT '发奖周期，D天',
  `img` varchar(255) DEFAULT NULL COMMENT '奖品图片',
  `displayorder` int(11) DEFAULT NULL COMMENT '位置序号，小的排在前面',
  `gtype` int(11) DEFAULT NULL COMMENT '奖品类型，0虚拟币，1虚拟券，2实物-小奖，3实物-大奖',
  `gdata` varchar(255) DEFAULT NULL COMMENT '扩展数据，如：虚拟币数量',
  `time_begin` int(11) DEFAULT NULL COMMENT '开始时间',
  `time_end` int(11) DEFAULT NULL COMMENT '结束时间',
  `prize_data` mediumtext COMMENT '发奖计划，[{时间1, 数量1}, {时间2, 数量2}]',
  `prize_begin` int(11) DEFAULT NULL COMMENT '发奖周期的开始',
  `prize_end` int(11) DEFAULT NULL COMMENT '发奖周期的结束',
  `sys_status` smallint(6) DEFAULT NULL COMMENT '状态，0正常，1删除',
  `sys_created` int(11) DEFAULT NULL COMMENT '创建时间',
  `sys_updated` int(11) DEFAULT NULL COMMENT '修改时间',
  `sys_ip` varchar(50) DEFAULT NULL COMMENT '修改人ip',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of lw_gift
-- ----------------------------
BEGIN;
INSERT INTO `lw_gift` (`id`, `title`, `prize_num`, `left_num`, `prize_code`, `prize_time`, `img`, `displayorder`, `gtype`, `gdata`, `time_begin`, `time_end`, `prize_data`, `prize_begin`, `prize_end`, `sys_status`, `sys_created`, `sys_updated`, `sys_ip`) VALUES (1, 'iphone 14 Pro Max', 2, 2, '0-1', 3, 'https://img2.baidu.com/it/u=1791804638,2063291533&fm=253&fmt=auto&app=120&f=JPEG?w=653&h=914', 1, 3, NULL, 20230606, 20230610, '1', 1, 1, 0, 1, 1, '1');
COMMIT;

-- ----------------------------
-- Table structure for lw_result
-- ----------------------------
DROP TABLE IF EXISTS `lw_result`;
CREATE TABLE `lw_result` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `gift_id` int(11) DEFAULT NULL COMMENT '奖品id，关联lw_gift表',
  `gift_name` varchar(255) DEFAULT NULL COMMENT '奖品名称',
  `gift_type` int(11) DEFAULT NULL COMMENT '奖品类型，同lw_gift.gtype',
  `uid` int(11) DEFAULT NULL COMMENT '用户id',
  `username` varchar(50) DEFAULT NULL COMMENT '用户名',
  `prize_code` int(11) DEFAULT NULL COMMENT '抽奖编号（4位数的随机数）',
  `gift_data` varchar(255) DEFAULT NULL COMMENT '获奖信息',
  `sys_created` int(11) DEFAULT NULL COMMENT '创建时间',
  `sys_ip` varchar(50) DEFAULT NULL COMMENT '用户抽奖的IP',
  `sys_status` smallint(6) DEFAULT NULL COMMENT '状态，0正常，1删除，2作弊',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of lw_result
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for lw_user
-- ----------------------------
DROP TABLE IF EXISTS `lw_user`;
CREATE TABLE `lw_user` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `username` varchar(50) DEFAULT NULL COMMENT '用户名',
  `blacktime` int(11) DEFAULT NULL COMMENT '黑名单限制到期时间',
  `realname` varchar(50) DEFAULT NULL COMMENT '联系人',
  `mobile` varchar(50) DEFAULT NULL COMMENT '手机号',
  `address` varchar(255) DEFAULT NULL COMMENT '联系地址',
  `sys_created` int(11) DEFAULT NULL COMMENT '创建时间',
  `sys_updated` int(11) DEFAULT NULL COMMENT '修改时间',
  `sys_ip` varchar(50) DEFAULT NULL COMMENT 'IP地址',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of lw_user
-- ----------------------------
BEGIN;
INSERT INTO `lw_user` (`id`, `username`, `blacktime`, `realname`, `mobile`, `address`, `sys_created`, `sys_updated`, `sys_ip`) VALUES (1, '小新', 0, '刘文强', '18752073970', '北京市', NULL, NULL, NULL);
COMMIT;

-- ----------------------------
-- Table structure for lw_userday
-- ----------------------------
DROP TABLE IF EXISTS `lw_userday`;
CREATE TABLE `lw_userday` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `uid` int(11) DEFAULT NULL COMMENT '用户id',
  `day` int(11) DEFAULT NULL COMMENT '日期，如20180725',
  `num` int(11) DEFAULT NULL COMMENT '次数',
  `sys_created` int(11) DEFAULT NULL COMMENT '创建时间',
  `sys_updated` int(11) DEFAULT NULL COMMENT '修改时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of lw_userday
-- ----------------------------
BEGIN;
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
