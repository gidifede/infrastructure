import boto3
import os
import logging

import mysql.connector
from mysql.connector import errorcode

import json

logger = logging.getLogger()

logger.setLevel(logging.INFO)


def lambda_handler(event, context):

    secret_name = event.get("secret_name", None)

    secrets_client = boto3.client('secretsmanager', region_name = "eu-central-1")
    
    response = secrets_client.get_secret_value(SecretId=secret_name)
    secrets_dict = json.loads(response['SecretString'])
    
    db_host = secrets_dict['host']
    db_port = secrets_dict['port']
    db_user = secrets_dict['username']
    db_pass = secrets_dict['password']
    db_name = secrets_dict['dbname']

    print(db_host, db_user, db_pass, db_name)

    projectionDb = mysql.connector.connect(
        host=db_host,
        user=db_user,
        password=db_pass,
        database=db_name
    )
    cursor = projectionDb.cursor()

    query = event['query']
    query_type = event['query_type']

    if query_type == "INSERT":
        with projectionDb.cursor() as cur:
            cur.execute(query)
            projectionDb.commit()
            return {
                'statusCode': 200,
            }
    else:
        item_count = 0
        cursor.execute(query)
        for _ in cursor:
            item_count += 1
        return {
            'statusCode': 200,
            'body': json.dumps({"count": item_count}),
            'headers': {
                'Content-Type': 'application/json',
                'Access-Control-Allow-Origin': '*'
            }
        }
