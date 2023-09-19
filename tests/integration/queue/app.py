from aws_cdk import (
    aws_dynamodb as dynamodb,
    Stack,
)
from constructs import Construct
import aws_cdk as cdk
import os
from infrastructure.constructs.configuration.s3 import Configuration
from infrastructure.constructs.contextconfig import ContextConfig
from infrastructure.constructs.base.bucket import Bucket
from infrastructure.constructs.queue.sqs_queue import SQSQueueCommandHandlerLambdaConstruct


class SQS_Queue_CommandHandler_LambdaConstruct_TestStack(Stack):

    def __init__(self, scope: Construct, construct_id: str, config: ContextConfig, **kwargs) -> None:
        super().__init__(scope, construct_id, tags=config.get("tags", {}), **kwargs)

        # Create the SQS and Lambda construct  (Command Handler)
        
        cfg = Configuration(self, "Configuration", config=config)

        es = dynamodb.Table(
            self, "EventStoreTest",
            table_name="EventStoreTest",
            partition_key=dynamodb.Attribute(
                name="aggregate_id", type=dynamodb.AttributeType.STRING),
            sort_key=dynamodb.Attribute(
                name="timestamp", type=dynamodb.AttributeType.NUMBER),
            removal_policy=config.get_dynamodb_table_removal_policy(),
            billing_mode=config.get_dynamodb_billing_mode(),
            stream=dynamodb.StreamViewType.NEW_AND_OLD_IMAGES
            # replication_regions=["us-east-1", "us-east-2", "us-west-2"]
        )

        SQSQueueCommandHandlerLambdaConstruct(self,
                                              "SQSQueue",
                                              config=config,
                                              cfg=cfg.get_config(),
                                              table=es
                                              )


# the app
app = cdk.App()
cfg = ContextConfig.build_config(app)

test_id = app.node.try_get_context("test_id")
if test_id is None:
    raise Exception("no test_id provided")

cfg.configuration['tags']['test_id'] = test_id
s_name = f"TestSQSCommandHandler-{test_id}"

SQS_Queue_CommandHandler_LambdaConstruct_TestStack(app, s_name, env=cdk.Environment(account=os.getenv(
    'CDK_DEFAULT_ACCOUNT'), region=os.getenv('CDK_DEFAULT_REGION')), config=cfg)

app.synth()
