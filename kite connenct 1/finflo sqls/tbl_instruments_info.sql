SELECT * FROM finflo_base_db.tbl_instruments_info;CREATE TABLE `tbl_instruments_info` (
  `index` bigint DEFAULT NULL,
  `instrument_token` text COLLATE utf8mb4_unicode_ci,
  `exchange_token` bigint DEFAULT NULL,
  `tradingsymbol` text COLLATE utf8mb4_unicode_ci,
  `name` text COLLATE utf8mb4_unicode_ci,
  `last_price` bigint DEFAULT NULL,
  `expiry` text COLLATE utf8mb4_unicode_ci,
  `strike` double DEFAULT NULL,
  `tick_size` double DEFAULT NULL,
  `lot_size` bigint DEFAULT NULL,
  `instrument_type` text COLLATE utf8mb4_unicode_ci,
  `segment` text COLLATE utf8mb4_unicode_ci,
  `exchange` text COLLATE utf8mb4_unicode_ci,
  KEY `ix_tbl_instruments_info_index` (`index`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
