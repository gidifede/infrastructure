#!/usr/bin/env python3
import os

import aws_cdk as cdk
from aws_cdk import(
    Stack,
    aws_dynamodb as dynamodb,
    aws_sns_subscriptions as subscriptions,
    aws_sqs as sqs
)

from constructs import Construct
from infrastructure.constructs.eventstore_cdc.eventstore import EventStore_CDC
from infrastructure.constructs.contextconfig import ContextConfig


class  EventStore_CDC_TestStack(Stack):
    """
    Testing Configuration + Access Construct
    """
    def __init__(self, scope: Construct, construct_id: str, config: ContextConfig ,**kwargs) -> None:
        super().__init__(scope, construct_id,tags=config.get("tags",{}), **kwargs)
        
        es = dynamodb.Table(
            self,"EventStoreTest",
            table_name="EventStoreTest",
            partition_key=dynamodb.Attribute(name="aggregate_id",type=dynamodb.AttributeType.STRING),
            sort_key=dynamodb.Attribute(name="timestamp",type=dynamodb.AttributeType.NUMBER),
            removal_policy=config.get_dynamodb_table_removal_policy(),
            billing_mode=config.get_dynamodb_billing_mode(),
            stream=dynamodb.StreamViewType.NEW_AND_OLD_IMAGES
            # replication_regions=["us-east-1", "us-east-2", "us-west-2"]
        )

        bd_cdc = EventStore_CDC(self, "EventStore_CDC_TestStack", config=config,
                             table=es)

        # crate a temporary SQS queue to consume SNS message
        bd_cdc.topic.add_subscription(subscriptions.SqsSubscription(sqs.Queue(self, "TestQueue")))



# the app

app = cdk.App()
cfg = ContextConfig.build_config(app)

test_id=app.node.try_get_context("test_id")
if test_id is None:
    raise Exception("no test_id provided")

cfg.configuration['tags']['test_id'] = test_id
s_name = f"TestEventStoreCDC-{test_id}"

EventStore_CDC_TestStack(app, s_name,env=cdk.Environment(account=os.getenv('CDK_DEFAULT_ACCOUNT'), region=os.getenv('CDK_DEFAULT_REGION')), config=cfg)

app.synth()
