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


# test api gateway and lambda integration
@pytest.mark.integration
def test_command_handler(tid, create_stack, eventstore_table, sqs_command_handler, command_handler_lambda):
    sqs = boto3.client('sqs')
    dynamodb = boto3.client('dynamodb')
    local_path = str(pathlib.Path(__file__).parent.resolve())

    valid_json = os.path.join(local_path, "assests", "valid.json")

    with open(valid_json) as f:
        valid_data = json.loads(f.read())

    id = str(uuid.uuid4())
    valid_data["id"] = id
    valid_data["data"]["product"]["id"] = id

    # Get the queue
    queue_name = sqs_command_handler["ResourceARN"].split(":")[-1]
    queue_url = sqs.get_queue_url(
        QueueName=queue_name
    )

    # wait until lambda is deployed completely. the cdk deploy create the function but doesn't wait for the complete deploy
    util.wait_for_lambda(command_handler_lambda)

    # Create a new message
    response = sqs.send_message(
        QueueUrl=queue_url["QueueUrl"],
        MessageBody=json.dumps(valid_data),
        MessageGroupId='product-command'
    )

    print(response.get('MessageId'))
    print(response.get('MD5OfMessageBody'))

    # give time to process the request
    time.sleep(2)

    response = dynamodb.query(TableName=eventstore_table,
                              KeyConditions={'aggregate_id': {'AttributeValueList': [{'S': id}],
                                                              'ComparisonOperator': 'EQ'}})

    # check 1 row is inserted and that the generated timestamp is different from the timestamp from the CMD
    count = response["Count"]
    assert count > 0


@pytest.mark.integration
def test_serverside_generated_timestamp(tid, create_stack, eventstore_table, sqs_command_handler, command_handler_lambda):
    sqs = boto3.client('sqs')
    dynamodb = boto3.client('dynamodb')
    local_path = str(pathlib.Path(__file__).parent.resolve())

    valid_json = os.path.join(local_path, "assests", "valid.json")

    with open(valid_json) as f:
        valid_data = json.loads(f.read())

    id = str(uuid.uuid4())
    valid_data["id"] = id
    valid_data["data"]["product"]["id"] = id

    # Get the queue
    queue_name = sqs_command_handler["ResourceARN"].split(":")[-1]
    queue_url = sqs.get_queue_url(
        QueueName=queue_name
    )

    # wait until lambda is deployed completely. the cdk deploy create the function but doesn't wait for the complete deploy
    util.wait_for_lambda(command_handler_lambda)

    # Create a new message
    response = sqs.send_message(
        QueueUrl=queue_url["QueueUrl"],
        MessageBody=json.dumps(valid_data),
        MessageGroupId='product-command'
    )

    # give time to process the request
    time.sleep(2)

    response = dynamodb.query(
        TableName=eventstore_table,
        KeyConditions={
            'aggregate_id': {
                'AttributeValueList': [{'S': id}],
                'ComparisonOperator': 'EQ'
            }
        }
    )

    # check 1 row is inserted and that the generated timestamp is different from the timestamp from the CMD
    count = response["Count"]
    if count > 0:
        item = response["Items"][0]
        timestamp = item["timestamp"]
        timestampsent = item["timestampSent"]
        assert timestamp != timestampsent
    else:
        assert False, "Item not returned"


@pytest.mark.integration
def test_command_handler_invalid_messages(tid, create_stack, eventstore_table, sqs_command_handler, command_handler_lambda, sqs_dlq_command_handler):
    sqs = boto3.resource('sqs')
    local_path = str(pathlib.Path(__file__).parent.resolve())

    invalid_json = os.path.join(local_path, "assests", "invalid.json")
    with open(invalid_json) as f:
        invalid_data = json.loads(f.read())

    valid_json = os.path.join(local_path, "assests", "valid.json")

    with open(valid_json) as f:
        valid_data = json.loads(f.read())

    id = str(uuid.uuid4())
    valid_data["id"] = id
    valid_data["data"]["product"]["id"] = id
    # dataWithRandomUUID = data.replace("changeId", id)

    # Get the queue
    queue_name = sqs_command_handler["ResourceARN"].split(":")[-1]
    
    # wait until lambda is deployed completely. the cdk deploy create the function but doesn't wait for the complete deploy
    util.wait_for_lambda(command_handler_lambda)

    queue = sqs.get_queue_by_name(QueueName=queue_name)
    # Send a new  messages
    response = queue.send_messages(
        Entries=[
            {
                "Id": str(uuid.uuid4()),
                "MessageGroupId": 'product-command',
                "MessageBody": json.dumps(valid_data)
            },
            {
                "Id": str(uuid.uuid4()),
                "MessageGroupId": 'product-command',
                "MessageBody": json.dumps(invalid_data)
            }

        ]
    )

    # give time to process the request
    time.sleep(2)

    # Get the DLQ
    dlq_queue_name = sqs_dlq_command_handler["ResourceARN"].split(":")[-1]

    # Retrieve the message from dlq
    dlq_queue = sqs.get_queue_by_name(QueueName=dlq_queue_name)
    for _ in range(20):
        print("Retrieving messages from ",  dlq_queue_name)
        messages = dlq_queue.receive_messages()
        print(f"Got  {messages}")
        if len(messages) != 0:
            assert len(messages) == 1
            break
        time.sleep(2)
    assert len(messages) != 0


@pytest.mark.integration
def test_command_handler_idempotency(tid, create_stack, eventstore_table, sqs_command_handler, command_handler_lambda):
    sqs = boto3.resource('sqs')
    dynamodb = boto3.client('dynamodb')
    local_path = str(pathlib.Path(__file__).parent.resolve())

    valid_json = os.path.join(local_path, "assests", "valid.json")

    with open(valid_json) as f:
        msgOK = json.loads(f.read())

    id = str(uuid.uuid4())
    msgOK["id"] = id
    msgOK["data"]["product"]["id"] = id
    msgOKExpType = 'Logistic.PCL.Product.Accept.Accepted'

    msgDuplicated = msgOK.copy()
    msgDuplicated['time'] = '2023-01-01T101:23:45.678Z'
    msgDifferentType = msgOK.copy()
    msgDifferentType["type"] = "Logistic.PCL.Product.Processing.StartProcessing"
    msgDifferentTypeExpType = 'Logistic.PCL.Product.Processing.ProcessingStarted'

    # Get the queue
    queue_name = sqs_command_handler["ResourceARN"].split(":")[-1]
    queue = sqs.get_queue_by_name(QueueName=queue_name)

    # wait until lambda is deployed completely. the cdk deploy create the function but doesn't wait for the complete deploy
    util.wait_for_lambda(command_handler_lambda)

    # Send msgs
    response = queue.send_messages(
        Entries=[
            {
                "Id": str(uuid.uuid4()),
                "MessageGroupId": 'product-command',
                "MessageBody": json.dumps(msgOK)
            },
            {
                "Id": str(uuid.uuid4()),
                "MessageGroupId": 'product-command',
                "MessageBody": json.dumps(msgDuplicated)
            },
            {
                "Id": str(uuid.uuid4()),
                "MessageGroupId": 'product-command',
                "MessageBody": json.dumps(msgDifferentType)
            }
        ]
    )

    # give time to process the request
    time.sleep(2)
    
    response = dynamodb.query(TableName=eventstore_table,
                              KeyConditions={'aggregate_id': {'AttributeValueList': [{'S': id}],
                                                              'ComparisonOperator': 'EQ'}})

    # check inserted values
    print(response)
    count = response["Count"]
    assert count == 2
    assert response['Items'][0]['aggregate_id']['S'] == msgOK['id']
    assert response['Items'][0]['type']['S'] == msgOKExpType
    assert response['Items'][1]['aggregate_id']['S'] == msgDifferentType['id']
    assert response['Items'][1]['type']['S'] == msgDifferentTypeExpType