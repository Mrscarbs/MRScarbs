DELIMITER $$
CREATE DEFINER=`root`@`localhost` PROCEDURE `stp_Update_Instrument`(
    IN InstrumentID INT,
    IN NewTimestamp BIGINT,
    IN NewQuote INT
)
BEGIN
    UPDATE tbl_current_quotes
    SET 
        ntimestamp = NewTimestamp,
        nquote = NewQuote
    WHERE
        ninstrument_token = InstrumentID;
END$$
DELIMITER ;
