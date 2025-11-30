import sys
import os
sys.path.append(os.path.abspath(os.path.dirname(__file__)))
from utils import send_and_print, BASE_URL, load_config, save_config

url = f"{BASE_URL}/users"
token = load_config("access_token")

payload = {
    "name": "John Doe",
    "email": "john.doe@example.com",
    "password": "password123",
    "role": "user"
}

headers = {
    "Content-Type": "application/json",
    "Authorization": f"Bearer {token}"
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
    # Save this specific user ID to test update/delete later
    save_config("target_user_id", str(data['id']))
    print("\n[INFO] Target User ID saved to secrets.json")