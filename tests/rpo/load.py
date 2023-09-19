import requests
import random
from datetime import datetime, timezone
import threading
import sys

API_KEY = "marco21Sy27CMBD8Fu6OSYvEudDXoSClQHuNlmQbXBw7zTpFUZR"
API_URL = "https://e14cdb0hje.execute-api.eu-central-1.amazonaws.com/dev/product/v1/accept"

ACCEPT_COMMAND = {
    "specversion": "1.0",
    "type": "Logistic.PCL.Product.Accept.Accept",
    "source": "Logistic.PCL.UP.OMP",
    "subject": "Command",
    "id": "123456789015623456789",
    "time": "2023-02-25T21:15:05.232Z",
    "datacontenttype": "application/json",
    "data": {
        "product": {
            "name": "Poste Delivery Web",
            "id": "1a1a1a",
            "type": "BOX",
            "attributes": []
        },
        "location": {
            "type": "UP",
            "address": "Viale Europa 190",
            "zipcode": "00144",
            "city": "Roma",
            "nation": "Italia",
            "locationCode": "55Y90",
            "attributes": []
        },
        "sender": {
            "name": "Pippo Franco",
            "province": "RM",
            "city": "Roma",
            "address": "Via Nepal 51",
            "zipcode": "00144",
            "attributes": []
        },
        "receiver": {
            "name": "Paolo Bonolis",
            "province": "RM",
            "city": "Roma",
            "address": "Via della Camilluccia 649",
            "cap": "00135",
            "number": "3333456987",
            "email": "paolo.bonolis@gmail.com",
            "note": "Consegna presso portiere",
            "attributes": []
        },
        "timestamp": "2023-02-25T21:15:05.232Z",
        "attributes": []
    }
}


def product_id():
    return random.randint(1, 99999999)


def send_request(idx, command):

    headers = {
        'Content-type': 'application/json',
        'x-api-key': API_KEY,
    }

    response = requests.post(
        f"{API_URL}", json=command,
        headers=headers
    )
    print(f"Sent request #{idx} - {response.status_code} - ID: ", command["id"])


def get_timestamp():
    return datetime.now(timezone.utc).isoformat().replace('+00:00', 'Z')


def generate_random_ids(n):
    unique_ids = set()
    for _ in range(n):
        unique_ids.add(product_id())
    return unique_ids


NUM_OF_REQUESTS = 100
PARALLEL = False
print(sys.argv)

try:
    NUM_OF_REQUESTS = sys.argv[1]
except IndexError:
    pass

try:
    PARALLEL = sys.argv[2] == "parallel"
except IndexError:
    pass

print(f"Sending {NUM_OF_REQUESTS} using parallel {PARALLEL}")

ids = generate_random_ids(NUM_OF_REQUESTS)

for idx, id in enumerate(ids):
    
    command = ACCEPT_COMMAND.copy()
    command["id"] = str(id)
    command["data"]["product"]["id"] = str(id)
    command["time"] = get_timestamp()
    command["data"]["timestamp"] = get_timestamp()
    
    if PARALLEL:
        threading.Thread(target=send_request, args=[idx, command]).start()
    else:
        send_request(idx, command)
