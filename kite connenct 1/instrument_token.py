import requests
import pandas as pd
import json
import csv
import sqlalchemy
from sqlalchemy import create_engine
import datetime
import time

url = "https://api.kite.trade/instruments"
headers = {
    "X-Kite-Version": "3",
    "Authorization": "246y3a7zlg83xv2l:SuparnSaumitra"
}

def instrument_writer():
    try:
        
        response = requests.get(url, headers=headers)
        response.raise_for_status()  
        print("Status Code:", response.status_code)
    
        text_response = response.text

    
        with open("instrument_data.csv", "w", newline="") as csv_file:
            writer = csv.writer(csv_file)
            lines = response.text.splitlines()
            c=0
            for line in lines:
                if c>=1 and c<=len(lines):
                    writer.writerow([line])
                c= c+1
            print("Response data written to response_data.csv successfully.")

        column_names = [
            'instrument_token', 'exchange_token', 'tradingsymbol', 'name', 'last_price',
            'expiry', 'strike', 'tick_size', 'lot_size', 'instrument_type',
            'segment', 'exchange'
        ]

        # Read CSV with specified column names
        instrument_csv = pd.read_csv("instrument_data.csv", quoting=csv.QUOTE_NONE,names=column_names)
        print(instrument_csv)
        print(instrument_csv['exchange_token'])
        conn_string = 'mysql+pymysql://root:Karma100%@localhost/finflo_base_db'
        engine = create_engine(conn_string)
        instrument_csv.to_sql("tbl_instruments_info", engine, if_exists='replace' )
   
    
    except requests.exceptions.RequestException as e:
        print("Error:", e)



now = datetime.datetime.now()

if now.hour == 23 and now.minute<=10:
    instrument_writer()
    time.sleep(700)