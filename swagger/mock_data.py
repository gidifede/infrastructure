
from swagger.builder import BackboneSwagger
from infrastructure.constructs.utils import generate_random_data
from faker import Faker
import json
import os


localpath = os.path.dirname(__file__)

fake = Faker()
bb_swagger = BackboneSwagger.load("bb-swagger", "v1")

mock_responses = {}

for api in bb_swagger.apis:
    if api.response_schema:
        print(api.path, generate_random_data(api.response_schema, fake))
        mock_responses.update({api.path: generate_random_data(api.response_schema, fake)})
        # Directly from dictionary
        swagger_file = json.dumps(mock_responses, indent=2)
        filename = f"{localpath}/mocks.json"
        with open(filename, 'w') as outfile:
            outfile.write(swagger_file)
        