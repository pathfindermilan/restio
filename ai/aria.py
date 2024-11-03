import base64
import os
from dotenv import load_dotenv
from langchain_openai import ChatOpenAI
from langchain_core.messages import HumanMessage, SystemMessage
from typing import List, Dict, Any

load_dotenv()
base_url = 'https://api.rhymes.ai/v1/'
api_key = os.getenv('ARIA_API_KEY')
client = ChatOpenAI(
    model="aria",
    api_key=api_key,
    base_url=base_url,
    streaming=False,
    temperature=0.6,
    max_tokens=1024
)

def analyze_image_with_text(base64_image: str, prompt: str) -> str:
    """
    Envoie une image encodée en base64 et un texte à l'API Aria pour obtenir une analyse en retour.
    :param base64_image: L'image encodée en base64
    :param prompt: Un texte décrivant l'intention ou le contexte de l'analyse
    :return: La réponse de l'API en texte
    """
    messages = [
        HumanMessage(
            content=[
                {
                    "type": "image_url",
                    "image_url": {
                        "url": f"data:image/jpeg;base64,{base64_image}"
                    }
                },
                {
                    "type": "text",
                    "text": f"<image>\n{prompt}"
                }
            ]
        )
    ]

    response = client.invoke(messages)
    return response.content

def analyze_text(text: str) -> str:
    """
    Envoie uniquement un texte à l'API Aria pour obtenir une analyse en retour.
    :param text: Le texte à analyser
    :return: La réponse de l'API en texte
    """
    messages = [
        HumanMessage(
            content=[
                {
                    "type": "text",
                    "text": text
                }
            ]
        )
    ]

    response = client.invoke(messages)
    return response.content
