import json
import boto3
import pytest
import random

@pytest.mark.integration
def test_projection_database_query_handler(tid, create_stack, query_handler_secret, test_query_lambda):
    
    query =  "select * from PRODUCT"
    payload = run_query(query, "SELECT", query_handler_secret, test_query_lambda)

    statusCode = payload.get('statusCode')
    
    assert statusCode == 200         
    

@pytest.mark.integration
def test_projection_database_projection(tid, create_stack, projection_secret, test_query_lambda):
    
    query =  "select * from PRODUCT"
    payload = run_query(query, "SELECT", projection_secret, test_query_lambda)

    statusCode = payload.get('statusCode')
    assert statusCode == 200
    
    productId = random.randint(0, 100000)
    
    query =  f"INSERT INTO PRODUCT (id, name, type) VALUES ({productId}, 'iPhone 13', 'Electronics')"
    payload = run_query(query, "INSERT", projection_secret, test_query_lambda)

    statusCode = payload.get('statusCode')
    
    assert statusCode == 200
    
    query = f"select * from PRODUCT where id = {productId}"
    payload = run_query(query, "SELECT", projection_secret, test_query_lambda)

    statusCode = payload.get('statusCode')
    body = json.loads(payload.get('body', "{}"))
    
    assert statusCode == 200
    assert body["count"] == 1


def run_query(query, query_type, secret, function):
    l = boto3.client("lambda")
    query = {
        "query": query,
        "query_type": query_type,
        "secret_name": secret
    }
    
    response = l.invoke(
        FunctionName=function,
        Payload=json.dumps(query)
    )
    
    return json.loads(response['Payload'].read())
