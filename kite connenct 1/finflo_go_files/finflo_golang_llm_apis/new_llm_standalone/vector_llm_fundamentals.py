from langchain.indexes import VectorstoreIndexCreator
import os
from langchain.document_loaders import TextLoader
import mysql.connector
import sys

conn = mysql.connector.connect(
        host="localhost",
        user="root",
        password="Karma100%",
        database="finflo_base_db"
    )
mycur = conn.cursor()
mycur.callproc("stp_get_api_config", [7,])

for result in mycur.stored_results():
    rest=result.fetchall()
    rest_final = rest[0]


conn.commit()
conn.close()
api_key = rest_final[0]


os.environ['OPENAI_API_KEY'] = api_key #"sk-proj-qp6oyXYelSJ5dJCWMnEtT3BlbkFJ8pd3TUT61oU8CT1cKJW1"
loader = TextLoader(r'C:\Users\karma\OneDrive\Documents\finfloapps\MRScarbs\kite connenct 1\finflo_go_files\finflo_llm_bot_data_creator\llm_doc.txt')
loader2 = TextLoader(r'C:\Users\karma\OneDrive\Documents\finfloapps\MRScarbs\kite connenct 1\finflo_go_files\finflo_llm_bot_data_creator\llm_cashflow.txt')
loader3 = TextLoader(r'C:\Users\karma\OneDrive\Documents\finfloapps\MRScarbs\kite connenct 1\finflo_go_files\finflo_llm_bot_data_creator\llm_companies_contex.txt')


index = VectorstoreIndexCreator().from_loaders([loader, loader2,loader3])
def llm_cashflow_query(query):
    index_res = index.query(str(query))
    return index_res
query = str(sys.argv[1])
index_res=llm_cashflow_query(query)
print(index_res)




