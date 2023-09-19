from unittest import mock
from unittest.mock import patch, mock_open
import json
from swagger.builder import BackboneSwagger
from swagger import builder


SWAGGER = {
    "swagger": "2.0",
    "info": {
        "title": "Simple API overview",
        "version": "v2"
    },
    "paths": {
        "/": {
            "get": {
                "operationId": "listVersionsv2",
                "summary": "List API versions",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "test"
                ],
                "responses": {
                    "200": {
                        "description": "200 300 response",
                        "examples": {
                            "application/json": "{\n    \"versions\": [\n        {\n            \"status\": \"CURRENT\",\n            \"updated\": \"2011-01-21T11:33:21Z\",\n            \"id\": \"v2.0\",\n            \"links\": [\n                {\n                    \"href\": \"http://127.0.0.1:8774/v2/\",\n                    \"rel\": \"self\"\n                }\n            ]\n        },\n        {\n            \"status\": \"EXPERIMENTAL\",\n            \"updated\": \"2013-07-23T11:33:21Z\",\n            \"id\": \"v3.0\",\n            \"links\": [\n                {\n                    \"href\": \"http://127.0.0.1:8774/v3/\",\n                    \"rel\": \"self\"\n                }\n            ]\n        }\n    ]\n}"
                        }
                    },
                    "300": {
                        "description": "200 300 response",
                        "examples": {
                            "application/json": "{\n    \"versions\": [\n        {\n            \"status\": \"CURRENT\",\n            \"updated\": \"2011-01-21T11:33:21Z\",\n            \"id\": \"v2.0\",\n            \"links\": [\n                {\n                    \"href\": \"http://127.0.0.1:8774/v2/\",\n                    \"rel\": \"self\"\n                }\n            ]\n        },\n        {\n            \"status\": \"EXPERIMENTAL\",\n            \"updated\": \"2013-07-23T11:33:21Z\",\n            \"id\": \"v3.0\",\n            \"links\": [\n                {\n                    \"href\": \"http://127.0.0.1:8774/v3/\",\n                    \"rel\": \"self\"\n                }\n            ]\n        }\n    ]\n}"
                        }
                    }
                }
            }
        },
        "/v2": {
            "get": {
                "operationId": "getVersionDetailsv2",
                "summary": "Show API version details",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "test", "test_v2"
                ],
                "responses": {
                    "200": {
                        "description": "200 203 response",
                        "examples": {
                            "application/json": "{\n    \"version\": {\n        \"status\": \"CURRENT\",\n        \"updated\": \"2011-01-21T11:33:21Z\",\n        \"media-types\": [\n            {\n                \"base\": \"application/xml\",\n                \"type\": \"application/vnd.openstack.compute+xml;version=2\"\n            },\n            {\n                \"base\": \"application/json\",\n                \"type\": \"application/vnd.openstack.compute+json;version=2\"\n            }\n        ],\n        \"id\": \"v2.0\",\n        \"links\": [\n            {\n                \"href\": \"http://127.0.0.1:8774/v2/\",\n                \"rel\": \"self\"\n            },\n            {\n                \"href\": \"http://docs.openstack.org/api/openstack-compute/2/os-compute-devguide-2.pdf\",\n                \"type\": \"application/pdf\",\n                \"rel\": \"describedby\"\n            },\n            {\n                \"href\": \"http://docs.openstack.org/api/openstack-compute/2/wadl/os-compute-2.wadl\",\n                \"type\": \"application/vnd.sun.wadl+xml\",\n                \"rel\": \"describedby\"\n            },\n            {\n              \"href\": \"http://docs.openstack.org/api/openstack-compute/2/wadl/os-compute-2.wadl\",\n              \"type\": \"application/vnd.sun.wadl+xml\",\n              \"rel\": \"describedby\"\n            }\n        ]\n    }\n}"
                        }
                    },
                    "203": {
                        "description": "200 203 response",
                        "examples": {
                            "application/json": "{\n    \"version\": {\n        \"status\": \"CURRENT\",\n        \"updated\": \"2011-01-21T11:33:21Z\",\n        \"media-types\": [\n            {\n                \"base\": \"application/xml\",\n                \"type\": \"application/vnd.openstack.compute+xml;version=2\"\n            },\n            {\n                \"base\": \"application/json\",\n                \"type\": \"application/vnd.openstack.compute+json;version=2\"\n            }\n        ],\n        \"id\": \"v2.0\",\n        \"links\": [\n            {\n                \"href\": \"http://23.253.228.211:8774/v2/\",\n                \"rel\": \"self\"\n            },\n            {\n                \"href\": \"http://docs.openstack.org/api/openstack-compute/2/os-compute-devguide-2.pdf\",\n                \"type\": \"application/pdf\",\n                \"rel\": \"describedby\"\n            },\n            {\n                \"href\": \"http://docs.openstack.org/api/openstack-compute/2/wadl/os-compute-2.wadl\",\n                \"type\": \"application/vnd.sun.wadl+xml\",\n                \"rel\": \"describedby\"\n            }\n        ]\n    }\n}"
                        }
                    }
                }
            }
        }
    },
    "consumes": [
        "application/json"
    ]
}
SWAGGER_METHODS = {
    "swagger": "2.0",
    "info": {
        "title": "Simple API overview",
        "version": "v2"
    },
    "paths": {
        "/": {
            "post": {
                "operationId": "listVersionsv2",
                "summary": "List API versions",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "test"
                ],
                "responses": {
                    "200": {
                        "description": "200 300 response",
                        "examples": {
                            "application/json": "{\n    \"versions\": [\n        {\n            \"status\": \"CURRENT\",\n            \"updated\": \"2011-01-21T11:33:21Z\",\n            \"id\": \"v2.0\",\n            \"links\": [\n                {\n                    \"href\": \"http://127.0.0.1:8774/v2/\",\n                    \"rel\": \"self\"\n                }\n            ]\n        },\n        {\n            \"status\": \"EXPERIMENTAL\",\n            \"updated\": \"2013-07-23T11:33:21Z\",\n            \"id\": \"v3.0\",\n            \"links\": [\n                {\n                    \"href\": \"http://127.0.0.1:8774/v3/\",\n                    \"rel\": \"self\"\n                }\n            ]\n        }\n    ]\n}"
                        }
                    },
                    "300": {
                        "description": "200 300 response",
                        "examples": {
                            "application/json": "{\n    \"versions\": [\n        {\n            \"status\": \"CURRENT\",\n            \"updated\": \"2011-01-21T11:33:21Z\",\n            \"id\": \"v2.0\",\n            \"links\": [\n                {\n                    \"href\": \"http://127.0.0.1:8774/v2/\",\n                    \"rel\": \"self\"\n                }\n            ]\n        },\n        {\n            \"status\": \"EXPERIMENTAL\",\n            \"updated\": \"2013-07-23T11:33:21Z\",\n            \"id\": \"v3.0\",\n            \"links\": [\n                {\n                    \"href\": \"http://127.0.0.1:8774/v3/\",\n                    \"rel\": \"self\"\n                }\n            ]\n        }\n    ]\n}"
                        }
                    }
                }
            }
        },
        "/v2": {
            "post": {
                "operationId": "getVersionDetailsv2",
                "summary": "Show API version details",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "test", "test_v2"
                ],
                "responses": {
                    "200": {
                        "description": "200 203 response",
                        "examples": {
                            "application/json": "{\n    \"version\": {\n        \"status\": \"CURRENT\",\n        \"updated\": \"2011-01-21T11:33:21Z\",\n        \"media-types\": [\n            {\n                \"base\": \"application/xml\",\n                \"type\": \"application/vnd.openstack.compute+xml;version=2\"\n            },\n            {\n                \"base\": \"application/json\",\n                \"type\": \"application/vnd.openstack.compute+json;version=2\"\n            }\n        ],\n        \"id\": \"v2.0\",\n        \"links\": [\n            {\n                \"href\": \"http://127.0.0.1:8774/v2/\",\n                \"rel\": \"self\"\n            },\n            {\n                \"href\": \"http://docs.openstack.org/api/openstack-compute/2/os-compute-devguide-2.pdf\",\n                \"type\": \"application/pdf\",\n                \"rel\": \"describedby\"\n            },\n            {\n                \"href\": \"http://docs.openstack.org/api/openstack-compute/2/wadl/os-compute-2.wadl\",\n                \"type\": \"application/vnd.sun.wadl+xml\",\n                \"rel\": \"describedby\"\n            },\n            {\n              \"href\": \"http://docs.openstack.org/api/openstack-compute/2/wadl/os-compute-2.wadl\",\n              \"type\": \"application/vnd.sun.wadl+xml\",\n              \"rel\": \"describedby\"\n            }\n        ]\n    }\n}"
                        }
                    },
                    "203": {
                        "description": "200 203 response",
                        "examples": {
                            "application/json": "{\n    \"version\": {\n        \"status\": \"CURRENT\",\n        \"updated\": \"2011-01-21T11:33:21Z\",\n        \"media-types\": [\n            {\n                \"base\": \"application/xml\",\n                \"type\": \"application/vnd.openstack.compute+xml;version=2\"\n            },\n            {\n                \"base\": \"application/json\",\n                \"type\": \"application/vnd.openstack.compute+json;version=2\"\n            }\n        ],\n        \"id\": \"v2.0\",\n        \"links\": [\n            {\n                \"href\": \"http://23.253.228.211:8774/v2/\",\n                \"rel\": \"self\"\n            },\n            {\n                \"href\": \"http://docs.openstack.org/api/openstack-compute/2/os-compute-devguide-2.pdf\",\n                \"type\": \"application/pdf\",\n                \"rel\": \"describedby\"\n            },\n            {\n                \"href\": \"http://docs.openstack.org/api/openstack-compute/2/wadl/os-compute-2.wadl\",\n                \"type\": \"application/vnd.sun.wadl+xml\",\n                \"rel\": \"describedby\"\n            }\n        ]\n    }\n}"
                        }
                    }
                }
            }
        }
    },
    "consumes": [
        "application/json"
    ]
}

SWAGGER_NO_METHODS = {
    "swagger": "2.0",
    "info": {
        "title": "Simple API overview",
        "version": "v2"
    },
    "paths": {
        "/": {
            "": {
                "operationId": "listVersionsv2",
                "summary": "List API versions",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "test"
                ],
                "responses": {
                    "200": {
                        "description": "200 300 response",
                        "examples": {
                            "application/json": "{\n    \"versions\": [\n        {\n            \"status\": \"CURRENT\",\n            \"updated\": \"2011-01-21T11:33:21Z\",\n            \"id\": \"v2.0\",\n            \"links\": [\n                {\n                    \"href\": \"http://127.0.0.1:8774/v2/\",\n                    \"rel\": \"self\"\n                }\n            ]\n        },\n        {\n            \"status\": \"EXPERIMENTAL\",\n            \"updated\": \"2013-07-23T11:33:21Z\",\n            \"id\": \"v3.0\",\n            \"links\": [\n                {\n                    \"href\": \"http://127.0.0.1:8774/v3/\",\n                    \"rel\": \"self\"\n                }\n            ]\n        }\n    ]\n}"
                        }
                    },
                    "300": {
                        "description": "200 300 response",
                        "examples": {
                            "application/json": "{\n    \"versions\": [\n        {\n            \"status\": \"CURRENT\",\n            \"updated\": \"2011-01-21T11:33:21Z\",\n            \"id\": \"v2.0\",\n            \"links\": [\n                {\n                    \"href\": \"http://127.0.0.1:8774/v2/\",\n                    \"rel\": \"self\"\n                }\n            ]\n        },\n        {\n            \"status\": \"EXPERIMENTAL\",\n            \"updated\": \"2013-07-23T11:33:21Z\",\n            \"id\": \"v3.0\",\n            \"links\": [\n                {\n                    \"href\": \"http://127.0.0.1:8774/v3/\",\n                    \"rel\": \"self\"\n                }\n            ]\n        }\n    ]\n}"
                        }
                    }
                }
            }
        },
        "/v2": {
            "": {
                "operationId": "getVersionDetailsv2",
                "summary": "Show API version details",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "test", "test_v2"
                ],
                "responses": {
                    "200": {
                        "description": "200 203 response",
                        "examples": {
                            "application/json": "{\n    \"version\": {\n        \"status\": \"CURRENT\",\n        \"updated\": \"2011-01-21T11:33:21Z\",\n        \"media-types\": [\n            {\n                \"base\": \"application/xml\",\n                \"type\": \"application/vnd.openstack.compute+xml;version=2\"\n            },\n            {\n                \"base\": \"application/json\",\n                \"type\": \"application/vnd.openstack.compute+json;version=2\"\n            }\n        ],\n        \"id\": \"v2.0\",\n        \"links\": [\n            {\n                \"href\": \"http://127.0.0.1:8774/v2/\",\n                \"rel\": \"self\"\n            },\n            {\n                \"href\": \"http://docs.openstack.org/api/openstack-compute/2/os-compute-devguide-2.pdf\",\n                \"type\": \"application/pdf\",\n                \"rel\": \"describedby\"\n            },\n            {\n                \"href\": \"http://docs.openstack.org/api/openstack-compute/2/wadl/os-compute-2.wadl\",\n                \"type\": \"application/vnd.sun.wadl+xml\",\n                \"rel\": \"describedby\"\n            },\n            {\n              \"href\": \"http://docs.openstack.org/api/openstack-compute/2/wadl/os-compute-2.wadl\",\n              \"type\": \"application/vnd.sun.wadl+xml\",\n              \"rel\": \"describedby\"\n            }\n        ]\n    }\n}"
                        }
                    },
                    "203": {
                        "description": "200 203 response",
                        "examples": {
                            "application/json": "{\n    \"version\": {\n        \"status\": \"CURRENT\",\n        \"updated\": \"2011-01-21T11:33:21Z\",\n        \"media-types\": [\n            {\n                \"base\": \"application/xml\",\n                \"type\": \"application/vnd.openstack.compute+xml;version=2\"\n            },\n            {\n                \"base\": \"application/json\",\n                \"type\": \"application/vnd.openstack.compute+json;version=2\"\n            }\n        ],\n        \"id\": \"v2.0\",\n        \"links\": [\n            {\n                \"href\": \"http://23.253.228.211:8774/v2/\",\n                \"rel\": \"self\"\n            },\n            {\n                \"href\": \"http://docs.openstack.org/api/openstack-compute/2/os-compute-devguide-2.pdf\",\n                \"type\": \"application/pdf\",\n                \"rel\": \"describedby\"\n            },\n            {\n                \"href\": \"http://docs.openstack.org/api/openstack-compute/2/wadl/os-compute-2.wadl\",\n                \"type\": \"application/vnd.sun.wadl+xml\",\n                \"rel\": \"describedby\"\n            }\n        ]\n    }\n}"
                        }
                    }
                }
            }
        }
    },
    "consumes": [
        "application/json"
    ]
}

CONFIG = {
    "config": {
        "title": "Test API",
        "description": "Test API",
        "version": "1.0.0",
        "json_schema_refs": [
            "../configuration/JSON/JSONSCHEMA/common.json"
        ]
    },
    "apis": {
        "product": [
            {
                "name": "Accept",
                "description": "Sends the command to accept the product",
                "path": "/product/v1/accept",
                "method": "post",
                "json_schema_location": "../configuration/JSON/JSONSCHEMA/Product/Commands/External/Accept/Accept.json",
                "responses": [
                    {
                        "code": "202",
                        "content_type": "text/plain",
                        "content_json_schema_location": None,
                        "content_schema_type": "string",
                        "description": "Command accepted"
                    }
                ]
            },
            {
                "name": "TestGet",
                "description": "Test Get description",
                "path": "/test",
                "method": "get",
                "json_schema_location": None,
                "parameters": [
                    {"name": "param", "type": "string", "description": "Product ID"}
                ],
                "responses": [
                    {
                        "code": "202",
                        "content_type": "text/plain",
                        "content_json_schema_location": None,
                        "content_schema_type": "string",
                        "description": "Test Description"
                    }
                ]
            },

        ]
    }
}


def test_load():
    with patch("builtins.open", mock_open(read_data=json.dumps(SWAGGER))) as mock_file:
        bs = BackboneSwagger.load("test", "1")
        assert bs != None

        apis = bs.apis
        assert len(apis) == 2
        assert apis[0].path == "/"
        assert apis[1].path == "/v2"

        assert apis[0].methods == ["get"]
        assert apis[1].methods == ["get"]


def test_build():
    
    with mock.patch('swagger.builder.os.path.isfile') as mocked_isfile, \
        mock.patch('swagger.builder.os.listdir') as mocked_listdir:
        mocked_isfile.return_value = True
        mocked_listdir.return_value = ['test.json']
    
        with patch("builtins.open", mock_open(read_data=json.dumps(CONFIG))) as mock_file:
            bs = BackboneSwagger(name="test", version="v1")
            assert bs != None
            apis = bs.apis
            assert len(apis) == 2
            assert apis[0].path == "/product/v1/accept"
            assert apis[0].methods == ["post"]
            
            assert apis[1].path == "/test"
            assert apis[1].methods == ["get"]


def test_filter():
    with patch("builtins.open", mock_open(read_data=json.dumps(SWAGGER))) as mock_file:
        bs = BackboneSwagger.load("test", "1")
        assert bs != None

        apis = bs.filter_by_tag("test_v2")
        assert len(apis) == 1
        
        apis = bs.filter_by_tag("test")
        assert len(apis) == 2

def test_filter_by_methods():
    with patch("builtins.open", mock_open(read_data=json.dumps(SWAGGER_METHODS))) as mock_file:
        bs = BackboneSwagger.load("test", "1")
        assert bs != None

        apis = bs.filter_by_methods(["post"])
        assert len(apis) == 2
        
        apis = bs.filter_by_methods(["get"])
        assert len(apis) == 0


    with patch("builtins.open", mock_open(read_data=json.dumps(SWAGGER))) as mock_file:
        bss = BackboneSwagger.load("test", "1")
        assert bss != None

        apis = bss.filter_by_methods(["post"])
        assert len(apis) == 0 

        apis = bss.filter_by_methods(["get"])
        assert len(apis) == 2


    with patch("builtins.open", mock_open(read_data=json.dumps(SWAGGER_NO_METHODS))) as mock_file:
        bss = BackboneSwagger.load("test", "1")
        assert bss != None

        apis = bss.filter_by_methods(["post"])
        assert len(apis) == 0 

        apis = bss.filter_by_methods(["get"])
        assert len(apis) == 0
        

@patch('swagger.builder.os')
def test_swagger_save(mock_os):
    with mock.patch('swagger.builder.os.path.isfile') as mocked_isfile, \
        mock.patch('swagger.builder.os.listdir') as mocked_listdir:
        mocked_isfile.return_value = True
        mocked_listdir.return_value = ['test.json']
    
        with patch("builtins.open", mock_open(read_data=json.dumps(CONFIG))) as mock_file:
    # with patch("builtins.open", mock_open(read_data=json.dumps(CONFIG))) as mock_file:
            bs = BackboneSwagger(name="test", version="v1")
    assert bs != None

    with patch("builtins.open", mock_open()) as mock_file:
        bs.save_swagger()
        mock_file.assert_called_once_with(
            f"{builder.localpath}/test/v1/swagger.json", mode='w')

        # format content before comparing with raw string
        mock_file().write.assert_called_once_with(
            json.dumps(json.loads(bs.swagger_content), indent=2))
