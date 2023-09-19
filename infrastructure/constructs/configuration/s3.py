"""
This module implements the configuration with DynamoDB
"""
import os
from aws_cdk import(
    aws_dynamodb as dynamodb,
    aws_s3 as s3,
    CfnTag,
    aws_lambda_event_sources,
    aws_lambda as lambda_,
    aws_sns as sns
)
from constructs import Construct
from infrastructure.constructs.base.bucket import Bucket
from infrastructure.constructs.contextconfig import ContextConfig
from infrastructure.constructs.base.baseconstruct import BaseConstruct


class Configuration(BaseConstruct):
    def __init__(self, scope: "Construct", id: str, config: ContextConfig) -> None:
        super().__init__(scope, id, config)
        
        self.__conf_bucket = Bucket(self, "ConfigurationBucket", config=config, bucket_name="ConfigurationBucket")
        self.__conf_bucket.add_asset(path=config.get_command_handler_state_machine_config_origin(), prefix=config.get_command_handler_state_machine_config_prefix())
        
        folders_to_upload = [
            'Facility',
            'Fleet',
            'Logistic',
            'Network',
            'Parcel',
            'Product',
            'Route',
            'Utils',
            'Common'
            # Add more folders as needed
        ]
        for folder in folders_to_upload:
            # Get the full path of the folder
            file_path = f"configuration/JSON/JSONSCHEMA/{folder}"
            # print(f"path={file_path}, prefix=commandValidator/v1/JSONSCHEMA")
            self.__conf_bucket.add_asset(path=file_path, prefix=f"commandValidator/v1/JSONSCHEMA/{folder}", key=file_path)
        
        self.__conf_bucket.add_asset(path=config.get_command_validator_command_api_map_config_origin(), prefix=config.get_command_validator_command_api_map_config_prefix())


    def get_config(self) -> s3.Bucket:
        return self.__conf_bucket.bucket
