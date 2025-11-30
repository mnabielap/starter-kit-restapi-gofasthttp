import sys
import os
sys.path.append(os.path.abspath(os.path.dirname(__file__)))
from utils import send_and_print, BASE_URL, save_config

url = f"{BASE_URL}/auth/register"

payload = {
    "name": "Super Admin",
    "email": "admin@example.com",
    "password": "password123",
    "role": "admin" 
}

headers = {
    "Content-Type": "application/json"
}

response = send_and_print(
    url=url,
    method="POST",
    headers=headers,
    body=payload,
    output_file=f"{os.path.splitext(os.path.basename(__file__))[0]}.json",
    print_pretty_response=True
)

if response and response.status_code == 201:
    data = response.json()
    # Save tokens for next scripts
    save_config("access_token", data['tokens']['access']['token'])
    save_config("refresh_token", data['tokens']['refresh']['token'])
    save_config("user_id", str(data['user']['id']))
    print("\n[INFO] Tokens and User ID saved to secrets.json")