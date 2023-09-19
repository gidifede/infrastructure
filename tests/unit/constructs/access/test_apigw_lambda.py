from aws_cdk import aws_sns as sns
import aws_cdk as cdk
from aws_cdk.assertions import Template
from infrastructure.constructs import constants
from infrastructure.constructs.access.apigw_lambda import ApiGwLambda
from infrastructure.constructs.base.function import GoFunction

def test_apigw_lambda(config):
    stack = cdk.Stack()

    apig = ApiGwLambda(stack, "LogisticStack", config=config)

    function = GoFunction(
        stack,
        "test_Validator",
        name="test_Validator",
        asset=f"{constants.LAMBDA_ASSET_ROOT}/commandvalidator",
        handler="main",
        description="Lambda Validator: Triggered by Api Gateway",
        config=config)

    apig.add_lambda_methods(verb="GET", paths=["/test1/v1/test2", "/test", "/test1/v1/test3", "/test3"], handler=function)

    template = Template.from_stack(stack)

    template.has_resource_properties(
        "AWS::Lambda::Function",
        {
            "Handler": "main",  # the handler is managed by GoFunction library
            "Runtime": "go1.x",
        },
    )

    # Assert that we have created all expected resources
    template.resource_count_is("AWS::ApiGateway::RestApi", 1)
    template.resource_count_is("AWS::ApiGateway::Deployment", 1)
    template.resource_count_is("AWS::ApiGateway::Stage", 1)
    template.resource_count_is("AWS::ApiGateway::Resource", 6)
    template.resource_count_is("AWS::Lambda::Function", 2)  # created because of log retention = 1 Day

    template.has_resource_properties("AWS::ApiGateway::RestApi",
                                    {
                                    "Description": "Logistic Backbone Apis",
                                    "Name": "default-logisticbackboneapi",
                                    },
                                )
    template.has_resource_properties("AWS::ApiGateway::Stage",
                                     {
                                    "StageName":"test",    #lo prendo dal file di output perché nel test non viene passato nulla, è corretto??                                
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

def test_apigw_lambda_from_swagger(config):
    stack = cdk.Stack()

    apig = ApiGwLambda(stack, "LogisticStack", config=config)

    function = GoFunction(
        stack,
        "test_Validator",
        name="test_Validator",
        asset=f"{constants.LAMBDA_ASSET_ROOT}/commandvalidator",
        handler="main",
        description="Lambda Validator: Triggered by Api Gateway",
        config=config)

    apig.add_lambda_methods(handler=function, verb="post", paths=["/product/v1/accept"])
    apig.add_swagger_method(verb="get", path="commands/v1", response_content="{}")

    template = Template.from_stack(stack)
    

    template.has_resource_properties(
        "AWS::Lambda::Function",
        {
            "Handler": "main",  # the handler is managed by GoFunction library
            "Runtime": "go1.x",
        },
    )

    # Assert that we have created all expected resources
    template.resource_count_is("AWS::ApiGateway::RestApi", 1)
    template.resource_count_is("AWS::ApiGateway::Deployment", 1)
    template.resource_count_is("AWS::ApiGateway::Stage", 1)
    template.resource_count_is("AWS::ApiGateway::Resource", 8)
    template.resource_count_is("AWS::Lambda::Function", 3)  # created because of log retention = 1 Day

    template.has_resource_properties("AWS::ApiGateway::RestApi",
                                    {
                                    "Description": "Logistic Backbone Apis",
                                    "Name": "default-logisticbackboneapi",
                                    },
                                )
    template.has_resource_properties("AWS::ApiGateway::Stage",
                                     {
                                    "StageName":"test",    #lo prendo dal file di output perché nel test non viene passato nulla, è corretto??                                
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


def test_apigw_with_swagger(config):

    config.configuration["api"]["access-logs-enabled"] = True
    config.configuration["api"]["api-key"]["random"] = False
    config.configuration["api"]["publish-swagger"] = True

    stack = cdk.Stack()

    apig = ApiGwLambda(stack, "LogisticStack", config=config)

    function = GoFunction(
        stack,
        "test_Validator",
        name="test_Validator",
        asset=f"{constants.LAMBDA_ASSET_ROOT}/commandvalidator",
        handler="main",
        description="Lambda Validator: Triggered by Api Gateway",
        config=config)

    apig.add_lambda_methods(verb="GET", paths=["/test1/v1/test2", "/test", "/test1/v1/test3", "/test3"], handler=function)

    publish_swagger_flag = config.get("api", {}).get("publish-swagger", False)
    if publish_swagger_flag:
        apig.add_swagger_method("get", f"commands/v1", "\{\}") 


    template = Template.from_stack(stack)

    # Assert that we have created all expected resources
    template.resource_count_is("AWS::ApiGateway::RestApi", 1)
    template.resource_count_is("AWS::ApiGateway::Deployment", 1)
    template.resource_count_is("AWS::ApiGateway::Stage", 1)
    template.resource_count_is("AWS::ApiGateway::Resource", 11)
    template.resource_count_is("AWS::Lambda::Function", 3)  # created because of log retention = 1 Day
