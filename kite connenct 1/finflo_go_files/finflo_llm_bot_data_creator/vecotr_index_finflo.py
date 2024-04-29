from langchain.indexes import VectorstoreIndexCreator
import os
from langchain.document_loaders import TextLoader
os.environ['OPENAI_API_KEY'] = "sk-proj-qp6oyXYelSJ5dJCWMnEtT3BlbkFJ8pd3TUT61oU8CT1cKJW1"
loader = TextLoader('llm_doc.txt')
loader2 = TextLoader('llm_cashflow.txt')
loader3 = TextLoader("llm_companies_contex.txt")

index = VectorstoreIndexCreator().from_loaders([loader, loader2,loader3])
index.query("give me detail for google")
def query_index(request_text):


    result=index.query(str(request_text))

    print(result)

query_index("give me the pe and pb for microsoft")