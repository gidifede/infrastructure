from aws_cdk import (
    aws_sqs as sqs, 
    aws_lambda_event_sources as eventsources,
    aws_dynamodb as dynamodb,
    aws_s3 as s3
)
from constructs import Construct

from infrastructure.constructs.contextconfig import ContextConfig
from infrastructure.constructs.base.baseconstruct import BaseConstruct
from infrastructure.constructs.queue.commandhandler import CommandHandler

class SQSQueueCommandHandlerLambdaConstruct(BaseConstruct):

    def __init__(self, scope: Construct, id: str,
                 config: ContextConfig, cfg: s3.Bucket, table: dynamodb.Table, **kwargs) -> None:
        super().__init__(scope, id, config=config, **kwargs)

        # Create an SQS queue and DLQ
        self.dlq_queue = sqs.Queue(self, "Queue.DLQ", content_based_deduplication=True)
        dlq = sqs.DeadLetterQueue(queue=self.dlq_queue, max_receive_count=1)
        self.queue = sqs.Queue(self, "Queue", content_based_deduplication=True, dead_letter_queue=dlq)

        # Create a Lambda function that consumes messages from the queue        
        chl = CommandHandler(self, "CommandHandler", table=table, config=config, queue=self.queue, cfg=cfg)
        # Allow the consumer Lambda to read messages from the queue
        self.queue.grant_consume_messages(chl.command_handler_lambda.function)
        
        # Add an SQS trigger to the lambda
        eventSource = eventsources.SqsEventSource(self.queue, report_batch_item_failures=True)
        chl.command_handler_lambda.function.add_event_source(eventSource)
        
        self.commandhandler_loggroup = chl.command_handler_lambda.function.log_group
