import sys
import os
sys.path.append(os.path.abspath(os.path.dirname(__file__)))
from utils import send_and_print, BASE_URL, load_config

url = f"{BASE_URL}/auth/logout"
refresh_token = load_config("refresh_token")

payload = {
    "refreshToken": refresh_token
}

headers = {
    "Content-Type": "application/json"
}

send_and_print(
    url=url,
    method="POST",
    headers=headers,
    body=payload,
    output_file=f"{os.path.splitext(os.path.basename(__file__))[0]}.json",
    print_pretty_response=True
)