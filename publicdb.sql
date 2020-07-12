-- phpMyAdmin SQL Dump
-- version 4.0.4.1
-- http://www.phpmyadmin.net
--
-- Host: localhost
-- Generation Time: Feb 23, 2018 at 04:01 PM
-- Server version: 5.1.73
-- PHP Version: 5.3.3

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;

--
-- Database: `publicdb`
--
CREATE DATABASE IF NOT EXISTS `publicdb` DEFAULT CHARACTER SET latin1 COLLATE latin1_swedish_ci;
USE `publicdb`;

-- --------------------------------------------------------

--
-- Table structure for table `AcquisitionData`
--

CREATE TABLE IF NOT EXISTS `AcquisitionData` (
  `DEF_id` int(11) NOT NULL,
  `DEF_timestamp` timestamp NULL DEFAULT NULL,
  `acquisitionname` text NOT NULL,
  `TXT_notes` text NOT NULL,
  PRIMARY KEY (`DEF_id`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `DataFile`
--

CREATE TABLE IF NOT EXISTS `DataFile` (
  `DEF_id` int(20) NOT NULL,
  `DEF_timestamp` timestamp NULL DEFAULT NULL,
  `REF|TiltSeriesData|tiltseries` varchar(24) DEFAULT NULL,
  `TXT_notes` text,
  `filetype` text,
  `auto` int(1) NOT NULL DEFAULT '0',
  `filename` text,
  `grab` tinyint(4) NOT NULL,
  `zoom` float NOT NULL,
  `xcenter` int(11) NOT NULL,
  `ycenter` int(11) NOT NULL,
  `zcenter` int(11) NOT NULL,
  `xangle` float NOT NULL,
  `yangle` float NOT NULL,
  `zangle` float NOT NULL,
  `REF|ThreeDFile|image` text,
  `status` tinyint(4) NOT NULL DEFAULT '0',
  `delete_text` text,
  `delete_time` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`DEF_id`),
  KEY `REF|TiltSeriesData|tiltseries` (`REF|TiltSeriesData|tiltseries`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `FeatureData`
--

CREATE TABLE IF NOT EXISTS `FeatureData` (
  `DEF_id` int(20) NOT NULL,
  `featurename` text NOT NULL,
  PRIMARY KEY (`DEF_id`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `GroupData`
--

CREATE TABLE IF NOT EXISTS `GroupData` (
  `DEF_id` int(2) NOT NULL,
  `DEF_timestamp` timestamp NULL DEFAULT NULL,
  `name` text,
  `description` text,
  PRIMARY KEY (`DEF_id`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `ListFeature`
--

CREATE TABLE IF NOT EXISTS `ListFeature` (
  `DEF_id` int(11) NOT NULL,
  `REF|TiltSeriesData|tiltseries` varchar(24) DEFAULT NULL,
  `REF|FeatureData|featurename` text,
  `checked` text,
  PRIMARY KEY (`DEF_id`),
  KEY `REF|TiltSeriesData|tiltseries` (`REF|TiltSeriesData|tiltseries`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `PublicationData`
--

CREATE TABLE IF NOT EXISTS `PublicationData` (
  `DEF_id` int(4) NOT NULL,
  `DEF_timestamp` timestamp NULL DEFAULT NULL,
  `title` text NOT NULL,
  `authors` text,
  `journal` text,
  `issue` text,
  `pages` text,
  `year` text NOT NULL,
  PRIMARY KEY (`DEF_id`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `ScopeData`
--

CREATE TABLE IF NOT EXISTS `ScopeData` (
  `DEF_id` int(11) NOT NULL,
  `DEF_timestamp` timestamp NULL DEFAULT NULL,
  `scopename` text,
  `TXT_notes` text,
  PRIMARY KEY (`DEF_id`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `SpeciesData`
--

CREATE TABLE IF NOT EXISTS `SpeciesData` (
  `DEF_id` int(3) NOT NULL,
  `DEF_timestamp` timestamp NULL DEFAULT NULL,
  `SpeciesName` text NOT NULL,
  `strain` text,
  `tax_id` int(11) DEFAULT '0',
  `TXT_notes` text,
  `count` int(15) NOT NULL DEFAULT '0',
  PRIMARY KEY (`DEF_id`),
  FULLTEXT KEY `SpeciesName` (`SpeciesName`,`strain`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `ThreeDFile`
--

CREATE TABLE IF NOT EXISTS `ThreeDFile` (
  `DEF_id` int(20) NOT NULL,
  `DEF_timestamp` timestamp NULL DEFAULT NULL,
  `REF|TiltSeriesData|tiltseries` varchar(24) DEFAULT NULL,
  `title` text,
  `TXT_notes` text,
  `classify` text,
  `filename` text,
  `pixel_size` double DEFAULT '0',
  `status` tinyint(4) NOT NULL DEFAULT '0',
  `delete_text` text,
  `delete_time` timestamp NULL DEFAULT NULL,
  `tag` tinyint(4) NOT NULL DEFAULT '0',
  PRIMARY KEY (`DEF_id`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `TiltSeriesData`
--

CREATE TABLE IF NOT EXISTS `TiltSeriesData` (
  `DEF_id` int(20) NOT NULL,
  `DEF_timestamp` timestamp NULL DEFAULT NULL,
  `tiltseriesID` varchar(24) DEFAULT NULL,
  `title` text,
  `TXT_notes` text,
  `time_modified` timestamp NULL DEFAULT NULL,
  `REF|UserData|user` text,
  `tomo_date` date DEFAULT NULL,
  `keywords` text,
  `roles` text,
  `REF|SpeciesData|specie` text,
  `treatment` text,
  `sample` text,
  `single_dual` int(1) DEFAULT NULL,
  `tilt_min` double DEFAULT NULL,
  `tilt_max` double DEFAULT NULL,
  `tilt_step` text,
  `tilt_constant` int(1) DEFAULT NULL,
  `dosage` double DEFAULT NULL,
  `defocus` double DEFAULT NULL,
  `magnification` double DEFAULT NULL,
  `voxel` double DEFAULT '0',
  `scope` text NOT NULL,
  `software_acquisition` text,
  `software_process` text,
  `REF_PublicationData_publication` text,
  `loadmethod` text,
  `loadpath` text,
  `searchtext` text,
  `pubtext` text,
  `raptorcheck` tinyint(4) NOT NULL DEFAULT '0',
  `keyimg` int(1) NOT NULL DEFAULT '0',
  `keymov` int(1) NOT NULL DEFAULT '0',
  `visited` int(11) NOT NULL DEFAULT '0',
  `feature` char(100) DEFAULT NULL,
  `status` tinyint(4) NOT NULL DEFAULT '0',
  `ispublic` tinyint(4) NOT NULL DEFAULT '0',
  `emdb` varchar(24) DEFAULT NULL,
  `delete_text` text,
  `delete_time` timestamp NULL DEFAULT NULL,
  `pipeline` tinyint(4) NOT NULL DEFAULT '0',
  `proj1` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`DEF_id`),
  FULLTEXT KEY `title` (`title`,`TXT_notes`,`tiltseriesID`,`roles`,`treatment`,`sample`,`software_process`,`pubtext`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `UserData`
--

CREATE TABLE IF NOT EXISTS `UserData` (
  `DEF_id` int(3) NOT NULL,
  `DEF_timestamp` timestamp NULL DEFAULT NULL,
  `username` text,
  `fullname` text,
  `var` text NOT NULL,
  `email` text,
  `count` int(15) NOT NULL DEFAULT '0',
  `REF|GroupData|group` text,
  PRIMARY KEY (`DEF_id`),
  FULLTEXT KEY `fullname` (`fullname`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
