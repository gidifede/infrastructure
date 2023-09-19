from aws_cdk import (CfnOutput, aws_apigateway as apigw, aws_lambda as _lambda)
from constructs import Construct

from infrastructure.constructs.base.function import GoFunction
from infrastructure.constructs.base.baseconstruct import BaseConstruct
from infrastructure.constructs.contextconfig import ContextConfig


class ApiGatewayQueryConstruct(Construct):

    def __init__(self, scope: Construct, id: str, rest_api_name: str,
                 rest_api_description: str, stage_name: str, config: ContextConfig) -> None:
        super().__init__(scope, id)

        # Create API Gateway
        rest_api = apigw.RestApi(self,
                                 id=id,
                                 rest_api_name=rest_api_name,
                                 description=rest_api_description,
                                 deploy_options=apigw.StageOptions(
                                     stage_name=stage_name
                                 ))
        
        # Add usage plan to api gateway
        usage_plan = rest_api.add_usage_plan(f"{rest_api.rest_api_name}UsagePlan",
                                             name=f"{rest_api.rest_api_name}UsagePlan",
                                             api_stages=[apigw.UsagePlanPerApiStage(
                                                 api=rest_api, stage=rest_api.deployment_stage)]
                                             )
        
        # if config.get("api").get("api-key").get("enabled", False):
        #     api_key_value=config.get("api").get("api-key").get("value")
        #     key = rest_api.add_api_key(
        #         "api-key", api_key_name=config.get("api").get("api-key").get("name"), value=api_key_value)

        #     usage_plan.add_api_key(api_key=key)

        #     CfnOutput(self, "APIKey", value=api_key_value)

        self.api = rest_api

    def add_lambda_method(self, verb: str, path: str,
                          handler: GoFunction):
        """
        Create a new resource and method for the provided path
        """
        # Create resource(s)
        resource = self.api.root.resource_for_path(path)

        # Create method
        api_method = resource.add_method(verb,
                                         apigw.LambdaIntegration(handler.function), api_key_required=True)
