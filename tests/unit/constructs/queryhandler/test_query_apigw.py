from aws_cdk import aws_sns as sns
import aws_cdk as cdk
from aws_cdk.assertions import Template
from infrastructure.constructs import constants
from infrastructure.constructs.access.apigw_lambda import ApiGwLambda
from infrastructure.constructs.base.function import GoFunction


def test_query_apigw(config):
    stack = cdk.Stack()

    apig = ApiGwLambda(stack, "LogisticStack", config=config)

    function = GoFunction(
        stack,
        "test_QueryHandler",
        name="test_QueryHandler",
        asset=f"{constants.LAMBDA_ASSET_ROOT}/queryhandler",
        handler="main",
        description="Lambda QueryHandler: Triggered by Api Gateway",
        config=config)

    apig.add_lambda_method("GET", "/test1/v1/test2", function)
    apig.add_lambda_method("GET", "/test", function)
    apig.add_lambda_method("GET", "/test1/v1/test3", function)
    apig.add_lambda_method("GET", "/test3", function)

    template = Template.from_stack(stack)
    # print(template.to_json())

    # Assert that we have created all expected resources
    template.resource_count_is("AWS::ApiGateway::RestApi", 1)
    template.resource_count_is("AWS::ApiGateway::Deployment", 1)
    template.resource_count_is("AWS::ApiGateway::Stage", 1)
    template.resource_count_is("AWS::ApiGateway::Resource", 6)
    # created because of log retention = 1 Day
    template.resource_count_is("AWS::Lambda::Function", 2)

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

