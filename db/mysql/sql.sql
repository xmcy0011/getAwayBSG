SELECT table_name FROM information_schema.TABLES WHERE table_name ='lianjia';

CREATE TABLE `lianjia` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `Title` varchar(200) DEFAULT NULL,
  `ListCrawlTime` datetime DEFAULT NULL,
  `DetailStatus` tinyint(4) DEFAULT NULL,
  `TotalPrice` int(11) DEFAULT NULL,
  `UnitPrice` int(11) DEFAULT NULL,
  `Link` varchar(128) DEFAULT NULL,
  `HouseRecordLJ` varchar(32) DEFAULT NULL,
  `DetailCrawlTime` datetime DEFAULT NULL,
  `AreaName` varchar(64) DEFAULT NULL,
  `BaseAttr` varchar(500) DEFAULT NULL,
  `BeOnlineTime` varchar(64) DEFAULT NULL,
  `CompletedInfo` varchar(64) DEFAULT NULL,
  `DecorateInfo` varchar(128) DEFAULT NULL,
  `DirectionInfo` varchar(32) DEFAULT NULL,
  `FloorInfo` varchar(32) DEFAULT NULL,
  `RoomInfo` varchar(32) DEFAULT NULL,
  `Size` float DEFAULT NULL,
  `TransactionAttr` varchar(128) DEFAULT NULL,
  `VillageName` varchar(64) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;