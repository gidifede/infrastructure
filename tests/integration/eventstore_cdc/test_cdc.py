import json
import boto3
import pytest
from datetime import datetime, timedelta
import time
from tests.integration import util
import random


# test dynamodb cdc published on sns topic by the lamda function
@pytest.mark.integration
def test_insert_event_ok(create_stack, eventstore_table, event_topic, sqs_queue, cdc_lambda):

    dynamodb = boto3.client('dynamodb')
    sqs = boto3.client('sqs')
    # function = boto3.client('lambda')

    queue_name=sqs_queue["ResourceARN"].split(":")[-1]
    queue_url = sqs.get_queue_url(
        QueueName=queue_name
    )

    # Add items
    item1 = { 'id': datetime.now().strftime('%Y-%m-%d %H:%M:%S.%f'), 'timestamp': str(random.randint(0, 30000000))}
    item2 = { 'id': datetime.now().strftime('%Y-%m-%d %H:%M:%S.%f'), 'timestamp': str(random.randint(0, 30000000))}
    item2Updated = { 'id': datetime.now().strftime('%Y-%m-%d %H:%M:%S.%f'), 'timestamp': str(random.randint(0, 30000000))}

    response = dynamodb.put_item(
        TableName=eventstore_table,
        Item={
            'aggregate_id': {"S": item1['id']},
            'timestamp': {"N": item1['timestamp']},
        }
    )
    response = dynamodb.put_item(
        TableName=eventstore_table,
        Item={
            'aggregate_id': {"S": item2['id']},
            'timestamp': {"N": item2['timestamp']},
        }
    )
    # Update item
    response = dynamodb.put_item(
        TableName=eventstore_table,
        Item={
            'aggregate_id': {"S": item2Updated['id']},
            'timestamp': {"N": item2Updated['timestamp']},
        }
    )
    # Remove item
    response = dynamodb.delete_item(
        TableName=eventstore_table,
        Key={
            'aggregate_id': {"S": item1['id']},
            'timestamp': {"N": item1['timestamp']},
        }
    )

    # wait until lambda is deployed completely. the cdk deploy create the function but doesn't wait for the comlete deploy
    # waiter = function.get_waiter('function_active_v2')
    # waiter.wait(FunctionName=cdc_lambda)
    util.wait_for_lambda(cdc_lambda)
    
    # From docs: sqs.receive_message for small size message could need to be called multiple times in order to get the message ie: one time could not be enough
    msgList = []
    for _ in range(6):
        messages = sqs.receive_message(
            QueueUrl=queue_url["QueueUrl"],
            MaxNumberOfMessages=1,
            WaitTimeSeconds=5
        )
        for msg in messages.get("Messages", [None]):
            if msg != None:
                msgList.append(msg)
        
        if len(msgList) == 3:
            break

    if len(msgList) == 3:
        for msg in msgList: 
            body = json.loads(msg["Body"])
            m = json.loads(body["Message"])
            find = False
            if m["aggregate_id"] == item1['id'] and m["timestamp"] == item1['timestamp']:
                find = True 
            if m["aggregate_id"] == item2['id'] and m["timestamp"] == item2['timestamp']:
                find = True 
            if m["aggregate_id"] == item2Updated['id'] and m["timestamp"] == item2Updated['timestamp']:
                find = True 
            
            if find:
                sqs.delete_message( 
                    QueueUrl=queue_url["QueueUrl"],
                    ReceiptHandle=msg['ReceiptHandle'])
            else:
                assert False
      
    else:
        assert False, f"expected 3 messages after CDC trigger, received {len(msgList)}."
