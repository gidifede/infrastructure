from aws_cdk import aws_sns as sns
from aws_cdk import aws_secretsmanager as secretmanager
import aws_cdk as cdk
from aws_cdk.assertions import Template

from infrastructure.constructs.projection.projection import Projection


def test_projection(config):
    stack = cdk.Stack()


    Projection(stack, "Projection_product",
               config=config, 
               topic=sns.Topic(stack, "TestTopic"),
               secret=secretmanager.Secret(stack, "TestSecrets", secret_name="Test")
               )

    template = Template.from_stack(stack)

    template.has_resource_properties(
        "AWS::Lambda::Function",
        {
            "Handler": "projection",
            "Runtime": "go1.x",
        },
    )

    # Assert that we have created all expected resources
    template.resource_count_is("AWS::SQS::Queue", 2)
    # created because of log retention = 1 Day
    template.resource_count_is("AWS::Lambda::Function", 2)
