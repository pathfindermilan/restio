import os
import requests

from time import sleep
from dotenv import load_dotenv

load_dotenv()

def generate_video(prompt, token):
    url = "https://api.rhymes.ai/v1/generateVideoSyn"
    headers = {
        "Authorization": f"Bearer {token}",
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
        return f"An error occurred: {str(e)}"

def query_video_status(token, request_id):
    url = "https://api.rhymes.ai/v1/videoQuery"
    headers = {
        "Authorization": f"Bearer {token}",
    }
    params = {
        "requestId": request_id
    }

    try:
        response = requests.get(url, headers=headers, params=params)
        response.raise_for_status()
        return response.json()
    except requests.exceptions.RequestException as e:
        return f"An error occurred: {str(e)}"

token = os.getenv("ALLEGRO_API_KEY")
prompt = input("\n-> Enter the prompt: \n")

response_data = generate_video(prompt, token)
print(response_data)

counter = 0
while True:
    sleep(20)
    print((counter + 1) * 20, "s")
    video_data = query_video_status(token, response_data.get("data")).get("data")
    if video_data:
        print(video_data)
        break
    counter = counter + 1


