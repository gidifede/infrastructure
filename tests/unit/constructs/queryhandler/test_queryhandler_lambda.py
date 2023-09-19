from aws_cdk import aws_secretsmanager as secretmanager
from aws_cdk import aws_s3 as s3
import aws_cdk as cdk
from aws_cdk.assertions import Template
from infrastructure.constructs.access.queryhandler import QueryHandler

from infrastructure.constructs.access.apigw_lambda import ApiGwLambda

def test_queryhandler_lambda(config):
    stack = cdk.Stack()

    apig = ApiGwLambda(stack, "LogisticStack", config=config)

    queryhandler = QueryHandler(stack,
                        f"QueryHandler_product",
                        config=config,
                        cfg=s3.Bucket(stack, "Test"),
                        secret=secretmanager.Secret(stack, "TestSecrets", secret_name="Test"))

    apig.add_lambda_methods(handler=queryhandler.query_handler_lambda, verb="get", paths=["query/v1/status"])

    template = Template.from_stack(stack)
    # print(template.to_json())

    template.has_resource_properties(
        "AWS::Lambda::Function",
        {
            "Handler": "queryhandler",  # the handler is managed by GoFunction library
            "Runtime": "go1.x",
        },
    )

    # Assert that we have created all expected resources
    template.resource_count_is("AWS::ApiGateway::RestApi", 1)
    template.resource_count_is("AWS::ApiGateway::Deployment", 1)
    template.resource_count_is("AWS::ApiGateway::Stage", 1)
    template.resource_count_is("AWS::ApiGateway::Resource", 3)
    template.resource_count_is("AWS::Lambda::Function", 2)  # created because of log retention = 1 Day

    template.has_resource_properties("AWS::ApiGateway::RestApi",
                                    {
                                    "Description": "Logistic Backbone Apis",
                                    "Name": "default-logisticbackboneapi",
                                    },
                                )
    template.has_resource_properties("AWS::ApiGateway::Stage",
                                     {
                                    "StageName":"test",                                    
                                    "TracingEnabled": True,
                                    "MethodSettings":[{
                                        "DataTraceEnabled": True,
                                        "LoggingLevel": "INFO",
                                        "MetricsEnabled": True,
                                        "HttpMethod":"*",
                                        },
                                        ],
                                    }
                                )
