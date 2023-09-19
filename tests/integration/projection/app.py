#!/usr/bin/env python3
import os

import aws_cdk as cdk
import infrastructure

from aws_cdk import (
    SecretsManagerSecretOptions,
    Stack,
    aws_sns as sns,

)
from constructs import Construct
from infrastructure.constructs.configuration.s3 import Configuration
from infrastructure.constructs.contextconfig import ContextConfig
from infrastructure.constructs.projection.projection import Projection
from infrastructure.constructs.projection_database import projection_database

class SQS_Queue_Projection_LambdaConstruct_TestStack(Stack):

    def __init__(self, scope: Construct, construct_id: str, config: ContextConfig, **kwargs) -> None:
        super().__init__(scope, construct_id, tags=config.get("tags", {}), **kwargs)

        # Create the SQS and Lambda construct  (Projection)
        
        # cfg = Configuration(self, "Configuration", config=config)
        db = projection_database.ProjectionDatabase(self,
                                "Projection_DB",
                                config =config,
                                )
        topic = sns.Topic(self, "EventTopic")

        projection = Projection(self,
                    f"Projection",
                    config=config,
                    secret=db.db_secret_projection,
                    topic=topic
                            )

# the app
app = cdk.App()
cfg = ContextConfig.build_config(app)

test_id = app.node.try_get_context("test_id")
if test_id is None:
    raise Exception("no test_id provided")

cfg.configuration['tags']['test_id'] = test_id
s_name = f"TestProjection-{test_id}"

SQS_Queue_Projection_LambdaConstruct_TestStack(app, s_name, env=cdk.Environment(account=os.getenv(
    'CDK_DEFAULT_ACCOUNT'), region=os.getenv('CDK_DEFAULT_REGION')), config=cfg)

app.synth()
