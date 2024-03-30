import mysql.connector
from configparser import ConfigParser
import pandas as pd
from enum import Enum


config = ConfigParser()

config.read('db_config.ini')

host = config['database']['host']
password = config['database']['password']
username = config['database']['username']
dbname = config['database']['db_name']

class api_sys_enum(Enum):
     api_key = 0
     secret_key = 1
     api_provider = 2
     access_token = 3
     last_prchase_time = 4
     first_purchase_time = 5
     historical = 6
     instrument_type =7
     api_id = 8




class sys_config_fethcer():

    def __init__(self, api_id, column_to_fetch) -> None:

        self.api_id = api_id
        self.column_to_fetch = column_to_fetch

    def fetcher(self):

        conn = mysql.connector.connect(
                    host=host,
                    user=username,
                    password=password,
                    database=dbname
        )

        mycursor = conn.cursor()

        mycursor.callproc("stp_get_api_config", [self.api_id,])
        for result in mycursor.stored_results():
                
                my_result = result.fetchall()

                result1 = my_result[0]
    
        result = list(result1)
             
        resultdf2 = pd.DataFrame(result, columns = ['details'])

        column = resultdf2['details'].iloc[self.column_to_fetch]
        #print(resultdf2)
        conn.commit()
        conn.close()
        return column
    
    def dataframe_fetcher(self):
        conn = mysql.connector.connect(
                    host=host,
                    user=username,
                    password=password,
                    database=dbname
        )

        mycursor = conn.cursor()

        mycursor.callproc("stp_get_api_config", [self.api_id,])
        for result in mycursor.stored_results():
                
                my_result = result.fetchall()

                result1 = my_result[0]
    
        result = list(result1)
             
        resultdf2 = pd.DataFrame(result, columns = ['details'])
        conn.commit()
        conn.close()
        return resultdf2
         
    

obg1 =sys_config_fethcer(1, api_sys_enum.api_key.value)
fetch =obg1.fetcher()
dataframe = obg1.dataframe_fetcher()
print(fetch)
print(dataframe)

        
        




    