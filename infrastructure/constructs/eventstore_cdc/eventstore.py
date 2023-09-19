"""
This module implements the configuration with DynamoDB
"""
import os
from aws_cdk import(
    aws_dynamodb as dynamodb,
    aws_s3 as s3,
    CfnTag,
    aws_lambda_event_sources,
    aws_lambda as lambda_,
    aws_sns as sns
)
from constructs import Construct
from infrastructure.constructs import constants
from infrastructure.constructs.eventstore_cdc.cdc_function import CDC
from infrastructure.constructs.contextconfig import ContextConfig
from infrastructure.constructs.base.baseconstruct import BaseConstruct


class EventStore_CDC(BaseConstruct):
    """
    Configuration Construct.
    The storage is, for now, on a S3 Bucket
    """
    def __init__(self,
        scope: Construct,
        construct_id: str,
        config:ContextConfig,
        table: dynamodb.Table,
        **kwargs) -> None:
        super().__init__(scope, construct_id, config=config, **kwargs)


        topic = sns.Topic(self, "EventTopic")

        self.topic =  topic


        CDC(self, "CDC", config=config, topic=topic, table=table)

