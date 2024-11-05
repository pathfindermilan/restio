from typing import Optional, List, Dict, Any
from pydantic import BaseModel
from datetime import datetime
import os
import requests
from time import sleep
from langchain.schema import HumanMessage

class ContentGenerator:
    def __init__(self, client, allegro_token):
        self.client = client
        self.allegro_token = allegro_token

    class UserData(BaseModel):
        name: str
        age: int
        image_summary: str
        document_summary: str
        user_transcript: str
        content_type: str
        feeling_level: str

    def analyze_text(self, text: str) -> str:
        """Call Aria API for text content generation."""
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
        response = self.client.invoke(messages)
        return response.content

    def generate_video(self, prompt: str) -> Dict:
        """Generate video using Allegro API."""
        url = "https://api.rhymes.ai/v1/generateVideoSyn"
        headers = {
            "Authorization": f"Bearer {self.allegro_token}",
            "Content-Type": "application/json"
        }
        data = {
            "refined_prompt": prompt,
            "num_step": 100,
            "cfg_scale": 7.5,
            "user_prompt": prompt,
            "rand_seed": 12345
        }
        try:
            response = requests.post(url, headers=headers, json=data)
            response.raise_for_status()
            return response.json()
        except requests.exceptions.RequestException as e:
            raise Exception(f"Video generation error: {str(e)}")

    def query_video_status(self, request_id: str) -> Dict:
        """Query video generation status."""
        url = "https://api.rhymes.ai/v1/videoQuery"
        headers = {
            "Authorization": f"Bearer {self.allegro_token}",
        }
        params = {
            "requestId": request_id
        }
        try:
            response = requests.get(url, headers=headers, params=params)
            response.raise_for_status()
            return response.json()
        except requests.exceptions.RequestException as e:
            raise Exception(f"Video status query error: {str(e)}")

    def generate_content(self, data: UserData) -> Dict:
        def get_base_prompt(content_type: str) -> str:
            """Get base prompt based on content type."""
            prompts = {
                "joke": f"""Generate a humorous and uplifting story for {data.name} who is {data.age} years old and feeling {data.feeling_level}.
                Focus on these elements:
                - Create an engaging story that brings joy and laughter
                - Use age-appropriate humor and references
                - Include elements from their interests: {data.user_transcript}
                - Keep the tone light and positive
                - Make sure the humor is respectful and uplifting
                - Include a clever punchline or resolution""",

                "exercise": f"""Create an engaging exercise routine for {data.name} who is {data.age} years old and feeling {data.feeling_level}.
                Routine requirements:
                - Focus on mood-lifting exercises
                - Keep exercises safe and achievable
                - Include clear instructions for each movement
                - Provide modifications if needed
                - Add encouraging messages throughout
                Current situation: {data.user_transcript}""",

                "speech": f"""Compose an inspiring motivational speech for {data.name} who is {data.age} years old and feeling {data.feeling_level}.
                Speech elements:
                - Address their current situation: {data.user_transcript}
                - Use age-appropriate language and examples
                - Include personal development insights
                - Build confidence and motivation
                - Provide actionable steps
                - End with an inspiring call to action""",

                "video": f"""Create a short 5-second uplifting video scene for {data.name} who is {data.age} years old and feeling {data.feeling_level}.
                Scene requirements:
                - Visual mood: bright and positive
                - Scene action: dynamic and engaging
                - Color palette: vibrant and uplifting
                - Movement: smooth and purposeful
                - Context: {data.user_transcript}
                Make the scene visually represent a transformation from {data.feeling_level} to a positive emotional state."""
            }
            return prompts.get(content_type, "Invalid content type")

        try:
            base_prompt = get_base_prompt(data.content_type)

            if data.image_summary:
                base_prompt += f"\nVisual context to consider: {data.image_summary}"
            if data.document_summary:
                base_prompt += f"\nAdditional context: {data.document_summary}"

            if data.content_type == "video":
                video_response = self.generate_video(base_prompt)
                request_id = video_response.get("data")

                counter = 0
                video_data = None

                while counter < 30:
                    sleep(20)
                    video_status = self.query_video_status(request_id)
                    video_data = video_status.get("data")
                    if video_data:
                        break
                    counter += 1

                return {
                    "content_type": "video",
                    "prompt": base_prompt,
                    "content": video_data,
                    "status": "completed" if video_data else "timeout",
                    "generated_at": datetime.now().isoformat(),
                }
            else:
                generated_content = self.analyze_text(base_prompt)

                return {
                    "content_type": data.content_type,
                    "prompt": base_prompt,
                    "content": generated_content,
                    "generated_at": datetime.now().isoformat(),
                }

        except Exception as e:
            return {
                "content": "",
                "content_type": data.content_type,
                "generated_at": datetime.now().isoformat(),
            }
