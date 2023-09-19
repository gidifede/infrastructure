import json
from datetime import datetime, timezone
from jsonschema import RefResolver, Draft7Validator
from cloudevents.conversion import to_structured
from cloudevents.http import CloudEvent
import requests


def generate_cloudevent(source, type, data):
    attributes = {
        "specversion": "1.0",
        "type": type,
        "source": source,
        "time": datetime.now(timezone.utc).isoformat().replace('+00:00', 'Z'),
        "datacontenttype": "application/json",
        "subject": "Command",
    }
    event = CloudEvent(attributes, data)

    # Creates the HTTP request representation of the CloudEvent in structured content mode
    headers, body = to_structured(event)
    return headers, body


class BaseCommand:

    common_schema_path = "configuration/JSON/JSONSCHEMA/common.json"
    command_schema_path = None
    payload = None
    api_path = None

    def validate(self):
        pass
        # with open(self.command_schema_path) as f:
        #     command = json.load(f)
        # with open(self.common_schema_path) as f:
        #     common_schema = json.load(f)
        # schema_store = {
        #     command['$id']: command,
        #     common_schema['$id']: common_schema,
        # }

        # resolver = RefResolver.from_schema(command, store=schema_store)
        # validator = Draft7Validator(command, resolver=resolver)
        # validator.validate(self.payload)

    def to_json(self):
        return json.loads(self.body)

    def send(self, base_api_url, api_key=None):
        headers = {
                'Content-type': 'application/json'
            }
        if api_key is not None:
            headers['x-api-key'] = api_key 
        
        headers.update(self.headers)
        response = requests.post(
            f"{base_api_url}{self.api_path}", json=self.to_json(),
            headers=headers
        )
        return response
