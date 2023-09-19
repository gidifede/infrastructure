from http import client
import json
import boto3
import pytest
import requests
import os
from datetime import datetime, timedelta
import time
import pathlib
import uuid
from tests.integration import util



@pytest.mark.integration
def test_projection_invalid_messages(tid, create_stack, sqs_projection, sqs_dlq_projection, sns_projection):
    sqs = boto3.resource('sqs')
    local_path = str(pathlib.Path(__file__).parent.resolve())

    invalid_json = os.path.join(local_path, "assests", "invalid.json")
    with open(invalid_json) as f:
        invalid_data = json.loads(f.read())

    valid_json = os.path.join(local_path, "assests", "valid.json")

    with open(valid_json) as f:
        valid_data = json.loads(f.read())

    id = str(uuid.uuid4())
    valid_data["aggregate_id"] = id
    valid_data["data"] = json.dumps(valid_data["data"])
    # Get the queue
    queue_name = sqs_projection["ResourceARN"].split(":")[-1]
    queue = sqs.get_queue_by_name(QueueName=queue_name)

    sns_client = boto3.client('sns')


    response =sns_client.publish_batch(
        TopicArn=sns_projection["ResourceARN"],
        PublishBatchRequestEntries=[
            {
                'Id': str(uuid.uuid4()),
                'Message': json.dumps(valid_data),
            },
            {
                'Id': str(uuid.uuid4()),
                'Message': json.dumps(invalid_data),
            },
           {
                'Id': str(uuid.uuid4()),
                'Message': json.dumps(invalid_data),
            },
        ]
    )
    # give time to process the request
    time.sleep(2)

    # Get the DLQ
    dlq_queue_name = sqs_dlq_projection["ResourceARN"].split(":")[-1]

    # Retrieve the message from dlq
    dlq_queue = sqs.get_queue_by_name(QueueName=dlq_queue_name)
    messages = []

    for _ in range(15):
        print("Retrieving messages from ",  dlq_queue_name)
        messages = messages + dlq_queue.receive_messages()
        print(f"Got  {messages}")
        time.sleep(3)
    assert len(messages) == 2
    for message in messages: #delete the messages inside the dlq after the assert
        message.delete()