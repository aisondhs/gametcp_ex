
CREATE DATABASE IF NOT EXISTS test DEFAULT CHARSET utf8 COLLATE utf8_general_ci;

 CREATE TABLE `user` (

  `Uid` int(11) NOT NULL AUTO_INCREMENT,

  `Account` varchar(64) NOT NULL DEFAULT '',

  `Pwd` varchar(64) NOT NULL DEFAULT '',

  `Ctime` int(11) NOT NULL DEFAULT 0,

  PRIMARY KEY (`Uid`)

) ENGINE=MyISAM AUTO_INCREMENT=100001 DEFAULT CHARSET=utf8;

CREATE TABLE `role` (

  `RoleId` int(11) NOT NULL AUTO_INCREMENT,

  `Uid` int(11) NOT NULL DEFAULT 0,

  `Areaid` mediumint(11) NOT NULL DEFAULT 0,

  `Ctime` int(11) NOT NULL DEFAULT 0,

  PRIMARY KEY (`RoleId`)

) ENGINE=MyISAM AUTO_INCREMENT=100001 DEFAULT CHARSET=utf8;