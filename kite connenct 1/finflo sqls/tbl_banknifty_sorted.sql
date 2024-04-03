CREATE TABLE `tbl_banknifty_sorted` (
  `ninstrument_token` int NOT NULL,
  `sname` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `sexchange` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `npct_change` int DEFAULT NULL,
  `npoints_change` int DEFAULT NULL,
  `nlast_update_time` bigint DEFAULT NULL,
  `sinterval` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`ninstrument_token`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
