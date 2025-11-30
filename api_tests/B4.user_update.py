import sys
import os
sys.path.append(os.path.abspath(os.path.dirname(__file__)))
from utils import send_and_print, BASE_URL, load_config

target_id = load_config("target_user_id")
if not target_id:
    print("No target user to update. Run user_create.py first.")
    exit()

url = f"{BASE_URL}/users/{target_id}"
token = load_config("access_token")

payload = {
    "name": "John Updated",
    "email": "john.updated@example.com"
}

headers = {
    "Content-Type": "application/json",
    "Authorization": f"Bearer {token}"
}

send_and_print(
    url=url,
    method="PATCH",
    headers=headers,
    body=payload,
    output_file=f"{os.path.splitext(os.path.basename(__file__))[0]}.json",
    print_pretty_response=True
)