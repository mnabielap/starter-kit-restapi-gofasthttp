import sys
import os
sys.path.append(os.path.abspath(os.path.dirname(__file__)))
from utils import send_and_print, BASE_URL, save_config

url = f"{BASE_URL}/auth/login"

payload = {
    "email": "admin@example.com",
    "password": "password123"
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

if response and response.status_code == 200:
    data = response.json()
    save_config("access_token", data['tokens']['access']['token'])
    save_config("refresh_token", data['tokens']['refresh']['token'])
    save_config("user_id", str(data['user']['id']))
    print("\n[INFO] Tokens updated in secrets.json")