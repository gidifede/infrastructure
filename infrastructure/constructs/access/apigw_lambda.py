import json
from aws_cdk import (aws_apigateway as apigw, aws_lambda as _lambda, aws_secretsmanager as secretsmanager,)
from constructs import Construct
from infrastructure.constructs.base.function import GoFunction
from infrastructure.constructs.contextconfig import ContextConfig
from infrastructure.constructs import constants
from infrastructure.constructs.access.apigateway import ApiGatewayConstruct
from infrastructure.constructs.access.route53 import Route53Construct
from infrastructure.constructs.base.baseconstruct import BaseConstruct
from aws_cdk.aws_apigateway import MockIntegration, PassthroughBehavior
from aws_cdk.aws_apigateway import IntegrationResponse, MethodResponse
from typing import List

class ApiGwLambda(BaseConstruct):

    def __init__(self, scope: Construct, id: str, config: ContextConfig) -> None:

        super().__init__(scope, id, config=config)

        stage = config.env

        self.apig = ApiGatewayConstruct(self, 
                                        "ApiGwLambda",
                                        rest_api_name=self.get_name("logisticBackboneApi"),
                                        rest_api_description="Logistic Backbone Apis", 
                                        stage_name=stage, 
                                        config=config)
        
        # add layer to share node_modules
        self.node_modules_layer = _lambda.LayerVersion(
            self,
            self.get_name("LambdaNodeAPIDocsLayer"),
            code=_lambda.Code.from_asset("swagger/function/layers/swagger.zip"),
            compatible_runtimes=[_lambda.Runtime.NODEJS_14_X],
            description="serverless-http, swagger-ui-express",
            layer_version_name="serverless-swagger-ui",
        )
        
    def add_lambda_method(self, verb: str, path: str,
                          handler: GoFunction):
        """
        Create a new resource and method for the provided path
        """
        # Create resource(s)
        resource = self.apig.api.root.resource_for_path(path)

        # Create method
        api_method = resource.add_method(verb,
                                         apigw.LambdaIntegration(handler.function), api_key_required=True)
        
    def add_lambda_methods(self, verb: str, paths: List[str],
                          handler: GoFunction):
        """
        Create a new resource and method for the provided path
        """
        # Create resource(s)
        for path in paths:
            self.add_lambda_method(verb=verb, path=path, handler=handler)
        
    
    def add_validator_endpoints(self, verb: str, paths: List[str]) -> GoFunction:

        validator_lambda = GoFunction(
                    self,
                    "Validator",
                    name="Validator",
                    handler="commandvalidator",
                    asset=f"{constants.LAMBDA_ASSET_ROOT}/commandvalidator/",
                    description="Lambda Validator: Triggered by Api Gateway",
                    config=self.config
                    )
        
        for path in paths:
            self.add_lambda_method(verb=verb, handler=validator_lambda, path=path)
        
        return validator_lambda
    
    def add_query_endpoints(self, verb: str, secret: secretsmanager.Secret,  paths: List[str]) -> GoFunction:
        query_handler_lambda = GoFunction(
            self,
            "QueryHandler",
            name="QueryHandler",
            handler="queryhandler",
            asset=f"{constants.LAMBDA_ASSET_ROOT}/queryhandler/",
            description="Lambda Query Handler: Triggered by Api Gateway",
             lambda_env={
                "DB_PROJECTION_CREDENTIAL_SECRET_NAME": secret.secret_name,
                "REGION": self.region},
            config=self.config)
        
        for path in paths:
            self.add_lambda_method(verb=verb, handler=query_handler_lambda, path=f"/query/{path}")
            
        query_handler_lambda.grant_secret_read(secret)
        
        return query_handler_lambda
        

    def add_mock_method(self, verb: str, path: str, response_content: str):
        api_prefix = "mock"
        
        resource = self.apig.api.root.resource_for_path(f"/{api_prefix}{path}")
        resource.add_method(verb,
                            MockIntegration(
                                integration_responses=[
                                    IntegrationResponse(
                                        status_code="200", 
                                        response_templates={"application/json": json.dumps(response_content)},
                                        response_parameters={
                                            'method.response.header.Access-Control-Allow-Origin': "'*'"
                                        }
                                        ), ],
                                passthrough_behavior=PassthroughBehavior.NEVER,
                                request_templates={
                                    "application/json": "{ \"statusCode\": 200 }"
                                }
                            ),
                            method_responses=[
                                MethodResponse(
                                    status_code="200",
                                    response_parameters={
                                        'method.response.header.Access-Control-Allow-Origin': True
                                        }
                                    )
                                ]
                            )
        # resource.add_cors_preflight(allow_origins=["*"], allow_methods=["GET"])

    def add_swagger_method(self, verb: str, path: str, response_content: str):
        """
        Create a new resource and method for the provided path
        """
        
        api_prefix = "swagger"
        
        resource = self.apig.api.root.resource_for_path(f"{api_prefix}/{path}")
        resource.add_method(verb,
                            MockIntegration(
                                integration_responses=[
                                    IntegrationResponse(status_code="200", response_templates={"application/json": response_content}), ],
                                passthrough_behavior=PassthroughBehavior.NEVER,
                                request_templates={
                                    "application/json": "{ \"statusCode\": 200 }"
                                }
                            ),
                            method_responses=[
                                MethodResponse(status_code="200")]
                            )
       
        path_name = path.replace("/", "_")
        
        # add layer to share swagger file
        swagger_layer = _lambda.LayerVersion(
            self,
            self.get_name(f"{path_name}_APIDocsLayerSwagger"),
            code=_lambda.Code.from_asset(f"{api_prefix}/{path}"),
            compatible_runtimes=[_lambda.Runtime.NODEJS_14_X],
            description=f"{api_prefix}/{path}",
            layer_version_name=f"{api_prefix}_{path_name}",
        )

        api_docs_func = _lambda.Function(
            self,
            self.get_name(f"{path_name}_LambdaSwagger"),
            function_name=self.get_name(f"{path_name}_LambdaSwagger"),
            handler="index.handler",
            runtime=_lambda.Runtime.NODEJS_14_X,
            code=_lambda.Code.from_asset('swagger/function/src'),
            layers=[self.node_modules_layer, swagger_layer],
        )
        
        resource = self.apig.api.root.resource_for_path(f"{api_prefix}/{path}/api-docs")
        resource.add_method("GET", apigw.LambdaIntegration(api_docs_func), api_key_required=False)
        
        # needed to return js and css
        proxy_plus = resource.add_proxy(any_method=False,)
        proxy_plus.add_method(http_method="GET", integration=apigw.LambdaIntegration(api_docs_func))
        
        api_docs_func.add_environment("API_PREFIX", f"/{api_prefix}/{path}")
        
    def add_health_check_method(self):
        resource = self.apig.api.root.resource_for_path("/")
        resource.add_method("GET",
                            MockIntegration(
                                integration_responses=[
                                    IntegrationResponse(status_code="200"), 
                                ],
                                passthrough_behavior=PassthroughBehavior.NEVER,
                                request_templates={
                                    "application/json": "{ \"statusCode\": 200 }"
                                }
                            ),
                            method_responses=[
                                MethodResponse(status_code="200")]
                            )
        
        
        api=f"{self.apig.api.rest_api_id}.execute-api.eu-central-1.amazonaws.com"
        
        Route53Construct(self, id="Route53", config=self.config, 
                         api_custom_domain=api, 
                         hosted_zone_name=self.config.get("api", {}).get("route53-health-check", {}).get("hosted-zone", ""), 
                         domain_name=self.config.get("api", {}).get("route53-health-check", {}).get("domain-name", ""))
        