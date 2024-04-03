import pandas as pd
import requests
import pandas as pd
from kiteconnect import KiteConnect
import json
from sqlalchemy import create_engine, Column, Integer, Float, String, BigInteger, Table, MetaData
from sqlalchemy.orm import sessionmaker
from sqlalchemy.ext.declarative import declarative_base
import datetime

api_key = "246y3a7zlg83xv2l"

secret_key = "2996e33pjuzh9mmzsoj2y4xs6fdjrmtu"

url = "https://api.kite.trade/instruments/historical/2916865/10minute"
headers = {
    "X-Kite-Version": "3",
    "Authorization": "token 246y3a7zlg83xv2l:0g0113UwBFHPfmsDp6PZnj0wq5IAEFKd"
}
parameters = {
    "from": "2023-10-15 09:15:00",
    "to": "2024-01-15 09:20:00",
    "continuous":0,
    "oi":1
    
}

try:
    response = requests.get(url, headers=headers, params=parameters)
    response.raise_for_status()
    print("Status Code:", response.status_code)
    response_output = response.text
    json = json.loads(response_output)
    print(json)
    
except requests.exceptions.RequestException as e:
    print("Error:", e)

Base = declarative_base()

class PriceData(Base):
    __tablename__ = 'tbl_pricedata_equity_india_10m'
    ntimestamp = Column(BigInteger, primary_key=True)
    nopen = Column(Float)
    nhigh = Column(Float)
    nlow = Column(Float)
    nclose = Column(Float)
    sexchange = Column(String(45))
    nvolume = Column(Float)
    nopen_interest = Column(Float)
    stickername = Column(String(45), primary_key=True)
    sapi = Column(String(45))

# JSON data
data_json =json

# Parse the JSON data
candles = data_json['data']['candles']

# Create a database connection
engine = create_engine('mysql+pymysql://root:Karma100%@localhost/finflo_base_db', echo=True)
Session = sessionmaker(bind=engine)
session = Session()

# Insert the data into the database
for candle in candles:
    timestamp = datetime.datetime.strptime(candle[0], '%Y-%m-%dT%H:%M:%S%z').timestamp()
    open_price = candle[1]
    high_price = candle[2]
    low_price = candle[3]
    close_price = candle[4]
    exchange = "zerodha"
    ticker = "ticker"
    api = "kite"
    volume = candle[5]
    open_interest = candle[6]
    
    price_data = PriceData(
        ntimestamp=int(timestamp),
        nopen=open_price,
        nhigh=high_price,
        nlow=low_price,
        nclose=close_price,
        sexchange=str(exchange),
        nvolume = volume,
        nopen_interest = open_interest,
        stickername=str(ticker),
        sapi=str(api)
    )
    session.add(price_data)

# Commit the changes
session.commit()
