'''
Infrastucture stack module
'''
import base64
import hashlib
import pathlib
import os

import aws_cdk as core
from aws_cdk import aws_lambda as _lambda
from aws_cdk import aws_dynamodb as dynamodb
from constructs import Construct

from infrastructure.constructs.access.apigw_lambda import ApiGwLambda
from infrastructure.constructs.access.confighandler import ConfigHandler
from infrastructure.constructs.access.queryhandler import QueryHandler
from infrastructure.constructs.aggregate import Aggregate
from infrastructure.constructs.configuration.s3 import Configuration
from infrastructure.constructs.contextconfig import ContextConfig
from infrastructure.constructs.eventstore_cdc.eventstore import EventStore_CDC
from infrastructure.constructs.projection_database.projection_database import \
    ProjectionDatabase
from swagger.builder import BackboneSwagger


class InfrastructureStack(core.Stack):
    '''
    Application stack
    '''
    def __init__(self, scope: Construct, construct_id: str,
                 config: ContextConfig, **kwargs) -> None:
        super().__init__(scope, construct_id, **kwargs)

        bb_swagger = BackboneSwagger.load("bb-swagger", "v1")
        # queries = BackboneSwagger.load("queries", "v1")

        cfg = Configuration(self, "Configuration", config=config)

        db = ProjectionDatabase(self,
                                "Projection_DB",
                                config=config,
                                )

        es = dynamodb.Table(
            self, "common_EventStore",
            table_name=self.get_name("common_EventStore"),
            partition_key=dynamodb.Attribute(
                name="aggregate_id", type=dynamodb.AttributeType.STRING),
            sort_key=dynamodb.Attribute(
                name="timestamp", type=dynamodb.AttributeType.NUMBER),
            removal_policy=config.get_dynamodb_table_removal_policy(),
            billing_mode=config.get_dynamodb_billing_mode(),
            stream=dynamodb.StreamViewType.NEW_AND_OLD_IMAGES
            # replication_regions=["us-east-1", "us-east-2", "us-west-2"]
        )

        esc = EventStore_CDC(self,
                             "EventStore_CDC",
                             config=config,
                             table=es
                             )


        apigw = ApiGwLambda(self, "ApiGwLambda", config=config)

        ### LAYER CONTAINING THE CERTIFICATE TO CONNECT TO DOCUMENTDB
        local_path = str(pathlib.Path(__file__).parent.resolve())
        lambda_layer_asset = os.path.join(local_path, "layer", "layer_cert.zip")
        
        layer = _lambda.LayerVersion(self, f"{self.stack_name}DocDBLayer",
            code=_lambda.Code.from_asset(lambda_layer_asset),
            compatible_runtimes=[_lambda.Runtime.GO_1_X],
            license="Apache-2.0",
            description="A layer to expose TLS cert for mongo"
        )
        ### END LAYER

        confighandler = ConfigHandler(self,
                            "ConfigHandler_common",
                            config=config,
                            layer=layer,
                            cluster=db.cluster)


        aggregate = Aggregate(self,
                            "Aggregate",
                            cfg=cfg.get_config(),
                            config=config,
                            es=es,
                            topic=esc.topic,
                            layer=layer,
                            projection_secret=db.cluster.secret)
        
        apis = bb_swagger.filter_by_methods(["post","put","patch","delete"], exclude_tuples= [
                {"method": "post","path": "/v1/network/setup"},
                {"method": "post","path": "/v1/route/setup"},
                {"method": "post","path": "/v1/product/setup"},
                {"method": "post","path": "/v1/fleet/setup"}
            ])

        apigw.add_lambda_methods(verb="post",
                                paths=[api.path for api in apis],
                                handler=aggregate.commandvalidator.validator_lambda)

        apis = bb_swagger.filter_by_method_and_paths(
            [
                {"method": "post","path": "/v1/network/setup"},
                {"method": "post","path": "/v1/route/setup"},
                {"method": "post","path": "/v1/product/setup"},
                {"method": "post","path": "/v1/fleet/setup"}
            ]
        )
        apigw.add_lambda_methods(verb="post",
                                paths=[api.path for api in apis],
                                handler=confighandler.config_handler_lambda)
        
        
        for qh in config.get("query_handlers", []):
        
            lambda_name = qh["lambda_name"]
            paths = qh["paths"]
            
            qh_lambda = QueryHandler(self,
                                f"QueryHandler_{lambda_name}",
                                lambda_name=lambda_name,
                                config=config,
                                layer=layer,
                                cluster = db.cluster)
    
            # checks if configured apis are decrared in swagger file
            apis = bb_swagger.filter_by_method_and_paths(
                [
                    {"method": "get","path": p} for p in paths
                ]
            )
            apigw.add_lambda_methods(verb="get",
                                    paths=[api.path for api in apis],
                                    handler=qh_lambda.query_handler_lambda)
                    
            
        ### Generating mock methods
        apis = bb_swagger.filter_by_methods(["get"], [])
        
        for api in apis:
            if api.response_mock:
                apigw.add_mock_method(verb="get",
                                        path=api.path,
                                        response_content=api.response_mock)
            


        publish_swagger_flag = config.get("api", {}).get("publish-swagger", False)

        if publish_swagger_flag:
            apigw.add_swagger_method("get",
                                     f"{bb_swagger.name}/{bb_swagger.version}",
                                     bb_swagger.swagger_content)

        if config.get("api", {}).get("route53-health-check", {}).get("enabled", False):
            apigw.add_health_check_method()

    def get_name(self, contruct_name: str):
        """ Class method to return the full contruct name"""

        name = f"{self.stack_name.lower()}-{contruct_name.lower()}"
        if len(name) >= 63:
            m_hash = hashlib.sha256()
            m_hash.update(name.encode())
            suffix = (base64.b64encode(m_hash.digest())).decode()[0:12]
            suffix = suffix.replace("+", "-")
            suffix = suffix.replace("/", "-")
            suffix = suffix.lower()
            name = name[0:50]+suffix

        return name
        