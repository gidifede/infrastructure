#!/usr/bin/env python3
import os

import aws_cdk as cdk
from aws_cdk import (
    RemovalPolicy,
    BundlingOptions,
    Stack,
    aws_lambda as lambda_,
    aws_s3 as s3,
    aws_kms as kms,
    aws_sns_subscriptions as subscriptions,
    aws_sqs as sqs
)
from aws_cdk import (aws_apigateway as apigw, aws_lambda as _lambda)
from constructs import Construct
from infrastructure.constructs import constants
from infrastructure.constructs.access.commandvalidator import CommandValidator
from infrastructure.constructs.contextconfig import ContextConfig
from infrastructure.constructs.base.function import GoFunction
from infrastructure.constructs.access.apigw_lambda import ApiGwLambda
from infrastructure.constructs.base.bucket import Bucket
from infrastructure.constructs.configuration.s3 import Configuration



class API_GW_Lambda_TestStack(Stack):
    """
    Testing Configuration + Access Construct
    """

    def __init__(self, scope: Construct, construct_id: str, config: ContextConfig, **kwargs) -> None:
        super().__init__(scope, construct_id, tags=config.get("tags", {}), **kwargs)
        
        config.set_command_validator_json_schema_config_origin("tests/integration/access-command-api-map/assets")
        cfg = Configuration(self, "Configuration", config=config)
        apigw = ApiGwLambda(self, "API_GW_Lambda_Testack", config=config)
        
        
        dlq_queue = sqs.Queue(self, "TestQueueDLQ", content_based_deduplication=True)
        dlq = sqs.DeadLetterQueue(queue=dlq_queue, max_receive_count=1)
        queue = sqs.Queue(self, "TestQueue", content_based_deduplication=True, dead_letter_queue=dlq)
        
        commandvalidator = CommandValidator(self,
                                            "CommandValidator",
                                            config=config,
                                            cfg=cfg.get_config(),
                                            queue=queue,
                                            commandApiMapEnabled=True,
                                            commandApiMapFilename="command-api-map/fakeMap.json")
        
        apigw.add_lambda_methods(verb="post", paths=["/product/v1/fakeEndpoint1", "/product/v1/fake_Endpoint2"], handler=commandvalidator.validator_lambda)

        

# the app
app = cdk.App()
cfg = ContextConfig.build_config(app)

test_id = app.node.try_get_context("test_id")
if test_id is None:
    raise Exception("no test_id provided")

cfg.configuration['tags']['test_id'] = test_id
s_name = f"TestAPIGWLambda-{test_id}"

API_GW_Lambda_TestStack(app, s_name, env=cdk.Environment(account=os.getenv(
    'CDK_DEFAULT_ACCOUNT'), region=os.getenv('CDK_DEFAULT_REGION')), config=cfg)

app.synth()
