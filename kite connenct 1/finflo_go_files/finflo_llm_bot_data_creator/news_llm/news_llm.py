from langchain.indexes import VectorstoreIndexCreator
from langchain.document_loaders import TextLoader
import mysql.connector
import sys
import os


conn = mysql.connector.connect(
        host="localhost",
        user="root",
        password="Karma100%",
        database="finflo_base_db"
    )


mycur = conn.cursor()
mycur.callproc("stp_get_api_config", [6,])

for result in mycur.stored_results():
    my_result = result.fetchall()
    res = my_result[0]
    
api_key = res[0]

os.environ["OPENAI_API_KEY"] = api_key
loader = TextLoader("news_docs.txt", encoding='utf-8')

index = VectorstoreIndexCreator().from_loaders([loader])

question = str(sys.argv[1])

answer = index.query(question)

print(answer)
