from fastapi import FastAPI, File, UploadFile
from openai import OpenAI

from dotenv import load_dotenv

from langchain_openai import ChatOpenAI
from langchain_core.messages import HumanMessage, SystemMessage

import os
import io
from PIL import Image

# Charger les variables d'environnement
load_dotenv()

# Configuration des paramètres API
base_url = 'https://api.rhymes.ai/v1'
api_key = os.getenv('ARIA_API_KEY')

# Initialiser FastAPI
app = FastAPI()

# Initialiser l'API ChatOpenAI avec le modèle Aria
chat = ChatOpenAI(
    model="aria",
    api_key=api_key,
    base_url=base_url,
    streaming=False,
)

# Fonction pour générer une description émotionnelle à partir de l'image
def generate_image_summary(image: Image.Image) -> str:
    system_prompt = "Analyze the image and provide a short summary of the user's emotions based on the image."
    human_message = "Here is the image."
    
    # Invocation de l'API avec les messages de système et d'utilisateur
    response = chat.invoke([
        SystemMessage(content=system_prompt),
        HumanMessage(content=human_message)
    ])
    
    return response.content

# Endpoint FastAPI pour analyser une image et générer un résumé émotionnel
@app.post("/analyze-emotion/")
async def analyze_emotion(image: UploadFile = File(...)) -> dict:
    # Lire le fichier image
    image_data = await image.read()
    img = Image.open(io.BytesIO(image_data))

    # Obtenir le résumé émotionnel en utilisant la fonction de génération
    summary = generate_image_summary(img)

    # Retourner la réponse en JSON
    return {"summary": summary}
