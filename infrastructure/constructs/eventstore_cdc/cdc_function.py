from aws_cdk import(
    Tags,
    aws_dynamodb as dynamodb,
    aws_iam as iam,
    aws_s3 as s3,
    aws_lambda_event_sources,
    aws_lambda as lambda_,
    aws_sns as sns
)
from constructs import Construct

from infrastructure.constructs import constants
from infrastructure.constructs.base.function import GoFunction
from infrastructure.constructs.contextconfig import ContextConfig
from infrastructure.constructs.base.baseconstruct import BaseConstruct
from infrastructure.constructs import utils


class CDC(BaseConstruct):
    def __init__(self, scope: Construct, id: str,
                 config: ContextConfig, topic: sns.Topic, table: dynamodb.Table, **kwargs) -> None:
        super().__init__(scope, id, config=config, **kwargs)


        vpc_id=None
        if config.get("vpc").get("enabled"):
                vpc_id=config.get("vpc").get("lambda").get("id")
        
        lambda_fn = GoFunction(self,"CDC",
                            name="CDC",
                            asset=f"{constants.LAMBDA_ASSET_ROOT}/cdc/",
                            description="Triggered by DynamoDB after inserting an event into it",
                            handler="cdc",
                            lambda_env={
                                "SNS_TOPIC_ARN": topic.topic_arn
                            },
                            config=config,
                            vpc_id=vpc_id
                            )
        lambda_fn.grant_publish_on(topic)

        pathVersion=f"{constants.LAMBDA_ASSET_ROOT}/cdc/version.go"
        version = utils.get_version(pathVersion)
        if version:
            Tags.of(self).add(key = "Lambda_Version", value =f"{version}" )
        
        
        lambda_fn.function.add_event_source(aws_lambda_event_sources.DynamoEventSource(table,
            starting_position=lambda_.StartingPosition.TRIM_HORIZON,
            batch_size=5,
            bisect_batch_on_error=True,
            # on_failure=SqsDlq(dead_letter_queue),
            retry_attempts=10
        ))
        