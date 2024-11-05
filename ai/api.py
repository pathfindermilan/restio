from pydantic import BaseModel
from utils import image_to_base64
from typing import Dict, Optional
from format import ContentGenerator
from langchain_openai import ChatOpenAI
from aria import analyze_image_with_text, analyze_text
from fastapi import FastAPI, HTTPException, UploadFile, Form, File

import os
import fitz
import random

app = FastAPI()

class UserData(BaseModel):
    name: str
    age: int
    image_summary: str
    document_summary: str
    user_transcript: str
    content_type: str
    feeling_level: str

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

class TextAnalysisRequest(BaseModel):
    text: str

@app.post("/describe-image")
async def describe_image_endpoint(image: UploadFile = File(...)):
    """
    Endpoint to describe an image using the Aria API.
    """
    try:
        image_data = await image.read()
        base64_image = image_to_base64(image_data)
        system_prompt = "Analyze the provided image and return a description of its content."
        result = analyze_image_with_text(base64_image, system_prompt)
        return {"image_summary": result}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/describe-document")
async def analyze_pdf_summary(document: UploadFile = File(...)):
    """
    Endpoint to extract text from a PDF and generate a summary using the Aria API.
    """
    try:
        pdf_data = await document.read()

        pdf_text = ""
        with fitz.open(stream=pdf_data, filetype="pdf") as pdf_document:
            for page_num in range(pdf_document.page_count):
                page = pdf_document[page_num]
                pdf_text += page.get_text("text")

        summary_prompt = (
            "Here is a PDF document. Please provide a clear and concise summary "
            "of the main information in the content."
        )

        result = analyze_text(f"{summary_prompt}\n\n{pdf_text}")
        return {"document_summary": result}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/generate-answer")
async def process_user_request(data: UserData) -> Dict[str, str]:
    try:
        content_generator = ContentGenerator(
            client=client,
            allegro_token=os.getenv("ALLEGRO_API_KEY")
        )

        response = content_generator.generate_content(data)

        return {
            "ai_summary": response["content"]
        }
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Error generating prompt: {str(e)}")
