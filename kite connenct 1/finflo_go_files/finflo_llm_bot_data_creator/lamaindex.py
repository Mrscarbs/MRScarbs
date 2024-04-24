from llama_index.core import VectorStoreIndex, SimpleDirectoryReader
import openai
import os
import mysql.connector
from configparser import ConfigParser
import pandas as pd

config = ConfigParser()

config.read("db_config.ini")


password = config["database"]["password"]
host = config["database"]["host"]
db_name = config["database"]["db_name"]
username = config["database"]["username"]
conn = mysql.connector.connect(host=host, password=password, user=username, database=db_name)
mycur = conn.cursor()
mycur.execute("SELECT api_key from api_sys_config where api_id = 2")
key = mycur.fetchall()

df = pd.DataFrame(key, columns=["key"])
print(df)


conn.commit()
conn.close()

os.environ["OPENAI_API_KEY"]=df["key"].iloc[0]
print(df["key"].iloc[0])
openai.api_key = df["key"].iloc[0]

documents = SimpleDirectoryReader("llm_docs").load_data()
documents

index = VectorStoreIndex.from_documents(documents)
index

query_engine = index.as_query_engine(api_key = df["key"].iloc[0])
response = query_engine.query("give me the pb for aapl")

print(response)



