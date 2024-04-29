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
mycur.callproc("stp_get_api_config", [7,])

for result in mycur.stored_results():
    my_result = result.fetchall()
    res = my_result[0]
    
api_key = res[0]

os.environ["OPENAI_API_KEY"] = api_key
loader = TextLoader(r'C:\Users\karma\OneDrive\Documents\finfloapps\MRScarbs\kite connenct 1\finflo_go_files\finflo_llm_bot_data_creator\economic_llm\economic_llm.txt', encoding='utf-8')

index = VectorstoreIndexCreator().from_loaders([loader])

question = str(sys.argv[1])

answer = index.query(question)

print(answer)
