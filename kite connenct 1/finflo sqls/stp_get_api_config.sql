DELIMITER $$
CREATE PROCEDURE `stp_get_api_config`(
    IN p_api_id INT
)
BEGIN
    select * from api_sys_config
    where api_sys_config.api_id = api_id;
END$$
DELIMITER ;
