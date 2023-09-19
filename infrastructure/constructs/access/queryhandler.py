import os
import pathlib
import re
from aws_cdk import (
    Tags,
    aws_lambda as _lambda,
    aws_s3 as s3,
    aws_secretsmanager as secretsmanager,
)
from constructs import Construct

from infrastructure.constructs import constants
from infrastructure.constructs.base.function import GoFunction
from infrastructure.constructs.contextconfig import ContextConfig
from infrastructure.constructs.base.baseconstruct import BaseConstruct
from infrastructure.constructs import utils



class QueryHandler(BaseConstruct):
    def __init__(self, scope: Construct, id: str, lambda_name: str, layer: _lambda.LayerVersion,
                 config: ContextConfig, cluster, **kwargs) -> None:
        super().__init__(scope, id, config=config, **kwargs)

        
        vpc_id=config.get("vpc", {}).get("lambda", {}).get("id", None)

        self.query_handler_lambda = GoFunction(
            self,
            f"{lambda_name}_QueryHandler",
            name=f"{lambda_name}_QueryHandler",
            handler="queryhandler",
            asset=f"{constants.LAMBDA_ASSET_ROOT}/queryhandlers/{lambda_name}",
            description="Lambda Query Handler: Triggered by Api Gateway",
             lambda_env={
                "DB_PROJECTION_CREDENTIAL_SECRET_NAME": cluster.secret.secret_name,
                "REGION": self.region},
            config=self.config,
            layers=[layer],
            vpc_id=vpc_id
            )

        cluster.secret.grant_read(self.query_handler_lambda.function)
        
        
        pathVersion=f"{constants.LAMBDA_ASSET_ROOT}/queryhandlers/{lambda_name}/version.go"
        version = utils.get_version(pathVersion)
        if version:
            Tags.of(self).add(key = "Lambda_Version", value =f"{version}" )
