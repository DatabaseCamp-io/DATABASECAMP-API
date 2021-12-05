-- phpMyAdmin SQL Dump
-- version 5.1.1
-- https://www.phpmyadmin.net/
--
-- Host: db
-- Generation Time: Dec 05, 2021 at 02:39 PM
-- Server version: 10.6.4-MariaDB-1:10.6.4+maria~focal
-- PHP Version: 7.4.20

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `database_camp`
--

-- --------------------------------------------------------

--
-- Table structure for table `Activity`
--

CREATE TABLE `Activity` (
  `activity_id` int(11) NOT NULL,
  `activity_type_id` int(11) NOT NULL,
  `content_id` int(11) DEFAULT NULL,
  `point` int(11) NOT NULL,
  `question` text NOT NULL,
  `story` text NOT NULL,
  `activity_order` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- Table structure for table `ActivityType`
--

CREATE TABLE `ActivityType` (
  `activity_type_id` int(11) NOT NULL,
  `name` varchar(20) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- Table structure for table `Badge`
--

CREATE TABLE `Badge` (
  `badge_id` int(11) NOT NULL,
  `name` varchar(50) NOT NULL,
  `icon_path` text NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- Table structure for table `CompletionChoice`
--

CREATE TABLE `CompletionChoice` (
  `completion_choice_id` int(11) NOT NULL,
  `activity_id` int(11) NOT NULL,
  `content` text NOT NULL,
  `question_first` text NOT NULL,
  `question_last` text NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- Table structure for table `Content`
--

CREATE TABLE `Content` (
  `content_id` int(11) NOT NULL,
  `content_group_id` int(11) NOT NULL,
  `name` varchar(50) NOT NULL,
  `video_path` text NOT NULL,
  `slide_path` text NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- Table structure for table `ContentExam`
--

CREATE TABLE `ContentExam` (
  `exam_id` int(11) NOT NULL,
  `content_group_id` int(11) NOT NULL,
  `activity_id` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- Table structure for table `ContentGroup`
--

CREATE TABLE `ContentGroup` (
  `content_group_id` int(11) NOT NULL,
  `name` varchar(50) NOT NULL,
  `badge_id` int(11) NOT NULL,
  `mini_exam_id` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- Table structure for table `Exam`
--

CREATE TABLE `Exam` (
  `exam_id` int(11) NOT NULL,
  `type` enum('PRE','MINI','POST') NOT NULL,
  `instruction` text NOT NULL,
  `created_timestamp` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- Table structure for table `ExamResult`
--

CREATE TABLE `ExamResult` (
  `exam_result_id` int(11) NOT NULL,
  `user_id` int(11) NOT NULL,
  `exam_id` int(11) NOT NULL,
  `is_passed` tinyint(1) NOT NULL,
  `created_timestamp` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- Table structure for table `ExamResultActivity`
--

CREATE TABLE `ExamResultActivity` (
  `exam_result_id` int(11) NOT NULL,
  `activity_id` int(11) NOT NULL,
  `score` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- Table structure for table `Hint`
--

CREATE TABLE `Hint` (
  `hint_id` int(11) NOT NULL,
  `activity_id` int(11) NOT NULL,
  `content` text NOT NULL,
  `point_reduce` int(11) NOT NULL,
  `level` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- Table structure for table `LearningProgression`
--

CREATE TABLE `LearningProgression` (
  `user_id` int(11) NOT NULL,
  `activity_id` int(11) NOT NULL,
  `created_timestamp` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- Table structure for table `MatchingChoice`
--

CREATE TABLE `MatchingChoice` (
  `matching_choice_id` int(11) NOT NULL,
  `activity_id` int(11) NOT NULL,
  `pair_item1` varchar(50) NOT NULL,
  `pair_item2` varchar(50) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- Table structure for table `MultipleChoice`
--

CREATE TABLE `MultipleChoice` (
  `multiple_choice_id` int(11) NOT NULL,
  `activity_id` int(11) NOT NULL,
  `content` text NOT NULL,
  `is_correct` tinyint(4) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- Stand-in structure for view `Profile`
-- (See below for the actual view)
--
CREATE TABLE `Profile` (
`user_id` int(11)
,`name` varchar(100)
,`point` int(11)
,`created_timestamp` timestamp
,`activity_count` bigint(21)
);

-- --------------------------------------------------------

--
-- Stand-in structure for view `Ranking`
-- (See below for the actual view)
--
CREATE TABLE `Ranking` (
`user_id` int(11)
,`name` varchar(100)
,`email` varchar(50)
,`point` int(11)
,`ranking` bigint(21)
);

-- --------------------------------------------------------

--
-- Table structure for table `User`
--

CREATE TABLE `User` (
  `user_id` int(11) NOT NULL,
  `name` varchar(100) NOT NULL,
  `email` varchar(50) NOT NULL,
  `password` text NOT NULL,
  `access_token` text NOT NULL,
  `point` int(11) NOT NULL,
  `expired_token_timestamp` timestamp NULL DEFAULT NULL,
  `created_timestamp` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `updated_timestamp` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- Table structure for table `UserBadge`
--

CREATE TABLE `UserBadge` (
  `user_id` int(11) NOT NULL,
  `badge_id` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- Table structure for table `UserHint`
--

CREATE TABLE `UserHint` (
  `user_id` int(11) NOT NULL,
  `hint_id` int(11) NOT NULL,
  `created_timestamp` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- Structure for view `Profile`
--
DROP TABLE IF EXISTS `Profile`;

CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `Profile`  AS SELECT `User`.`user_id` AS `user_id`, `User`.`name` AS `name`, `User`.`point` AS `point`, `User`.`created_timestamp` AS `created_timestamp`, count(`LearningProgression`.`activity_id`) AS `activity_count` FROM (`User` left join `LearningProgression` on(`User`.`user_id` = `LearningProgression`.`user_id`)) GROUP BY `User`.`user_id` ;

-- --------------------------------------------------------

--
-- Structure for view `Ranking`
--
DROP TABLE IF EXISTS `Ranking`;

CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `Ranking`  AS SELECT `User`.`user_id` AS `user_id`, `User`.`name` AS `name`, `User`.`email` AS `email`, `User`.`point` AS `point`, row_number()  ( order by `User`.`point` desc) AS `over` FROM `User` WHERE 1 ;

--
-- Indexes for dumped tables
--

--
-- Indexes for table `Activity`
--
ALTER TABLE `Activity`
  ADD PRIMARY KEY (`activity_id`),
  ADD UNIQUE KEY `content_id` (`content_id`,`activity_order`);

--
-- Indexes for table `ActivityType`
--
ALTER TABLE `ActivityType`
  ADD PRIMARY KEY (`activity_type_id`);

--
-- Indexes for table `Badge`
--
ALTER TABLE `Badge`
  ADD PRIMARY KEY (`badge_id`);

--
-- Indexes for table `CompletionChoice`
--
ALTER TABLE `CompletionChoice`
  ADD PRIMARY KEY (`completion_choice_id`);

--
-- Indexes for table `Content`
--
ALTER TABLE `Content`
  ADD PRIMARY KEY (`content_id`);

--
-- Indexes for table `ContentExam`
--
ALTER TABLE `ContentExam`
  ADD PRIMARY KEY (`content_group_id`,`activity_id`,`exam_id`);

--
-- Indexes for table `ContentGroup`
--
ALTER TABLE `ContentGroup`
  ADD PRIMARY KEY (`content_group_id`);

--
-- Indexes for table `Exam`
--
ALTER TABLE `Exam`
  ADD PRIMARY KEY (`exam_id`);

--
-- Indexes for table `ExamResult`
--
ALTER TABLE `ExamResult`
  ADD PRIMARY KEY (`exam_result_id`);

--
-- Indexes for table `ExamResultActivity`
--
ALTER TABLE `ExamResultActivity`
  ADD PRIMARY KEY (`exam_result_id`,`activity_id`);

--
-- Indexes for table `Hint`
--
ALTER TABLE `Hint`
  ADD PRIMARY KEY (`hint_id`);

--
-- Indexes for table `LearningProgression`
--
ALTER TABLE `LearningProgression`
  ADD PRIMARY KEY (`user_id`,`activity_id`);

--
-- Indexes for table `MatchingChoice`
--
ALTER TABLE `MatchingChoice`
  ADD PRIMARY KEY (`matching_choice_id`);

--
-- Indexes for table `MultipleChoice`
--
ALTER TABLE `MultipleChoice`
  ADD PRIMARY KEY (`multiple_choice_id`);

--
-- Indexes for table `User`
--
ALTER TABLE `User`
  ADD PRIMARY KEY (`user_id`),
  ADD UNIQUE KEY `email` (`email`);

--
-- Indexes for table `UserBadge`
--
ALTER TABLE `UserBadge`
  ADD PRIMARY KEY (`user_id`,`badge_id`);

--
-- Indexes for table `UserHint`
--
ALTER TABLE `UserHint`
  ADD PRIMARY KEY (`user_id`,`hint_id`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `Activity`
--
ALTER TABLE `Activity`
  MODIFY `activity_id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `ActivityType`
--
ALTER TABLE `ActivityType`
  MODIFY `activity_type_id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `Badge`
--
ALTER TABLE `Badge`
  MODIFY `badge_id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `CompletionChoice`
--
ALTER TABLE `CompletionChoice`
  MODIFY `completion_choice_id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `Content`
--
ALTER TABLE `Content`
  MODIFY `content_id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `ContentGroup`
--
ALTER TABLE `ContentGroup`
  MODIFY `content_group_id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `Exam`
--
ALTER TABLE `Exam`
  MODIFY `exam_id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `ExamResult`
--
ALTER TABLE `ExamResult`
  MODIFY `exam_result_id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `Hint`
--
ALTER TABLE `Hint`
  MODIFY `hint_id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `MatchingChoice`
--
ALTER TABLE `MatchingChoice`
  MODIFY `matching_choice_id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `MultipleChoice`
--
ALTER TABLE `MultipleChoice`
  MODIFY `multiple_choice_id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `User`
--
ALTER TABLE `User`
  MODIFY `user_id` int(11) NOT NULL AUTO_INCREMENT;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
