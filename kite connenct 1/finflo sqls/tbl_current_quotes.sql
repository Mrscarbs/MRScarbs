CREATE TABLE `tbl_current_quotes` (
  `ntimestamp` bigint NOT NULL,
  `nquote` int DEFAULT NULL,
  `ninstrument_token` int DEFAULT NULL,
  `sexchange` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `strading_symbol` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `sinstrument_type` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
