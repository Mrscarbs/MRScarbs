CREATE TABLE `api_sys_config` (
  `api_key` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `secret_key` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `api_provider` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `access_token` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `last_purchase_time` bigint DEFAULT NULL,
  `first_purchase_time` bigint DEFAULT NULL,
  `historical` tinyint DEFAULT NULL,
  `instrument_type` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `api_id` int NOT NULL,
  PRIMARY KEY (`api_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
