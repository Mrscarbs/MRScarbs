DELIMITER $$
CREATE DEFINER=`root`@`localhost` PROCEDURE `stp_get_instrument_info`(IN p_instrument_token VARCHAR(255))
BEGIN
    SELECT * FROM finflo_base_db.tbl_instruments_info WHERE ninstrument_token = p_instrument_token;
END$$
DELIMITER ;
