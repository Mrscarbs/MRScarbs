import pandas as pd
import requests
import hashlib
from kiteconnect import KiteConnect

api_key = "246y3a7zlg83xv2l"

secret_key = "2996e33pjuzh9mmzsoj2y4xs6fdjrmtu"

request_token ="UyQ3UvALdRiRSCbH2N4H1Xj2gDeAepzk"

kite = KiteConnect(api_key=api_key)

print(kite.login_url())



def generate_checksum(api_key, request_token, api_secret):
    data_to_hash = api_key.encode() + request_token.encode() + api_secret.encode()
    checksum = hashlib.sha256(data_to_hash).hexdigest()
    return checksum



checksum = generate_checksum(api_key, request_token, secret_key)
print("Checksum:", checksum)


url = "https://api.kite.trade/session/token"
headers = {
    "X-Kite-Version": "3"
}
data = {
    "api_key": api_key,
    "request_token": request_token,
    "checksum": checksum
}
response = requests.post(url, headers=headers, data=data)

print(response.json)
print(response.text)
