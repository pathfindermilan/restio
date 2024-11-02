from openai import OpenAI
from dotenv import load_dotenv
from langchain_openai import ChatOpenAI
from langchain_core.messages import HumanMessage, SystemMessage

import os

load_dotenv()

base_url = 'https://api.rhymes.ai/v1'
api_key = os.getenv('ARIA_API_KEY')

chat = ChatOpenAI(
    model="aria",
    api_key=api_key,
    base_url=base_url,
    streaming=False,
)

system_prompt = input("-> Enter the system prompt: \n")
human_message = input("\n-> Enter the human message: \n")

response = chat.invoke([
    SystemMessage(content=system_prompt),
    HumanMessage(content=human_message)
])

print("\n-> AI Response:\n", response.content)