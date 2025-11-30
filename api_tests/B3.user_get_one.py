import sys
import os
sys.path.append(os.path.abspath(os.path.dirname(__file__)))
from utils import send_and_print, BASE_URL, load_config

# We try to get the user created in user_create.py, or fallback to the logged in user
target_id = load_config("target_user_id") or load_config("user_id")
url = f"{BASE_URL}/users/{target_id}"
token = load_config("access_token")

headers = {
    "Content-Type": "application/json",
    "Authorization": f"Bearer {token}"
}

send_and_print(
    url=url,
    method="GET",
    headers=headers,
    output_file=f"{os.path.splitext(os.path.basename(__file__))[0]}.json",
    print_pretty_response=True
)