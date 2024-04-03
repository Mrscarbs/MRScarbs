DELIMITER $$
CREATE PROCEDURE `stp_get_limits`(in papi_id int, in pinterval varchar(45))
BEGIN
select tbl_max_historical_limit.nlimit from tbl_max_historical_limit
where tbl_max_historical_limit.napi_id = papi_id and tbl_max_historical_limit.sinterval =  pinterval;
END$$
DELIMITER ;
