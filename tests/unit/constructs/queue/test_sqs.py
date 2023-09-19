from aws_cdk import aws_dynamodb as dynamodb
from aws_cdk import aws_s3 as s3
import aws_cdk as cdk
from aws_cdk.assertions import Template
from infrastructure.constructs.queue import sqs_queue


def test_sqs_queue(config):
    stack = cdk.Stack()
    table = dynamodb.Table(stack, "Testtable",
                           partition_key=dynamodb.Attribute(
                               name="aggregate_id", type=dynamodb.AttributeType.STRING),
                           sort_key=dynamodb.Attribute(
                               name="timestamp", type=dynamodb.AttributeType.NUMBER),
                           stream=dynamodb.StreamViewType.NEW_AND_OLD_IMAGES
                           )
    cfg = s3.Bucket(stack, "Test")
    sqs_queue.SQSQueueCommandHandlerLambdaConstruct(stack, "test-sqs",
                                                    config=config,
                                                    cfg=cfg,
                                                    table=table)

    template = Template.from_stack(stack)

    # Assert that we have created all expected resources
    template.resource_count_is("AWS::SQS::Queue", 2)
    # created because of log retention = 1 Day
    template.resource_count_is("AWS::Lambda::Function", 2)
    template.resource_count_is("AWS::IAM::Role", 2)
    template.resource_count_is("AWS::IAM::Policy", 2)
    template.resource_count_is("AWS::Lambda::EventSourceMapping", 1)

    # Assert that constructs are being created with the expected properties
    template.has_resource_properties("AWS::SQS::Queue", props={
        "ContentBasedDeduplication": True,
        "FifoQueue": True
    },)

    template.has_resource_properties("AWS::Lambda::Function", props={
        "Handler": "command-handler",  # the handler is managed by GoFunction library
        "Runtime": "go1.x",
    },)

    template.has_resource_properties("AWS::IAM::Policy", props={
        'PolicyDocument': {
            'Statement': [
                {
                    'Action': [
                        "logs:PutRetentionPolicy",
                        "logs:DeleteRetentionPolicy"
                    ],
                    'Effect':'Allow',
                    'Resource':'*'
                }]
        }
    })

    queue_logical_id = None
    for q in template.find_resources("AWS::SQS::Queue").keys():
        if "DLQ" not in q:
            queue_logical_id = q

    table_logical_id = None
    for t in template.find_resources("AWS::DynamoDB::Table").keys():
        table_logical_id = t
        break

    bucket_logical_id = None
    for b in template.find_resources("AWS::S3::Bucket").keys():
        bucket_logical_id = b
        break

    template.has_resource_properties("AWS::IAM::Policy", props={
        'PolicyDocument': {
            'Statement': [
                {
                    'Action': [
                        'xray:PutTraceSegments',
                        'xray:PutTelemetryRecords'
                    ],
                    'Effect':'Allow',
                    'Resource':'*'
                },
                {
                    'Action': [
                        's3:GetObject',
                        's3:GetObjectVersion',
                        's3:GetObjectAttributes',
                        's3:GetObjectTagging',
                        's3:GetObjectVersionAttributes',
                        's3:GetObjectRetention',
                        's3:ListBucket'
                    ],
                    'Effect':'Allow',
                    'Resource': [{
                        'Fn::GetAtt': [bucket_logical_id, 'Arn'],
                    }, {
                        "Fn::Join": [
                            "",
                            [
                                {
                                    "Fn::GetAtt": [
                                        bucket_logical_id,
                                        "Arn"
                                    ]
                                },
                                "/*"
                            ]
                        ]
                    }]
                },
                {
                    "Action": [
                        "dynamodb:BatchGetItem",
                        "dynamodb:GetRecords",
                        "dynamodb:GetShardIterator",
                        "dynamodb:Query",
                        "dynamodb:GetItem",
                        "dynamodb:Scan",
                        "dynamodb:ConditionCheckItem",
                        "dynamodb:BatchWriteItem",
                        "dynamodb:PutItem",
                        "dynamodb:UpdateItem",
                        "dynamodb:DeleteItem",
                        "dynamodb:DescribeTable"
                    ],
                    "Effect":"Allow",
                    "Resource":[
                        {
                            "Fn::GetAtt": [table_logical_id, "Arn"],
                        }, {
                            "Ref": "AWS::NoValue"
                        }
                    ]
                },
                {
                    'Action': [
                        'sqs:ReceiveMessage',
                        'sqs:ChangeMessageVisibility',
                        'sqs:GetQueueUrl',
                        'sqs:DeleteMessage',
                        'sqs:GetQueueAttributes',
                    ],
                    'Effect':'Allow',
                    'Resource':{
                        'Fn::GetAtt': [queue_logical_id, 'Arn']
                    }
                }
            ]
        }
    })
