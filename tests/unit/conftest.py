import pytest
import aws_cdk as cdk
import json

from infrastructure.constructs.contextconfig import ContextConfig

# TEST_CONTEXT = {
#     "vpc": {
#         "lambda": {
#           "enable": False,
#           "id":"vpc-022598285c5f4ed46",
#           "availability_zones":[
#           ],
#           "private_subnet_ids": [
#             "subnet-0257c6c2e6901d6f4",
#             "subnet-0e2210ce3136822cb",
#             "subnet-0a9bb5cd12d28bd28"
#           ],
#           "private_subnet_route_table_ids":[
#           ]
#         },
#         "services": {
#           "enable": False,
#           "id":"vpc-022598285c5f4ed46",
#           "availability_zones":[
#           ],
#           "private_subnet_ids": [
#             "subnet-0257c6c2e6901d6f4",
#             "subnet-0e2210ce3136822cb",
#             "subnet-0a9bb5cd12d28bd28"
#           ],
#           "private_subnet_route_table_ids":[
#           ]
#         }
#       },
#     "api-stage": "dev",
#     "api": {
#         "api-key": {
#             "name": "test-logistic-api-key",
#             "enabled": "true",
#             "value": "abcdefghilmnopqrstuvz"
#         },
#         "access-logs-enabled": True
#     },
#     "dynamodb": {
#         "removal-policy": "DESTROY",
#         "billing-mode": "PAY_PER_REQUEST",
#         "stream-specification": "NEW_AND_OLD_IMAGES"
#     },
#     "tags": {
#         "app": "LogisticBackbone",
#         "env": "test"
#     },
#     "lambda": {
#         "build-image": "golang:1.19"
#     },
#     "command-handler": {
#         "idempotency-check-source": "true",
#     },
#     "projection": {
#         "db_projection_credential_secret_name": "product_projection_application_user",
#         "region": "eu-central-1"
#     },
#     "resources-tags": [
#         {
#             "key": "Test - Logistic Backbone Project Tag 1",
#             "value": "A-0123456789"
#         },
#         {
#             "key": "Test - Logistic Backbone Project tag 2",
#             "value": "B-0123456789"
#         },
#     ],
#           "query-handler":{
#         "db_projection_credential_secret_name": "product_projection_application_user",
#         "region": "eu-central-1"
#       },


# }


@pytest.fixture(scope="package")
def app():
    app = cdk.App()
    app.node.set_context("config", "test")
    
    with open("cdk.json") as f:
        content = f.read()
        cdk_conf = json.loads(content)
        test_context = cdk_conf["context"]["test"]
    
        app.node.set_context("test", test_context)
    return app


@pytest.fixture(scope="package")
def config(app):
    ctx = ContextConfig.build_config(app)
    return ctx
