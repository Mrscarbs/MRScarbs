CREATE TABLE `finflo_base_db`.`tbl_quant_stats` (
  `nsortino` DOUBLE NULL,
  `ninstrument_token` INT NOT NULL,
  `nsharpe` DOUBLE NULL,
  `nlast_update_time` BIGINT NULL,
  PRIMARY KEY (`ninstrument_token`));
