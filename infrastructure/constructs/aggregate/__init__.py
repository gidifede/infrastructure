"""
This module implements the configuration with DynamoDB
"""
import os
from aws_cdk import (
    aws_dynamodb as dynamodb,
    aws_s3 as s3,
    CfnTag,
    aws_lambda_event_sources,
    aws_lambda as lambda_,
    aws_sns as sns,
    aws_secretsmanager as secretsmanager
)
from constructs import Construct
from infrastructure.constructs.contextconfig import ContextConfig
from infrastructure.constructs.base.baseconstruct import BaseConstruct
from infrastructure.constructs.queue.sqs_queue import SQSQueueCommandHandlerLambdaConstruct
from infrastructure.constructs.eventstore_cdc.eventstore import EventStore_CDC
from infrastructure.constructs.projection_database.projection_database import ProjectionDatabase
from infrastructure.constructs.projection.projection import Projection
from infrastructure.constructs.access.commandvalidator import CommandValidator
from infrastructure.constructs.access.queryhandler import QueryHandler
from infrastructure.constructs.log_inspector.construct import LogInspector


class Aggregate(BaseConstruct):
    def __init__(self, 
                 scope: "Construct", 
                 id: str, 
                 config: ContextConfig, 
                 cfg: s3.Bucket, 
                 es:dynamodb.Table, 
                 layer: lambda_.LayerVersion,
                 topic:sns.Topic, 
                 projection_secret:secretsmanager.Secret) -> None:
        super().__init__(scope, id, config)

        queue_lambda_cst = SQSQueueCommandHandlerLambdaConstruct(self,
                                                                 "SQSQueue",
                                                                 config=config,
                                                                 cfg=cfg,
                                                                 table=es
                                                                 )
        
        for projection in config.get("projections", []):
            
            projection_name = projection["name"]
            projection_event_filters = projection["filters"]
            
            Projection(self,
                            f"Projection_{projection_name}",
                            projection_name=projection_name,
                            config=config,
                            layer=layer,
                            secret=projection_secret,
                            topic=topic,
                            filters=projection_event_filters
                            )
        
        self.commandvalidator = CommandValidator(self,
                                                 f"CommandValidator",
                                                 config=config,
                                                 cfg=cfg,
                                                 queue=queue_lambda_cst.queue,
                                                 commandApiMapEnabled=False,
                                                 commandApiMapFilename=f"command-api-map/commandApiMap.json")

        if config.get("enable-log-inspection", False):
            LogInspector(self, "li", log_groups=[
                self.commandvalidator.validator_lambda.function.log_group,
                queue_lambda_cst.commandhandler_loggroup
            ],
                config=config)
            
