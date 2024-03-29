import pandas as pd
import requests
import pandas as pd
from kiteconnect import KiteConnect

api_key = "246y3a7zlg83xv2l"

secret_key = "2996e33pjuzh9mmzsoj2y4xs6fdjrmtu"






url = "https://api.kite.trade/instruments/historical/2916865/5minute"
headers = {
    "X-Kite-Version": "3",
    "Authorization": "token 246y3a7zlg83xv2l:CYHuP62zMtkptwyOoB0Mu4jfTwOGm3Cs"
}
parameters = {
    "from": "2023-10-15 09:15:00",
    "to": "2024-01-15 09:20:00",
    "continuous":0,
    "oi":1
    
}

try:
    response = requests.get(url, headers=headers, params=parameters)
    response.raise_for_status()  # Raise an exception for 4xx and 5xx status codes
    print("Status Code:", response.status_code)
    print("Response Data:", response.text)  # Assuming the response is JSON
except requests.exceptions.RequestException as e:
    print("Error:", e)
