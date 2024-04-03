DELIMITER //
CREATE PROCEDURE `stp_GetSInterval` (IN api_id INT)
BEGIN
    SELECT `sinterval` FROM `tbl_max_historical_limit` WHERE `napi_id` = api_id;
END //
DELIMITER ;
