import pandas as pd
import mysql.connector

# Load the Nifty stocks list from the CSV
url = 'https://archives.nseindia.com/content/indices/ind_nifty50list.csv'
nifty_stocks = pd.read_csv(url)
symbols = nifty_stocks['Symbol'].tolist()

# Database connection parameters
db_user = 'root'
db_password = 'Karma100%'
db_host = 'localhost'
db_name = 'finflo_base_db'

# Create a connection to the database
try:
    connection = mysql.connector.connect(
        host=db_host,
        user=db_user,
        password=db_password,
        database=db_name
    )

    if connection.is_connected():
        db_info = connection.get_server_info()
        print("Connected to MySQL Server version ", db_info)
        cursor = connection.cursor()

        # Execute the stored procedure for each symbol
        for symbol in symbols:
            # The stored procedure call
            procedure_call = f"CALL stp_insert_current_ltp('{symbol}')"
            cursor.execute(procedure_call)
            connection.commit()  # Commit the transaction if necessary

        print("Stored procedure executed for all Nifty stocks.")

except mysql.connector.Error as e:
    print("Error while connecting to MySQL", e)

finally:
    if connection.is_connected():
        cursor.close()
        connection.close()
        print("MySQL connection is closed")

