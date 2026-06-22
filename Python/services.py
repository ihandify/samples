import os
import requests
from dotenv import load_dotenv

load_dotenv()

API_URL = os.getenv("API_URL", "")
API_KEY = os.getenv("API_KEY", "")

if not API_URL:
    print("⚠️ Warning: API_URL is not set in .env")
if not API_KEY:
    print("⚠️ Warning: API_KEY is not set in .env")

def generate_scoped_public_key(engines, expires_in_seconds):
    try:
        payload = {}
        if engines is not None:
            payload["engines"] = engines
        if expires_in_seconds is not None:
            payload["expiresInSeconds"] = expires_in_seconds

        headers = {
            "Content-Type": "application/json",
            "X-Api-Key": API_KEY
        }

        response = requests.post(
            f"{API_URL}/plan/auth/generate-scoped-public-key",
            json=payload,
            headers=headers,
            timeout=30
        )

        response.raise_for_status()

        print("✅ Generate Scoped Public Key API Request Successful!")

        return response.json()

    except requests.exceptions.HTTPError as e:
        print("❌ Generate Scoped Public Key API Request Failed!", e.response)

        status_code = e.response.status_code

        print(f"Error Status Code: {status_code}")
        print(f"Error Details: {e.response.text}")

        if status_code == 401:
            print("💡 Hint: API Key is missing or invalid.")
        elif status_code == 403:
            print("💡 Hint: Your IP address is not included in the IP allowlist.")

        return None

    except requests.exceptions.RequestException as e:
        print(
            "No response received from server. "
            "Please check your URL or ensure the backend server is running."
        )
        print(e)
        return None
