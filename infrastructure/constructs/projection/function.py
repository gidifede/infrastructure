from aws_cdk import (
    Tags,
    aws_dynamodb as dynamodb,
    aws_iam as iam,
    aws_s3 as s3,
    aws_lambda as lambda_,
    aws_secretsmanager as secretsmanager,
    aws_lambda_event_sources as eventsources,
    aws_sqs as sqs
)
from constructs import Construct

from infrastructure.constructs import constants
from infrastructure.constructs.base.function import GoFunction
from infrastructure.constructs.contextconfig import ContextConfig
from infrastructure.constructs.base.baseconstruct import BaseConstruct
from infrastructure.constructs import utils
from typing import Sequence
import json


class ProjectionFunction(BaseConstruct):
    def __init__(self, scope: Construct, id: str, projection_name: str, layer: lambda_.LayerVersion,
                 config: ContextConfig, queue: sqs.Queue, secret: secretsmanager.Secret, filters: Sequence[str], **kwargs) -> None:
        super().__init__(scope, id, config=config, **kwargs)

        # vpc_id=None
        # if config.get("vpc").get("enabled"):
        # force query handler to be in the VPC because of connection with aurora
        vpc_id=config.get("vpc", {}).get("lambda", {}).get("id", None)

        self.projection_lambda = GoFunction(
            self,
            f"{projection_name}",
            name=f"{projection_name}",
            asset=f"{constants.LAMBDA_ASSET_ROOT}/projections/{projection_name}",
            handler='projection',
            description="Triggered by SQS after inserting an message into it",
            lambda_env={
                "CREDENTIAL_SECRET_NAME": secret.secret_name,
                "DB_PROJECTION_CREDENTIAL_SECRET_NAME": secret.secret_name,
                "REGION": config.get_projection_region(),
            },
            config=config,
            layers=[layer],
            vpc_id=vpc_id
            )
        

        secret.grant_read(self.projection_lambda.function)

        queue.grant_consume_messages(self.projection_lambda.function)
        pathVersion=f"{constants.LAMBDA_ASSET_ROOT}/projections/{projection_name}/version.go"
        version = utils.get_version(pathVersion)
        if version:
            Tags.of(self).add(key = "Lambda_Version", value =f"{version}" )


        # Add an SQS trigger to the lambda          
        eventSource = eventsources.SqsEventSource(queue,report_batch_item_failures= True)
        self.projection_lambda.function.add_event_source(eventSource)
