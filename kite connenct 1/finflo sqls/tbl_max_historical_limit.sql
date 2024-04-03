CREATE TABLE `tbl_max_historical_limit` (
  `sinterval` varchar(45) COLLATE utf8mb4_unicode_ci NOT NULL,
  `napi_id` int NOT NULL,
  `nlimit` bigint DEFAULT NULL,
  PRIMARY KEY (`napi_id`,`sinterval`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
