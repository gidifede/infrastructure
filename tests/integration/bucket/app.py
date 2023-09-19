#!/usr/bin/env python3
import os

import aws_cdk as cdk
from aws_cdk import(
    RemovalPolicy,
    BundlingOptions,
    Stack,
    aws_lambda as lambda_,
    aws_s3 as s3,
    aws_kms as kms,
    aws_sns_subscriptions as subscriptions,
    aws_sqs as sqs,
    aws_s3_deployment as s3d
)

from constructs import Construct
from infrastructure.constructs.base.bucket import Bucket
from infrastructure.constructs.contextconfig import ContextConfig
import pathlib


class  Bucket_TestStack(Stack):
    """
    Testing Configuration + Access Construct
    """
    def __init__(self, scope: Construct, construct_id: str, config: ContextConfig ,**kwargs) -> None:
        super().__init__(scope, construct_id,tags=config.get("tags",{}), **kwargs)

        conf_bucket = Bucket(self, "TestBucket", config=config, bucket_name="TestBucket")
        local_path = str(pathlib.Path(__file__).parent.resolve())
        conf_bucket.add_asset(path=f"{local_path}/assets", prefix="test/prefix")


# the app

app = cdk.App()
cfg = ContextConfig.build_config(app)

test_id=app.node.try_get_context("test_id")
if test_id is None:
    raise Exception("no test_id provided")

cfg.configuration['tags']['test_id'] = test_id
s_name = f"TestBucket-{test_id}"

Bucket_TestStack(app, s_name,env=cdk.Environment(account=os.getenv('CDK_DEFAULT_ACCOUNT'), region=os.getenv('CDK_DEFAULT_REGION')), config=cfg)

app.synth()
