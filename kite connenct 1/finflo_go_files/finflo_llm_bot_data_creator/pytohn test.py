# import yfinance as yf

# # Replace "RELIANCE.NS" with the ticker symbol of Reliance Industries
# ticker = yf.Ticker("RELIANCE.NS")

# # Get the stock information
# stock_info = ticker.info

# # Extract the P/E ratio from the stock information
# pe_ratio = stock_info['trailingPE']

# print(f"Price to Earnings Ratio: {pe_ratio}")

import requests

import requests

ticker = "RELIANCE.NS"
_BASE_URL_ = 'https://query2.finance.yahoo.com'
url = f"{_BASE_URL_}/v7/finance/options/{ticker}"

headers = {
    'User-Agent': 'Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36'
}

response = requests.get(url=url, headers=headers)

print(response)

