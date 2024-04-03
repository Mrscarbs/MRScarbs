DELIMITER $$
CREATE DEFINER=`root`@`localhost` PROCEDURE `stp_get_tbl_current_ltp`()
BEGIN
select * from tbl_current_ltp;

END$$
DELIMITER ;
call stp_get_tbl_current_ltp();