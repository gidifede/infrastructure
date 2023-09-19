from typing import List
from apispec.exceptions import DuplicateComponentNameError
from apispec import APISpec
import json
import os
import numpy as np


localpath = os.path.dirname(__file__)


class API:
    def __init__(self, path, api, response_schema=None):
        self.path = path
        self.api = api
        self.methods = list(api.keys())
        self.tags = []
        for method in self.methods:
            self.tags.extend(api[method]["tags"])
        self.response_schema = response_schema
    
    @property
    def response_mock(self):
        with open(f'{localpath}/mocks.json') as json_file:
            data = json.load(json_file)
            return data[self.path]


class BackboneSwagger:
    def __init__(self, name, version, config_folder="jsonconfig", generateMap=False, swagger=None):

        self.name = name
        self.version = version
        self.config_folder = config_folder
        self.generateMap = generateMap
        self.swagger_content = None
        self.commandApiMap = None

        if swagger:
            self.swagger_content = swagger
            return

        self.config = {}
        res = []
        dirname = os.path.dirname(__file__)
        
        config_folder = os.path.join(dirname, self.config_folder)
        for file_path in os.listdir(config_folder):
            # check if current file_path is a file
            if os.path.isfile(os.path.join(config_folder, file_path)):
                # add filename to list
                res.append(os.path.join(config_folder, file_path))      

        for idx, config_file in enumerate(res):
            with open(f"{config_file}") as f:
                self.config = json.load(f)
            if idx == 0:
                self.spec = APISpec(
                    title=self.config["config"]["title"],
                    version=self.config["config"]["version"],
                    openapi_version="3.0.2",
                    info=dict(description=self.config["config"]["description"]),
                )

                api_key_scheme = {"type": "apiKey",
                                "in": "header", "name": "x-api-key"}
                self.spec.components.security_scheme("ApiKeyAuth", api_key_scheme)

            self.__load_ref_schemas()

            for aggregate, apis in self.config["apis"].items():
                for api in apis:

                    self.__load_command_schema(api=api, aggregate=aggregate)

                    self.__add_path(aggregate=aggregate, api=api)

        self.swagger_content = json.dumps(self.spec.to_dict())

        # print(22222222, self.swagger_content)

        if self.generateMap:
            commandApiMap = {}
            for aggregate, apis in self.config["apis"].items():
                aggregateMap = {}
                for api in apis:          
                    if api['json_schema_location'] is not None:   
                        apiName = api['path'].split('/')[-1].lower()
                        command = api['json_schema_location'].split('/')[-1].split('.')[0]
                        aggregateMap[apiName] = command
                commandApiMap[aggregate] = aggregateMap

            self.commandApiMap = commandApiMap


    def __load_ref_schemas(self):
        for ref in self.config["config"]["json_schema_refs"]:
            with open(f"{localpath}/{ref}") as common_schemas:
                schema = json.loads(common_schemas.read())

            for key, common_struct in schema.get("Definitions", {}).items():
                try:
                    self.spec.components.schema(
                        key, common_struct
                    )
                except DuplicateComponentNameError as e:
                    pass

    def __load_schema(self, name, json_schema_location):
        with open(f"{localpath}/{json_schema_location}") as schema:

            schema_str = schema.read()
            schema_str = schema_str.replace(
                "https://example.com/schemas/common#/Definitions/", "#/components/schemas/")

            schema = json.loads(schema_str)
            schema.pop("$id", None)
            schema.pop("$schema", None)
            self.spec.components.schema(
                name, schema
            )

    def __load_command_schema(self, api, aggregate):
        name = api["name"]
        json_schema_location = api["json_schema_location"]

        if json_schema_location:

            # loads cloudevents schema including into data field the schema related to the command
            aggregate = aggregate.title()
            schemaname = f"{name}{aggregate}"
            # self.spec.components.schema(
            #     f"{schemaname}Payload", {
            #         "type": "object",
            #         "allOf": [{"$ref": "https://raw.githubusercontent.com/cloudevents/spec/v1.0.1/spec.json"}],
            #         "properties": {
            #                 "data": {"$ref": f"#/components/schemas/{schemaname}"}
            #         }
            #     }
            # )

            if isinstance(json_schema_location, dict):
                schemas = json_schema_location["anyOf"]
                for idx, schema in enumerate(schemas):
                    split = schema.split('/')
                    jsonschema_file_name = split[-1].split('.')[0]
                                        
                    self.spec.components.schema(



                        f"{schemaname}{jsonschema_file_name}Payload", {
                            "type": "object",
                            "allOf": [{"$ref": "https://raw.githubusercontent.com/cloudevents/spec/v1.0.1/spec.json"}],
                            "properties": {
                                    "data": {"$ref": f"#/components/schemas/{schemaname}{jsonschema_file_name}"}
                            }
                        }
                    )
                    
                    self.__load_schema(f"{schemaname}{jsonschema_file_name}", schema)
            else:
                self.spec.components.schema(
                    f"{schemaname}Payload", {
                        "type": "object",
                        "allOf": [{"$ref": "https://raw.githubusercontent.com/cloudevents/spec/v1.0.1/spec.json"}],
                        "properties": {
                                "data": {"$ref": f"#/components/schemas/{schemaname}"}
                        }
                    }
                )
                self.__load_schema(schemaname, json_schema_location)

    def __build_swagger_responses(self, name, responses):
        """
        Converts a response configuration in a form

        {
            "code": "202",
            "content_type": "text/plain",
            "content_json_schema_location": null,
            "content_schema_type": "string",
            "description": "Command accepted"
        },

        to a swagger response

        "202": {
            "content": {
              "text/plain": {
                "schema": {
                  "type": "string"
                }
              }
            },
            "description": "Command accepted"
          }

        """
        swagger_responses = {}
        for r in responses:
            code = r["code"]
            content_type = r["content_type"]
            content_json_schema_location = r["content_json_schema_location"]
            content_schema_type = r["content_schema_type"]
            description = r["description"]

            response_schema_name = f"{code}{name}Response"
            if content_json_schema_location:
                self.__load_schema(response_schema_name,
                                   content_json_schema_location)
                response_schema = {
                    "$ref": f"#/components/schemas/{response_schema_name}"
                }
            else:
                response_schema = {
                    "type": content_schema_type
                }

            swagger_responses.update({
                code: {
                    "content": {
                        content_type: {
                            "schema": response_schema
                        }
                    },
                    "description": description
                }
            }
            )

        return swagger_responses

    def __add_path(self, aggregate, api):
        name = api["name"]
        name = name.replace(" ", "_")
        path = api["path"]
        aggregate = aggregate.title()
        schemaname = f"{name}{aggregate}"
        method = api["method"]
        parameters = api.get("parameters", [])
        description = api["description"]
        responses = api["responses"]
        json_schema_location = api["json_schema_location"]

        swagger_responses = self.__build_swagger_responses(name, responses)

        request_params = {}
        if method == "post":
            
            if isinstance(json_schema_location, dict):
                
                schemas = json_schema_location["anyOf"]
                schema_names = []
                for schema in schemas:
                    split = schema.split('/')
                    jsonschema_file_name = split[-1].split('.')[0]
                    schema_names.append(f"{schemaname}{jsonschema_file_name}")
                schema_config = {"schema": {"anyOf": schema_names}}
            else:
                schema_config =  {"schema": f"{schemaname}Payload"}
            
            swagger_parameters = []
            for p in parameters:
                swagger_parameters.append(
                    {
                        "in": p.get("in", "query"),
                        "name": p["name"],
                        "schema": {"type": p["type"]},
                        "description": p.get("description"),
                    }
                )
            request_params = dict(
                parameters = swagger_parameters,
                requestBody={
                    "content": {"application/json": schema_config},
                    "description": description
                    
                },
                
                responses=swagger_responses,
                tags=[aggregate],
                description =  description,
                security=[{"ApiKeyAuth": []}]
            )

        else:  # method = get

            swagger_parameters = []

            for p in parameters:
                swagger_parameters.append(
                    {
                        "in": p.get("in", "query"),
                        "name": p["name"],
                        "schema": {"type": p["type"]},
                        "description": p.get("description"),
                    }
                )
            request_params = {
                "parameters": swagger_parameters,
                "responses": swagger_responses,
                "tags": [aggregate],
                "description" : description,

                "security": [{"ApiKeyAuth": []}]
            }

        self.spec.path(
            path=path,
            operations={
                f"{method}": request_params
            },
        )

    def save_swagger(self):
        """
        save swagger file in JSON format
        """
        swagger_file = json.dumps(self.spec.to_dict(), indent=2)
        filename = f"{localpath}/{self.name}/{self.version}/swagger.json"
        os.makedirs(os.path.dirname(filename), exist_ok=True)

        with open(filename, mode="w") as swagger_json:
            swagger_json.write(swagger_file)

    def save_commandApiMap(self):
        """
        save command-api map file in JSON format
        """
        if self.generateMap:
            for aggregate, map in self.commandApiMap.items():
                fileContent = json.dumps(map, indent=2)
                filename = f"{localpath}/{self.name}/{self.version}/{aggregate}-commandApiMap.json"
                os.makedirs(os.path.dirname(filename), exist_ok=True)

                with open(filename, mode="w") as jsonFile:
                    jsonFile.write(fileContent)


    @classmethod
    def load(cls, name, version):
        """
        loads swagger file
        """
        filename = f"{localpath}/{name}/{version}/swagger.json"
        with open(f"{filename}") as swagger_file:
            swagger = swagger_file.read()

        return cls(name, version, None, False, swagger)

    @property
    def apis(self):
        swagger = json.loads(self.swagger_content)
        apis = []
        for path, api in swagger["paths"].items():
            
            response_schema = api.get("get", {}).get("responses", {}).get("200", {}).get("content", {}).get("application/json", {}).get("schema", {}).get("$ref", None)
            response_schema_json = None
            if response_schema:
                response_schema = response_schema.split("/")[-1]
                response_schema_json = swagger["components"]["schemas"].get(response_schema, None)
            
            apis.append(API(path, api, response_schema_json))
        return apis
    
    def filter_by_tag(self, tag) -> API:
        ret = []
        for api in self.apis:
            if tag in api.tags:
                ret.append(api)
        return ret
        
        
    # This function filters a list of API objects based on two criteria:
    # Whether the methods parameter has any intersection with the methods associated with each api.
    # Whether any of the dictionaries in exclude_tuples match the method and path of each api. If such a match is found, the api is excluded from the final list.
    def filter_by_methods(self, methods, exclude_tuples) -> List[API]:
        if not methods:
            return []

        methods_set = set(methods)

        return [
            api
            for api in self.apis
            if any(method in methods_set for method in api.methods)
            and not any(et['method'] in api.methods and et['path'] == api.path for et in exclude_tuples)
        ]
        

    def filter_by_methods_and_tag(self, methods, tag) -> API:
        save = []
        for api in self.apis:
            if tag in api.tags:
                if len(np.intersect1d(api.methods, methods ))>0:
                    save.append(api)
        return save

    def filter_by_method_and_paths(self, tuples) -> API:
        save = []
        for t in tuples:
            for api in self.apis:
                if t.get('method') in api.methods and api.path == t.get('path'):
                    save.append(api)
        return save