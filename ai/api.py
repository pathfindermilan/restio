from fastapi import FastAPI, HTTPException, UploadFile, Form, File
from pydantic import BaseModel
from typing import Optional
from aria import analyze_image_with_text, analyze_text
from utils import image_to_base64
import fitz

app = FastAPI()

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

@app.post("/analyze-text")
async def analyze_text_endpoint(request: TextAnalysisRequest):
    """
    Endpoint to analyze text using the Aria API.
    """
    try:
        result = analyze_text(request.text)
        return {"analysis_result": result}
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
