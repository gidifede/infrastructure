import aws_cdk as cdk
from aws_cdk.assertions import Template
from aws_cdk.aws_apigateway import (IntegrationResponse, LambdaIntegration,
                                    MethodResponse, MockIntegration,
                                    PassthroughBehavior)

from infrastructure.constructs import constants
from infrastructure.constructs.access.apigateway import ApiGatewayConstruct
from infrastructure.constructs.base.function import GoFunction


def create_stack(config):
    '''
    function to create a stack
    '''
    stack = cdk.Stack()
    apig = ApiGatewayConstruct(
        stack, "test-apig", "test-api", "test-desc", "test-stage", config=config)

    # add 1 fake resource to avoid the error
    # The REST API doesn't contain any methods
    resource = apig.api.root.resource_for_path("test_mock")
    resource.add_method("ANY",
                        MockIntegration(
                            integration_responses=[
                                IntegrationResponse(status_code="200"), ],
                            passthrough_behavior=PassthroughBehavior.NEVER,
                            request_templates={
                                "application/json": "{ \"statusCode\": 200 }"}
                        ),
                        method_responses=[MethodResponse(status_code="200")]
                        )

    function = GoFunction(
        stack,
        "test_Validator",
        name="test_Validator",
        asset=f"{constants.LAMBDA_ASSET_ROOT}/commandvalidator",
        handler="main",
        description="Lambda Validator: Triggered by Api Gateway",
        config=config)

    resource = apig.api.root.resource_for_path("test")
    resource.add_method("POST", LambdaIntegration(handler=function.function), api_key_required=True)

    return stack


def test_apigw_with_accesslog(config):
    '''
    test to check access log flag enabled
    '''
    config.configuration["api"]["access-logs-enabled"] = True
    config.configuration["api"]["api-key"]["random"] = False

    stack = create_stack(config)

    template = Template.from_stack(stack)
    # print(template.to_json())

    # Assert that we have created all expected resources
    template.resource_count_is("AWS::ApiGateway::RestApi", 1)
    template.resource_count_is("AWS::ApiGateway::Deployment", 1)
    template.resource_count_is("AWS::ApiGateway::Stage", 1)
    template.resource_count_is("AWS::ApiGateway::Resource", 2)
    # created because of log retention = 1 Day
    template.resource_count_is("AWS::Lambda::Function", 2)
    template.resource_count_is("AWS::Logs::LogGroup", 1)
    # API KEY must be already configured, the stack just link it
    template.resource_count_is("AWS::ApiGateway::ApiKey", 0)
    template.resource_count_is("AWS::ApiGateway::UsagePlan", 1)

    template.has_resource_properties("AWS::Logs::LogGroup",
                                        {
                                        "LogGroupName": "default-apigatewayaccesslog",
                                        "RetentionInDays": 30,
                                        },
                                    )
    template.has_resource_properties("AWS::ApiGateway::Stage",
                                     {
                                    "StageName":"test-stage",                                    
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

    template.has_resource_properties("AWS::ApiGateway::RestApi",
                                    {
                                    "Description": "test-desc",
                                    "Name": "test-api",
                                    },
                                )
    template.has_resource_properties("AWS::ApiGateway::Method",
                                    {
                                    "HttpMethod": "POST",
                                    "ApiKeyRequired": True,
                                    },
                                )
    template.has_resource_properties("AWS::ApiGateway::UsagePlan",
                                        {
                                        "UsagePlanName": "test-apiUsagePlan",
                                        },
                                    )


def test_apigw_without_accesslog(config):
    '''
    test to check access log flag disabled
    '''
    config.configuration["api"]["access-logs-enabled"] = False
    stack = create_stack(config)

    template = Template.from_stack(stack)

    # Assert that we have created all expected resources
    template.resource_count_is("AWS::ApiGateway::RestApi", 1)
    template.resource_count_is("AWS::ApiGateway::Deployment", 1)
    template.resource_count_is("AWS::ApiGateway::Stage", 1)
    template.resource_count_is("AWS::ApiGateway::Resource", 2)
    # created because of log retention = 1 Day
    template.resource_count_is("AWS::Lambda::Function", 2)
    template.resource_count_is("AWS::Logs::LogGroup", 0)
    # API KEY must be already configured, the stack just link it
    template.resource_count_is("AWS::ApiGateway::ApiKey", 0)
    template.resource_count_is("AWS::ApiGateway::UsagePlan", 1)
