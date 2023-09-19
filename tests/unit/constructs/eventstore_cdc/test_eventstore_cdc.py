from aws_cdk import aws_sns as sns
from aws_cdk import aws_s3 as s3
from aws_cdk import aws_dynamodb as dynamodb
import aws_cdk as cdk
from aws_cdk.assertions import Template

from infrastructure.constructs.eventstore_cdc.eventstore import EventStore_CDC

def test_synthesizes_eventstore_cdc(config):
    stack = cdk.Stack()

    table = dynamodb.Table(stack, "Testtable", 
                           partition_key=dynamodb.Attribute(name="aggregate_id",type=dynamodb.AttributeType.STRING), 
                           sort_key=dynamodb.Attribute(name="timestamp",type=dynamodb.AttributeType.NUMBER),
                           stream=dynamodb.StreamViewType.NEW_AND_OLD_IMAGES
                           )
    EventStore_CDC(stack, "LogisticStack", config=config, table=table)

    template = Template.from_stack(stack)



    template.has_resource_properties(
        "AWS::Lambda::Function",
        {
            "Handler": "cdc",  
            "Runtime": "go1.x",
        },
    )

    # Assert that we have created all expected resources
    template.resource_count_is("AWS::SNS::Topic", 1)
    template.resource_count_is("AWS::Lambda::Function", 2)  # created because of log retention = 1 Day
    template.resource_count_is("AWS::DynamoDB::Table", 1)
    template.has_resource_properties("AWS::DynamoDB::Table",
                                    {
                                    "KeySchema":[
                                        {
                                            "AttributeName":"aggregate_id",
                                        },
                                        {
                                            "AttributeName":"timestamp",
                                        },
                                        ],
                                    },
                                )
