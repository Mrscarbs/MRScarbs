DELIMITER $$

CREATE PROCEDURE `finflo_base_db`.`stp_insert_or_Update_QuantStats`(
    IN p_nsortino DOUBLE,
    IN p_ninstrument_token INT,
    IN p_nsharpe DOUBLE,
    IN p_nlast_update_time BIGINT
)
BEGIN
    INSERT INTO tbl_quant_stats (nsortino, ninstrument_token, nsharpe, nlast_update_time)
    VALUES (p_nsortino, p_ninstrument_token, p_nsharpe, p_nlast_update_time)
    ON DUPLICATE KEY UPDATE
        nsortino = VALUES(nsortino),
        nsharpe = VALUES(nsharpe),
        nlast_update_time = VALUES(nlast_update_time);
END$$

DELIMITER ;