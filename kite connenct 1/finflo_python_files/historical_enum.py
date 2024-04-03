import mysql.connector
import datetime
from datetime import timedelta

class TimeDiff:
    def __init__(self, interval, api_id):
        self.interval = interval
        self.api_id = api_id

    def get_db_details(self):
        # Connect to the database
        conn = mysql.connector.connect(host="localhost", password="Karma100%", user="root", database="finflo_base_db")
        mycursor = conn.cursor()

        mycursor.callproc("stp_get_limits", (self.api_id, self.interval))

        for result in mycursor.stored_results():
            myresult = result.fetchall()
            result1 = myresult[0]
        
        self.result = list(result1)
        self.result = self.result[0]
        print(self.result)

        mycursor.close()
        conn.close()

    def last_time(self):
        
        current_time = datetime.datetime.now()
        last_date = current_time - timedelta(minutes=self.result)
        return last_date
    
# time_diff = TimeDiff('minute', 1)
# time_diff.get_db_details()
# last_time = time_diff.last_time()
# print(last_time)