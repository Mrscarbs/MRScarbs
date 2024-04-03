DELIMITER $$
CREATE DEFINER=`root`@`localhost` PROCEDURE `stp_insert_current_quote`(IN p_name TEXT)
BEGIN
  DECLARE v_instrument_token INT;
  DECLARE v_quote INT;
  DECLARE v_timestamp BIGINT;
  
  -- Check if the name exists in tbl_instruments_info
  IF EXISTS (SELECT 1 FROM tbl_instruments_info WHERE stradingsymbol = p_name) THEN
    -- Get the instrument token for the name
    SELECT ninstrument_token INTO v_instrument_token 
    FROM tbl_instruments_info
    WHERE stradingsymbol = p_name
    LIMIT 1;
    
    -- Generate a random quote value (replace with actual logic)
    SET v_quote = 0;
    
    -- Get the current timestamp
    SET v_timestamp = UNIX_TIMESTAMP();
    
    -- Insert or update the record in tbl_current_quotes
    INSERT INTO tbl_current_quotes (ntimestamp, nquote, ninstrument_token)
    VALUES (v_timestamp, v_quote, v_instrument_token)
    ON DUPLICATE KEY UPDATE 
      ntimestamp = v_timestamp,
      nquote = v_quote;
  END IF;
END$$
DELIMITER ;
