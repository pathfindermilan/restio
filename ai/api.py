from fastapi import FastAPI, HTTPException, UploadFile, Form
from pydantic import BaseModel

from typing import Optional

from aria import analyze_image_with_text, analyze_text
from utils import image_to_base64

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
