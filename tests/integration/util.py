""" Test Utilities """
import re
import hashlib
import base64
import time
import json
from collections.abc import Sequence
from datetime import datetime, timezone
import requests
import boto3
from botocore.exceptions import ClientError


def s3_upload_file(tid: str,
                   bucket: str,
                   key: str,
                   file_content,
                   storage_class: str="STANDARD",
                   tags: str="",
                   content_type: str="application/json"):
    """ Uploads a file to S3 """
    content_hash = hashlib.new('sha256')
    content_hash.update(file_content)
    shasum = base64.b64encode(content_hash.digest())
    s3_client = boto3.client('s3')
    s3_client.put_object(
        Bucket=bucket,
        Key=key,
        Body=file_content,
        ContentType=content_type,
        ChecksumSHA256=shasum.decode(),
        StorageClass=storage_class,
        Tagging=tags,
        )

def get_topic(tid, t_map: dict):
    """ Gets the topic """
    topics = get_resources_by_tag(tid, ["sns:topic"])

    topic_map = {}
    for topic in topics:
        arn = topic['ResourceARN']
        for t_key in t_map.keys():
            if t_map[t_key] in arn:
                topic_map[t_key] = arn

    return topic_map


def get_dynamodb_tables(tid):
    """ Gets the DynamoDB Tables """
    tables = get_resources_by_tag(tid, ["dynamodb:table"])
    table_product_events = None
    
    for table in tables:
        if "EventStoreTest" in table['ResourceARN']:
            table_product_events = re.search("arn:aws:dynamodb:eu-central-1:[0-9]*:table/(.*)",table['ResourceARN']).group(1)

    return {
            "product_events": table_product_events
            }


def get_buckets(tid, b_map:dict):
    """ Gets the bucket from a map """
    buckets = get_resources_by_tag(tid, ["s3:bucket"])
    assert len(buckets) > 0

    buckets_map = {}
    for bucket in buckets:
        tags = { t['Key']:t for t in bucket['Tags']}
        logical_id = tags['aws:cloudformation:logical-id']['Value']
        for bucket_key in b_map.keys():
            lid = b_map[bucket_key]
            if logical_id.startswith(lid):
                buckets_map[bucket_key] = re.search("arn:aws:s3:::(.*)",bucket['ResourceARN']).group(1)

    return buckets_map


def get_lambdas(tid, l_map:dict):
    """ Gets the bucket from a map """
    lambdas = get_resources_by_tag(tid, ["lambda:function"])
    assert len(lambdas) > 0

    lambdas_map = {}
    for l in lambdas:
        arn = l['ResourceARN']
        for t_key in l_map.keys():
            if l_map[t_key] in arn:
                lambdas_map[t_key] = arn
    return lambdas_map


def get_secret(tid, t_map: dict):
    """ Gets the topic """
    secrets = get_resources_by_tag(tid, ["secretsmanager:secret"])
    secret_map = {}
    for secret in secrets:
        arn = secret['ResourceARN']
        for t_key in t_map.keys():
            if t_map[t_key] in arn:
                secret_map[t_key] = arn

    return secret_map


def get_resources_by_tag(tid: str, filters:Sequence[str]=None):
    """ Gets resources by tag and filter """
    if filters is None:
        filters=[]
    client = boto3.client('resourcegroupstaggingapi')
    response = client.get_resources(
        TagFilters=[
            { 'Key': 'app', 'Values': [ 'LogisticBackbone' ] },
            { 'Key': 'env', 'Values': [ 'test' ] },
            { 'Key': 'test_id', 'Values': [ tid ] },
        ],
        ResourceTypeFilters=filters,
    )

    return response['ResourceTagMappingList']


def get_apigateway_endpoint(resources):
    """ Gets the API Gateway endpoint """
    region = ''
    endpoint = ''
    stage = ''
    for res in resources:
        region = re.search("arn:aws:apigateway:(.*)::(.*)",  res['ResourceARN']).group(1)
        gateway = re.search("arn:aws:apigateway:(.*)::(.*)",  res['ResourceARN']).group(2)
        gateway_parts = gateway.split("/")
        if len(gateway_parts) > 3:
            endpoint = gateway_parts[2]
            stage = gateway_parts[4]

    assert region != ""
    assert endpoint  != ""
    assert stage != ""
    url=f"https://{endpoint}.execute-api.{region}.amazonaws.com/{stage}/"
    return url


def check_object(key: str,
                 bucket: str,
                 storage: str,
                 ret_days: datetime,
                 sleep_seconds: int=5):
    """ Checks object in S3"""
    tries=0
    resp_body = None
    while tries < 5:
        # check original object presence
        s3_client = boto3.client('s3')
        try:
            resp = s3_client.get_object(
                Bucket=bucket,
                Key=key,
            )
            resp_body = resp['Body'].read()
            break
        except ClientError:
            tries = tries +1
            time.sleep(sleep_seconds)
            continue

    assert resp_body is not None

    # head object
    resp = s3_client.head_object(
        Bucket=bucket,
        Key=key,
    )

    assert 'Expiration' in resp
    assert 'Metadata' in resp
    assert 'ObjectLockRetainUntilDate' in resp
    exp_date = re.search("expiry-date=\"(.*)\". rule-id=\"(.*)\"", resp['Expiration']).group(1)
    rule_id = re.search("expiry-date=\"(.*)\". rule-id=\"(.*)\"", resp['Expiration']).group(2)
    assert rule_id == storage

    expected_retention = ret_days.strftime("%Y%m%d")
    lockretention = resp['ObjectLockRetainUntilDate'].strftime("%Y%m%d")

    assert lockretention == expected_retention

def check_item_in_table(observer_table: str, key: str, sleep_seconds: int=5):
    """ Checks item in table """
    tries =0
    result = None
    dclient = boto3.client("dynamodb")
    while tries < 5:
        res = dclient.get_item(
            TableName=observer_table,
            Key={
                'testid':{
                    'S':key
                    }
                }
            )
        if 'Item' in res:
            result = res['Item']
            break
        tries=tries + 1
        time.sleep(sleep_seconds)     # sleep 5 seconds

    return result


def get_signed_uri( api_endpoint: str,
                customer_id: str,
                token: str,
                key: str):
    """ Gets a signed uri for GET """

    headers={
      "x-pagopa-safestorage-cx-id":customer_id,
    }
    if token != "":
        headers["Authorization"] = f"Bearer {token}"

    resp = requests.get(f"{api_endpoint}/safe-storage/v1/files/{key}",
                      headers=headers)

    return resp

def put_signed_uri(tid:str ,
                api_endpoint: str,
                customer_id: str,
                document_type: str,
                token: str,
                content_type: str,
                file_content: bytes):
    """ Gets a signed uri for PUT """
    content_hash = hashlib.new('sha256')
    content_hash.update(file_content)
    shasum = base64.b64encode(content_hash.digest())

    headers={
      "x-pagopa-safestorage-cx-id":customer_id,
      "x-checksum": "SHA-256",
      "x-checksum-value":shasum.decode(),
    }
    if token != "":
        headers["Authorization"] = f"Bearer {token}"

    payload={
            "contentType": content_type,
            "documentType":document_type,
            }
    resp = requests.post(f"{api_endpoint}/safe-storage/v1/files",
                      data=json.dumps(payload),
                      headers=headers)

    return resp,shasum

def upload_file(tid:str ,
                api_endpoint: str,
                customer_id: str,
                document_type: str,
                token: str,
                file_content: bytes,
                content_type: str="application/json"
                ):
    """ Uploads a file via API """


    (resp,shasum) = put_signed_uri(tid=tid,
                       api_endpoint=api_endpoint,
                       customer_id=customer_id,
                       document_type=document_type,
                       token=token,
                       content_type=content_type,
                       file_content=file_content)

    assert resp.status_code == 200

    uri = resp.json()['uploadUrl']
    secret = resp.json()['secret']
    key = resp.json()['key']

    headers={
        "Content-type":content_type,
        "x-amz-checksum-sha256":shasum.decode(),
        "x-amz-meta-secret":secret
        }
    resp = requests.put(uri, headers=headers, data=file_content)
    assert resp.status_code == 200

    return key


def wait_for_lambda(function_name, num_of_iterations=10, wait_time_in_loop=5):
    
    function = boto3.client("lambda")
    
    print(f"Start waiting for {function_name}...")
    
    for _ in range(num_of_iterations):
        time.sleep(wait_time_in_loop)
        now = datetime.utcnow()
        now_tzlocal = datetime.now()
        
        # check if triggers are updated
        es = function.list_event_source_mappings(FunctionName=function_name)
        
        lambda_esm = list(filter(lambda x: x["FunctionArn"] == function_name, es["EventSourceMappings"]))
        
        if len(lambda_esm) > 0:
            event_update_time = lambda_esm[0]["LastModified"]
            now_tzlocal = now_tzlocal.astimezone()
            if now_tzlocal > event_update_time:
                return
            else:
                print(f"EventSourceMapping of {function_name} not ready yet...")
                continue
        else:
            # check if lambda is updated
            response = function.get_function_configuration(FunctionName=function_name)
            function_update_time = datetime.fromisoformat(response["LastModified"].replace("+0000", ""))
            
            if now > function_update_time:
                return
            else:
                print(f"Function {function_name} not ready yet...")


def get_api_key(tid):
    client = boto3.client("apigateway")
    response = client.get_api_keys(includeValues=True)
    
    for api_key in response["items"]:    
        if f"{tid}-test-logistic-api-key".lower() in api_key["name"]:
            return api_key["value"]
