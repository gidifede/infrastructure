import json
import boto3
import pytest
import requests
import time
import pathlib
import uuid
import os
from tests.integration import util

# test api gateway and lambda integration
@pytest.mark.integration
def test_contact_api(tid, create_stack, sqs_queue, sqs_dlq_queue, validator_lambda, api_endpoint):

    sqs = boto3.client('sqs')
    local_path = str(pathlib.Path(__file__).parent.resolve())
    
    # wait until lambda is deployed completely. the cdk deploy create the function but doesn't wait for the comlete deploy
    util.wait_for_lambda(validator_lambda)
       
    endpoint = f"{api_endpoint}product/v1/accept"
    
    queue_name=sqs_queue["ResourceARN"].split(":")[-1]
    queue_url = sqs.get_queue_url(
        QueueName=queue_name
    )
    
    api_key = util.get_api_key(create_stack.tid)
    
    headers = {
        'x-api-key': api_key, #'abcdefghilmnopqrstuvz',
        'Content-Type': 'application/json'
    }
       
    valid_json = os.path.join(local_path,"assests","valid.json")

    validFile = open(valid_json)
    validData = json.load(validFile)
    validData['id'] =  str(uuid.uuid4())  
    response = requests.post(endpoint, headers=headers, json=validData)
    
    print(response)
    print(response.content)
    
    # From docs: sqs.receive_message for small size message could need to be called multiple times in order to get the message ie: one time could not be enough
    sqs_resource = boto3.resource('sqs')
    queue = sqs_resource.Queue(queue_url["QueueUrl"])
    # time.sleep(2)
    for _ in range(10):
        print("Retrieving messages from ",  queue_url["QueueUrl"])
        messages = queue.receive_messages()
        print(f"Got  {messages}")
        if len(messages) > 0:
            break
        time.sleep(2)
    
    message = messages[0] if len(messages) > 0 else None
    if message:
        body = json.loads(message.body)
        assert body == validData
    else:
        assert False, "no messages on queue after validation"    
          

@pytest.mark.integration
def test_contact_api_with_invalid_message(tid, create_stack, sqs_queue, sqs_dlq_queue, validator_lambda, api_endpoint):

    local_path = str(pathlib.Path(__file__).parent.resolve())
    
    # wait until lambda is deployed completely. the cdk deploy create the function but doesn't wait for the comlete deploy
    util.wait_for_lambda(validator_lambda)
    
    endpoint = f"{api_endpoint}product/v1/accept"
    
    # Test invalid
    api_key = util.get_api_key(create_stack.tid)
    headers = {
        'x-api-key': api_key, #'abcdefghilmnopqrstuvz',
        'Content-Type': 'application/json'
    }
    
    invalid_json = os.path.join(local_path,"assests","invalid.json")
    
    invalidFile = open(invalid_json)
    invalidData = json.load(invalidFile)
    invalidData['id'] =  str(uuid.uuid4()) 
    response = requests.post(endpoint, headers=headers,json=invalidData)
    
    print(response)
    print(response.content)
    
    assert response.status_code == 400


@pytest.mark.integration
def test_api_key(create_stack, validator_lambda, api_endpoint):

    # wait until lambda is deployed completely. the cdk deploy create the function but doesn't wait for the comlete deploy
    util.wait_for_lambda(validator_lambda)

    payload = json.dumps({
        "specversion": "1.0",
        "type": "Logistic.PCL.Product.Accept.Accept",
        "source": "Logistic.PCL.UP.OMP",
        "subject": "Command",
        "id": "C234-1234-125457",
        "time": "2023-01-24T14:14:00Z",
        "datacontenttype": "application/json",
        "data": {
            "product": {
                "name": "Poste Delivery Web",
                "id": "CARLOCONTI",
                "type:": "BOX",
                "attributes": []
            },
            "location": {
                "type": "UP",
                "address": "Viale Europa 190",
                "zipcode": "00144",
                "city": "Roma",
                "nation": "Italia",
                "locationCode": "55Y90",
                "attributes": []
            },
            "sender": {
                "name": "Pippo Franco",
                "province": "RM",
                "city": "Roma",
                "address": "Via Nepal 51",
                "cap": "00144",
                "attributes": []
            },
            "receiver": {
                "name": "Paolo Bonolis",
                "province": "RM",
                "city": "Roma",
                "address": "Via della Camilluccia 649",
                "cap": "00135",
                "number": "3333456987",
                "email": "paolo.bonolis@gmail.com",
                "note": "Consegna presso portiere",
                "attributes": []
            },
            "timestamp": "1673541473",
            "attributes": []
        }
    })
    api_key = util.get_api_key(create_stack.tid)
    headers = {
        'x-api-key': api_key, #'abcdefghilmnopqrstuvz',
        'Content-Type': 'application/json'
    }

    endpoint = f"{api_endpoint}product/v1/accept"
    
    print(endpoint)
    response = requests.request(
        "POST", endpoint, headers=headers, data=payload)

    assert response.status_code != 403

    headers["x-api-key"] = "fakekeyfakekeyfakekey"

    response = requests.request(
        "POST", endpoint, headers=headers, data=payload)
    assert response.status_code == 403

@pytest.mark.integration
def test_contact_wrong_api(tid, create_stack, sqs_queue, sqs_dlq_queue, validator_lambda, api_endpoint):
    sqs = boto3.client('sqs')
    local_path = str(pathlib.Path(__file__).parent.resolve())
    
    # wait until lambda is deployed completely. the cdk deploy create the function but doesn't wait for the comlete deploy
    util.wait_for_lambda(validator_lambda)
       
    endpoint = f"{api_endpoint}product/v1/fail_delivery"
    
    queue_name=sqs_queue["ResourceARN"].split(":")[-1]
    queue_url = sqs.get_queue_url(
        QueueName=queue_name
    )
    
    api_key = util.get_api_key(create_stack.tid)
    
    headers = {
        'x-api-key': api_key, #'abcdefghilmnopqrstuvz',
        'Content-Type': 'application/json'
    }
       
    valid_json = os.path.join(local_path,"assests","valid.json")

    validFile = open(valid_json)
    validData = json.load(validFile)
    validData['id'] =  str(uuid.uuid4())  
    response = requests.post(endpoint, headers=headers, json=validData)
    
    print(response)
    print(response.content)
    
    assert response.status_code == 400
    

@pytest.mark.integration
def test_contact_api_with_wrong_command(tid, create_stack, sqs_queue, sqs_dlq_queue, validator_lambda, api_endpoint):
    sqs = boto3.client('sqs')
    local_path = str(pathlib.Path(__file__).parent.resolve())
    
    # wait until lambda is deployed completely. the cdk deploy create the function but doesn't wait for the comlete deploy
    util.wait_for_lambda(validator_lambda)
       
    endpoint = f"{api_endpoint}product/v1/accept"
    
    queue_name=sqs_queue["ResourceARN"].split(":")[-1]
    queue_url = sqs.get_queue_url(
        QueueName=queue_name
    )
    
    api_key = util.get_api_key(create_stack.tid)
    
    headers = {
        'x-api-key': api_key, #'abcdefghilmnopqrstuvz',
        'Content-Type': 'application/json'
    }
       
    valid_json = os.path.join(local_path,"assests","wrong_command.json")

    validFile = open(valid_json)
    validData = json.load(validFile)
    validData['id'] =  str(uuid.uuid4())  
    response = requests.post(endpoint, headers=headers, json=validData)
    
    print(response)
    print(response.content)
    
    assert response.status_code == 400

