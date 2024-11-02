import base64
import os

from dotenv import load_dotenv

from langchain_openai import ChatOpenAI
from langchain_core.messages import HumanMessage, SystemMessage

load_dotenv()

base_url = 'https://api.rhymes.ai/v1'
api_key = os.getenv('ARIA_API_KEY')

client = ChatOpenAI(
    model="aria",
    api_key=api_key,
    base_url=base_url,
    streaming=False,
)

def analyze_image_with_text(base64_image: str, prompt: str) -> str:
    """
    Envoie une image encodée en base64 et un texte à l'API Aria pour obtenir une analyse en retour.

    :param image_b64: L'image encodée en base64
    :param text: Un texte décrivant l'intention ou le contexte de l'analyse
    :return: La réponse de l'API en texte
    """
    
    # prompt = "Analyze the image and text provided and return a short emotional summary."

    response = client.chat.completions.create(
                model="aria",
                messages=[
                    {
                        "role": "user",
                        "content": [
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
                    }
                ],
                temperature=0.6,
                max_tokens=1024,
                top_p=1,
                stop=["<|im_end|>"]
            )


    return response.choices[0].message.content


def analyze_text(text: str) -> str:
    """
    Envoie uniquement un texte à l'API Aria pour obtenir une analyse en retour.

    :param text: Le texte à analyser
    :return: La réponse de l'API en texte
    """
   
    response = client.chat.completions.create(
                model="aria",
                messages=[
                    {
                        "role": "user",
                        "content": [
                            {
                                "type": "text",
                                "text": f"{prompt}"
                            }
                        ]
                    }
                ],
                temperature=0.6,
                max_tokens=1024,
                top_p=1,
                stop=["<|im_end|>"]
            )


    return response.choices[0].message.content
