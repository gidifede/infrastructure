import zlib
import json
import boto3
import botocore
from base64 import b64decode
import os

table_name = os.getenv("TABLE", None)
dynamodb = boto3.resource('dynamodb')
table = dynamodb.Table(table_name)


def decode(data):
    compressed_payload = b64decode(data)
    json_payload = zlib.decompress(compressed_payload, 16+zlib.MAX_WBITS)
    return json.loads(json_payload)


def lambda_handler(event, context):
    # {
    # "messageType":"DATA_MESSAGE",
    # "owner":"065537138232",
    # "logGroup":"/aws/lambda/MarcoStack-product_CommandHandler",
    # "logStream":"2023/03/30/[$LATEST]d1a90c6c13774e11b73d709edaa18abc",
    # "subscriptionFilters":[
    #     "MarcoStack-Productproductlimarcostackproductsftokenresourceloggroupname9588EE7C4D0-mI0O8ohvxoct"
    # ],
    # "logEvents":[
    #     {
    #         "id":"37469826331875691519851463579940706474514167164586295302",
    #         "timestamp":1680205123116,
    #         "message":"{\"level\":\"debug\",\"time\":1680205123,\"message\":\"REQUEST_JSON_ID|123456789|1111115555555558888888889\"}\n"
    #     }
    # ]
    # }
    log = decode(event["awslogs"]["data"])

    loggroup = log["logGroup"]
    logevents = log["logEvents"]
    
    for logevent in logevents:
        message = json.loads(logevent["message"])
        ts = message["message"].split("|")[1]
        request_id = message["message"].split("|")[-1]

        function = loggroup.split("/")[-1]
        try:
            print("Inserting", request_id, function)
            table.put_item(
                Item={
                    'request_id': str(request_id),
                    'lambda_function': function,
                    'timestamp': str(ts),
                },
                ConditionExpression='attribute_not_exists(request_id) AND attribute_not_exists(lambda_function)'
            )
            print("Inserted", request_id, function)
        except Exception as e:
            print(e)
            continue

    return {
        'statusCode': 200,
        'body': json.dumps('Hello from Lambda!')
    }
