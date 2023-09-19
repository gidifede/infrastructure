# Generate swagger

From repository root, run

`task swagger`

the json structure used to configure the swagger builder could be like the following:

```
{
    "config": {
        "title": "Logistic Backbone Commands API",
        "description": "Logistic Backbone Commands API",
        "version":"1.0.0",
        "json_schema_refs": [
            "<path to common jsonschema>"
        ]
    },
    "apis": {
        "<tag>": [
            {
                "name": "TestPostEndpoint",
                "description": "description",
                "path": "/TestPostEndpointt",
                "method": "post",
                "json_schema_location": "<path to request jsonschema>",
                "responses": [
                    {
                        "code": "202",
                        "content_type": "text/plain",
                        "content_json_schema_location": null,
                        "content_schema_type": "string",
                        "description": "Command accepted"
                    },
                    {
                        "code": "400",
                        "content_type": "text/plain",
                        "content_json_schema_location": null,
                        "content_schema_type": "string",
                        "description": "Command not validated"
                    },
                    {
                        "code": "500",
                        "content_type": "text/plain",
                        "content_json_schema_location": null,
                        "content_schema_type": "string",
                        "description": "Command error"
                    }
                ]
            },
            ],
            "<tag2>": [
            {
                "name": "TestGetEndpoint",
                "description": "description",
                "path": "/TestGetEndpoint",
                "method": "get",
                "json_schema_location": null,
                "parameters": [
                    {"name": "param1", "type": "string", "description": "Param1"}
                ],
                "responses": [
                    {
                        "code": "200",
                        "content_type": "application/json",
                        "content_json_schema_location": "<path to response jsonschema>",
                        "content_schema_type": null,
                        "description": "The product status"
                    }
                ]
            }
        ]
    }
}

```

save it into a json file (example.json) and add the following lines into generate.py

```
example_swagger = BackboneSwagger(name="example", version="v1", config_file="example.json")
example_swagger.save_swagger()
```

then run `task swagger` to generate the swagger into `<name>/<version>/swagger.json` 

# Generate lambda layer for swaggerUI lambda (temporary solution to expose swagger, maybe cloudfront + s3 is better)

to expose both swagger and swaggerUI app, **ApiGwLambda** construct contains **add_swagger_method** to create 2 endpoints (swagger and swaggerUI (api-docs))

```
bs = BackboneSwagger.load("example", "v1")
...
apigw.add_swagger_method("get", f"{bs.name}/{bs.version}", bs.swagger_content)
```

From function/layers folder, run

`sh install.sh`

the command will install node dependencies (in a zip file) used by lambda function to generate swaggerUI app
