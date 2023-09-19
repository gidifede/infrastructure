import re
from aws_cdk import (
    Tags,
    aws_dynamodb as dynamodb,
    aws_iam as iam,
    aws_s3 as s3,
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


class CommandValidator(BaseConstruct):
    def __init__(self, scope: Construct, id: str,
                 config: ContextConfig, queue: sqs.Queue, cfg: s3.Bucket, 
                 commandApiMapEnabled: bool = False, commandApiMapFilename: str = "", **kwargs) -> None:
        super().__init__(scope, id, config=config, **kwargs)


        vpc_id=None
        if config.get("vpc").get("enabled"):
                vpc_id=config.get("vpc").get("lambda").get("id")

        self.validator_lambda = GoFunction(
            self,
            "Validator",
            name="Validator",
            handler="commandvalidator",
            asset=f"{constants.LAMBDA_ASSET_ROOT}/commandvalidator/",
            description="Lambda Validator: Triggered by Api Gateway",
            lambda_env={
                "SQS_QUEUE_NAME": queue.queue_name,
                "S3_BUCKET_NAME": cfg.bucket_name,
                "S3_CONFIG_PREFIX": config.get_command_validator_json_schema_config_prefix(),
                "MESSAGE_GROUP": "command",
                "CONFIG_VERSION": "v1",
                "COMMAND_ORIGIN": "External",
                "COMMAND_API_PATH_MAP_ENABLED": str(commandApiMapEnabled),
                "COMMAND_API_PATH_MAP_FILENAME": commandApiMapFilename
            },
            config=self.config,
            vpc_id=vpc_id
            )

        self.validator_lambda.grant_read_on_bucket(cfg)

        queue.grant_send_messages(self.validator_lambda.function)

        pathVersion=f"{constants.LAMBDA_ASSET_ROOT}/commandvalidator/version.go"
        version = utils.get_version(pathVersion)
        if version:
            Tags.of(self).add(key = "Lambda_Version", value =f"{version}" )
        
