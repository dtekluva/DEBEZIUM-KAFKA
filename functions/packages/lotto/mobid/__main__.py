"""
Mobid Tracker list serverless functions
"""

from typing import TypedDict, Any
import requests


class Args(TypedDict):
    limit: str
    page: str


def main(args: Args) -> Any:
    limit = args.get("limit", 10)
    page = args.get("page", 1)

    # Call the internal api service
    response_data = requests.get(
        f"http://24.199.101.18:4200/mobid-trackers?limit={limit}&page={page}"
    )
    print(response_data.json())

    return {"body": response_data.json()}
