from aws_cdk import (
    aws_sqs as sqs, 
    aws_lambda_event_sources as eventsources,
    aws_secretsmanager as secretsmanager,
    aws_sns as sns,
    aws_lambda as lambda_,
    aws_sns_subscriptions as subscriptions,
)
from constructs import Construct

from infrastructure.constructs import constants
from infrastructure.constructs.base.function import GoFunction
from infrastructure.constructs.contextconfig import ContextConfig
from infrastructure.constructs.base.baseconstruct import BaseConstruct
from infrastructure.constructs.projection.function import ProjectionFunction

from typing import Sequence


class Projection(BaseConstruct):

    def __init__(self, scope: Construct, id: str, secret: secretsmanager.Secret, layer: lambda_.LayerVersion,
                 config: ContextConfig, topic: sns.Topic, projection_name: str, filters: Sequence[str], **kwargs) -> None:
        super().__init__(scope, id, config=config, **kwargs)

        # Create an SQS queue and DLQ
        dlq_queue = sqs.Queue(self, "Projection_Queue.DLQ")
        dlq = sqs.DeadLetterQueue(queue=dlq_queue, max_receive_count=1)
        queue = sqs.Queue(self, "Projection_Queue", dead_letter_queue=dlq)

        # Create a Lambda function that consumes messages from the queue
        
        self.projection_lambda = ProjectionFunction(self, "ProjectionFunction", 
                                                    projection_name=projection_name,
                                                    config=config, 
                                                    queue=queue, 
                                                    layer=layer,
                                                    filters=filters,
                                                    secret=secret)
        
        filter_policy = {
            "Type": sns.SubscriptionFilter.string_filter(
                allowlist=filters
            ),
        }

        topic.add_subscription(subscriptions.SqsSubscription(queue=queue, dead_letter_queue=dlq_queue, filter_policy=filter_policy))