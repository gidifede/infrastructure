import re
import json
from datetime import datetime, timezone
from swagger.builder import BackboneSwagger



def get_timestamp():
    return  datetime.now(timezone.utc).isoformat().replace('+00:00', 'Z')
    

def get_version(version_file_path: str):
    regex = r"\d*[.]\d*[.]\d*[-]*SNAPSHOT*"
    try:
        with open(version_file_path) as f:
            test_str= f.read()
    except FileNotFoundError as fnfe:
        print(fnfe)
        return None
    
    matches = re.finditer(regex,test_str, re.MULTILINE)
    versions = [m for m in matches]
    for matchNum, match in enumerate(matches, start = 1):
        print ("Match {matchNum} was found at {start}-{end}:{match}".format(matchNum = matchNum, start = match.start(), end = match.end(), match = match.group()))

    if (len(versions)>0) and (versions[0].group() !=""):
        return versions[0].group()
    else:
        return None

        
def generate_random_data(schema, fake):
    
    if response_schema := schema.get("$ref", None):
        bb_swagger = BackboneSwagger.load("bb-swagger", "v1")
        swagger = json.loads(bb_swagger.swagger_content)
        response_schema = response_schema.split("/")[-1]
        response_schema_json = swagger["components"]["schemas"].get(response_schema, None)
        schema = response_schema_json
    
    if schema['type'] == 'object':
        properties = schema.get('properties', {})
        data = {}
        for prop, prop_schema in properties.items():
            data[prop] = generate_random_data(prop_schema, fake)
        return data
    elif schema['type'] == 'array':
        items_schema = schema['items']
        min_items = schema.get('minItems', 0)
        max_items = schema.get('maxItems', 5)
        num_items = fake.random_int(min=min_items, max=max_items)
        return [generate_random_data(items_schema, fake) for _ in range(num_items)]
    elif schema['type'] == 'string':
        if format := schema.get("format", False):
            if format == "date-time":
                return get_timestamp()
        return fake.word()
    elif schema['type'] == 'number':
        return fake.random_number()
    elif schema['type'] == 'integer':
        return fake.random_int()
    elif schema['type'] == 'boolean':
        return fake.boolean()
    else:
        return None