import pymysql
import json
import os
import boto3
import json
import uuid
from pymysql.constants import CLIENT



def handler(event, context):
    """
    Lambda function handler
    """
    print("Lambda running")
    payload = event[0]
    
    print(payload['admin_secret'])
    
    secret_id = payload['admin_secret']
    query_secret= payload['query_secret']
    projection_secret = payload['projection_secret']
    init_sql = payload['sql-script']
    
    
    secrets_client = boto3.client('secretsmanager', region_name = "eu-central-1")
    response = secrets_client.get_secret_value(SecretId=secret_id)
    secrets_dict = json.loads(response['SecretString'])
    
    db_host = secrets_dict['host']
    db_port = secrets_dict['port']
    db_user = secrets_dict['username']
    db_pass = secrets_dict['password']
    db_name = secrets_dict['dbname']
    
    secrets_client = boto3.client('secretsmanager', region_name = "eu-central-1")
    response = secrets_client.get_secret_value(SecretId=query_secret)
    secrets_dict = json.loads(response['SecretString'])
    
    query_user = secrets_dict['username']
    query_pass = secrets_dict['password']
    print(query_pass)
    
    secrets_client = boto3.client('secretsmanager', region_name = "eu-central-1")
    response = secrets_client.get_secret_value(SecretId=projection_secret)
    secrets_dict = json.loads(response['SecretString'])
    
    projection_user = secrets_dict['username']
    projection_password = secrets_dict['password']
    print(projection_password)
    
    init_sql_clean = init_sql.replace('\n', '').replace('\t', '')
    print(init_sql_clean)
        
    init_sql_partial_pass = init_sql_clean.replace('projection_handler_password', projection_password)
    print(init_sql_partial_pass)
    init_sql_complete = init_sql_partial_pass.replace('query_handler_password', query_pass)
    print(init_sql_complete)
    
    conn_config = {
        'host' : db_host,
        'port' : db_port,
        'user' : db_user,
        'password' : db_pass,
        'database' : db_name,
        'client_flag': CLIENT.MULTI_STATEMENTS
    }
    

    conn = pymysql.connect(**conn_config)
    cursor = conn.cursor()
    cursor.execute(init_sql_complete)
    conn.commit()
    
    # print("execute query succesfully")

    # query_handler_query = "CREATE USER '%s'@'%' IDENTIFIED BY '%s'"
    # projection_handler_query = "CREATE USER '%s'@'%' IDENTIFIED BY '%s'"

    # cursor.execute(query_handler_query, (query_user, query_pass))
    # cursor.execute(projection_handler_query, (projection_user, projection_password))
    # conn.commit()

    # print("created users succesfully")

    # grant_projection = "GRANT SELECT ON *.* TO '%s'@'%'"
    # grant_projection = "GRANT CREATE, SELECT, INSERT, UPDATE ON *.* TO '%s'@'%'"
    
    # cursor.execute(grant_query, (query_user,))
    # cursor.execute(grant_projection, (projection_user,))
    # conn.commit()
    # print("granted privileges succesfully")
    
    conn.close()