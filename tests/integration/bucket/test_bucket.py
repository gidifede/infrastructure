import json
import boto3
import pytest
import botocore


# test dynamodb cdc published on sns topic by the lamda function
@pytest.mark.integration
def test_add_asset_ok(create_stack, bucket):

    s3 = boto3.resource('s3')

    try:
        s3.Object(bucket["TestBucket"], 'test/prefix/test.json').load()
        assert True, "configuration file exists"
    except botocore.exceptions.ClientError as e:
        assert False, "configuration file doesn't exist"

    # delete all objects to delete the bucket at the end
    # bucket = s3.Bucket(bucket["TestBucket"])
    # bucket.objects.all().delete()

    
