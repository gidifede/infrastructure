"""
This module implements the configuration with DynamoDB
"""
import os
from aws_cdk import (
    aws_dynamodb as dynamodb,
    CfnTag,
    aws_lambda_event_sources,
    aws_logs_destinations as destinations,
    aws_logs as logs,
    aws_lambda as _lambda,
    aws_sns as sns
)
import pathlib
import os
from constructs import Construct
from infrastructure.constructs.contextconfig import ContextConfig
from infrastructure.constructs.base.baseconstruct import BaseConstruct
import typing


class LogInspector(BaseConstruct):
    """
    Configuration Construct.
    The storage is, for now, on a S3 Bucket
    """

    def __init__(self,
                 scope: Construct,
                 construct_id: str,
                 config: ContextConfig,
                 log_groups: typing.Optional[typing.Collection[logs.LogGroup]],
                 **kwargs) -> None:
        super().__init__(scope, construct_id, config=config, **kwargs)

        log_table = dynamodb.Table(
            self, "LogInspector",
            table_name=self.get_name("LogInspector"),
            partition_key=dynamodb.Attribute(
                name="request_id", type=dynamodb.AttributeType.STRING),
            sort_key=dynamodb.Attribute(
                name="lambda_function", type=dynamodb.AttributeType.STRING),
            removal_policy=config.get_dynamodb_table_removal_policy(),
            billing_mode=config.get_dynamodb_billing_mode()
        )

        local_path = str(pathlib.Path(__file__).parent.resolve())
        lambda_asset = os.path.join(local_path, "lambdafn")

        function = _lambda.Function(self, self.get_name("log-inspector-fn"),
                                    function_name=self.get_name("log-inspector-fn"),
                                    runtime=_lambda.Runtime.PYTHON_3_9,
                                    handler="handler.lambda_handler",
                                    environment={
                                        "TABLE": log_table.table_name
                                    },
                                    code=_lambda.Code.from_asset(lambda_asset))

        log_table.grant_write_data(function)

        for log_group in log_groups:
            logs.SubscriptionFilter(self,
                                    self.get_name("sf-{log_group.log_group_name}"),
                                    destination=destinations.LambdaDestination(fn=function),
                                    filter_pattern=logs.FilterPattern().any(
                                        logs.FilterPattern.string_value("$.message", "=", "REQUEST_JSON_ID*"),
                                    ),
                                    log_group=log_group
                                    )
