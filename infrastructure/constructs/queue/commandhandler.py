from aws_cdk import (
    Tags,
    aws_sqs as sqs, 
    aws_s3 as s3,
    aws_lambda_event_sources as eventsources,
    aws_dynamodb as dynamodb
)
from constructs import Construct

from infrastructure.constructs import constants
from infrastructure.constructs.base.function import GoFunction
from infrastructure.constructs.contextconfig import ContextConfig
from infrastructure.constructs.base.baseconstruct import BaseConstruct
from infrastructure.constructs import utils


class CommandHandler(BaseConstruct):
    def __init__(self, scope: Construct, id: str,
                 config: ContextConfig, queue: sqs.Queue, table: dynamodb.Table, cfg: s3.Bucket, **kwargs) -> None:
        super().__init__(scope, id, config=config, **kwargs)
        
        self.command_handler_lambda = GoFunction(
            self,
            "CommandHandler",
            name="CommandHandler",
            asset=f"{constants.LAMBDA_ASSET_ROOT}/commandhandler/",
            handler='command-handler',
            description="Triggered by SQS after inserting an message into it",
            lambda_env={
                "SQS_QUEUE_ARN": queue.queue_arn,
                "SOURCE_FIELD_IDEMPOTENCY_CHECK": config.get_command_handler_idempotency_check_source(),
                "S3_BUCKET_NAME": cfg.bucket_name, 
                "S3_CONFIG_PREFIX": config.get_command_handler_state_machine_config_prefix(),
                "MESSAGE_GROUP": "command",
                "CONFIG_VERSION": "v1",
                "COMMAND_ORIGIN": "External",
                "DYNAMO_TABLE_ARN": table.table_arn,
                "DYNAMO_TABLE_NAME": table.table_name
            },
            config=config)

        self.command_handler_lambda.grant_read_on_bucket(cfg)
        table.grant_read_write_data(self.command_handler_lambda.function)
        pathVersion=f"{constants.LAMBDA_ASSET_ROOT}/commandhandler/version.go"
        version = utils.get_version(pathVersion)
        if version:
            Tags.of(self).add(key = "Lambda_Version", value =f"{version}" )

        
