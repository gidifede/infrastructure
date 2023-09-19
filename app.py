#!/usr/bin/env python3
import os

import aws_cdk as cdk

from infrastructure.infrastructure_stack import InfrastructureStack
from infrastructure.constructs.contextconfig import ContextConfig

app = cdk.App()


config = ContextConfig.build_config(app=app)

STACK_NAME = config.get('stack-name', None)
if STACK_NAME is None:
    STACK_NAME = app.node.try_get_context("stack-name")
    if STACK_NAME is None:
        STACK_NAME = "LogisticBackboneStack"

API_KEY = app.node.try_get_context("api-key")
if API_KEY is not None:
    config.set_api_key(API_KEY)

InfrastructureStack(app, STACK_NAME,
                    # If you don't specify 'env', this stack will be environment-agnostic.
                    # Account/Region-dependent features and context lookups will not work,
                    # but a single synthesized template can be deployed anywhere.

                    # Uncomment the next line to specialize this stack for the AWS Account
                    # and Region that are implied by the current CLI configuration.

                    # env=cdk.Environment(account=os.getenv('CDK_DEFAULT_ACCOUNT'), region=os.getenv('CDK_DEFAULT_REGION')),

                    # Uncomment the next line if you know exactly what Account and Region you
                    # want to deploy the stack to. */

                    # env=cdk.Environment(account='123456789012', region='us-east-1'),

                    # For more information, see https://docs.aws.amazon.com/cdk/latest/guide/environments.html
                    env=cdk.Environment(account=os.getenv(
                        'CDK_DEFAULT_ACCOUNT'), region=os.getenv('CDK_DEFAULT_REGION')),
                    config=config
                    )

app.synth()
