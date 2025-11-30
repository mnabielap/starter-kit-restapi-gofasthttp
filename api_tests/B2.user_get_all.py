import sys
import os
sys.path.append(os.path.abspath(os.path.dirname(__file__)))
from utils import send_and_print, BASE_URL, load_config

url = f"{BASE_URL}/users?page=1&limit=5"
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