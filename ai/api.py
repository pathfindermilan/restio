from fastapi import FastAPI, HTTPException, UploadFile, Form
from pydantic import BaseModel

from typing import Optional

from aria import analyze_image_with_text, analyze_text
from utils import image_to_base64

from fastapi import File
import fitz  # PyMuPDF for reading PDF files

app = FastAPI()

class TextAnalysisRequest(BaseModel):
    text: str

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

@app.post("/analyze-image-text")
async def analyze_image_with_text_endpoint(file: UploadFile):
    """
    Endpoint to analyze an image along with text using the Aria API.
    """
    try:
        # Read the image file and convert it to base64
        image_data = await file.read()
        base64_image = image_to_base64(image_data)

        system_prompt = "Analyze the provided text and return a summary of the user's emotion or intent."

        # Perform the analysis with image and prompt
        result = analyze_image_with_text(base64_image, system_prompt)
        return {"analysis_result": result}
    
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/analyze-pdf-summary")
async def analyze_pdf_summary(file: UploadFile):
    """
    Endpoint to extract text from a PDF and generate a summary using the Aria API.
    """
    try:
        # Read the PDF file content
        pdf_data = await file.read()

        # Use PyMuPDF to extract text from the PDF
        pdf_text = ""
        with fitz.open(stream=pdf_data, filetype="pdf") as pdf_document:
            for page_num in range(pdf_document.page_count):
                page = pdf_document[page_num]
                pdf_text += page.get_text("text")

        # Prompt for summarizing the content
        summary_prompt = (
            "Here is a PDF document. Please provide a clear and concise summary of the main information in the content."
        )

        # Submit the extracted text to the model for analysis with the prompt
        result = analyze_text(f"{summary_prompt}\n\n{pdf_text}")

        return {"summary_result": result}
    
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
