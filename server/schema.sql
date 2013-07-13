# ************************************************************
# Sequel Pro SQL dump
# Version 4096
#
# http://www.sequelpro.com/
# http://code.google.com/p/sequel-pro/
#
# Host: localhost (MySQL 5.5.25a)
# Database: weather
# Generation Time: 2013-07-13 13:55:31 +0000
# ************************************************************


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


# Dump of table weather
# ------------------------------------------------------------

CREATE TABLE `weather` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `TimeSampled` datetime DEFAULT NULL,
  `BarometerTrend` float DEFAULT NULL,
  `Barometer` float DEFAULT NULL,
  `InsideTemperature` float DEFAULT NULL,
  `InsideHumidity` int(11) DEFAULT NULL,
  `OutsideTemerature` float DEFAULT NULL,
  `WindSpeed` float DEFAULT NULL,
  `WindDirection` int(11) DEFAULT NULL,
  `AverageWindSpeed10Minute` float DEFAULT NULL,
  `AverageWindSpeed2Minute` float DEFAULT NULL,
  `WindGust10Minute` float DEFAULT NULL,
  `WindDirectionGust10Minute` int(11) DEFAULT NULL,
  `DewPoint` float DEFAULT NULL,
  `OutsideHumidity` int(11) DEFAULT NULL,
  `HeatIndex` float DEFAULT NULL,
  `WindChill` float DEFAULT NULL,
  `THSWIndex` float DEFAULT NULL,
  `RainRate` float DEFAULT NULL,
  `UVIndex` int(11) DEFAULT NULL,
  `SolarRadiation` int(11) DEFAULT NULL,
  `StormRain` float DEFAULT NULL,
  `DailyRain` float DEFAULT NULL,
  `Last15MinuteRain` float DEFAULT NULL,
  `LastHourRain` float DEFAULT NULL,
  `DailyET` float DEFAULT NULL,
  `Last24Rain` float DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `TimeIndex` (`TimeSampled`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;




/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
