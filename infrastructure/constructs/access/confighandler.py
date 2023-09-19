import os
import pathlib
from aws_cdk import (
    Tags,
    aws_secretsmanager as secretsmanager,
    aws_lambda as _lambda
)
from constructs import Construct

from infrastructure.constructs import constants
from infrastructure.constructs.base.function import GoFunction
from infrastructure.constructs.contextconfig import ContextConfig
from infrastructure.constructs.base.baseconstruct import BaseConstruct
from infrastructure.constructs import utils



class ConfigHandler(BaseConstruct):
    def __init__(self, scope: Construct, id: str, layer: _lambda.LayerVersion,
                 config: ContextConfig, cluster, **kwargs) -> None:
        super().__init__(scope, id, config=config, **kwargs)

        # local_path = str(pathlib.Path(__file__).parent.resolve())
        # lambda_layer_asset = os.path.join(local_path, "layer", "layer_cert.zip")


        # layer = _lambda.LayerVersion(self, f"{self.stack_name}DocDBLayer",
        #     code=_lambda.Code.from_asset(lambda_layer_asset),
        #     compatible_runtimes=[_lambda.Runtime.PYTHON_3_8, _lambda.Runtime.GO_1_X, _lambda.Runtime.PROVIDED_AL2],
        #     license="Apache-2.0",
        #     description="A layer to expose TLS cert for mongo"
        # )

        vpc_id=config.get("vpc", {}).get("lambda", {}).get("id", None)

        self.config_handler_lambda = GoFunction(
            self,
            "ConfigHandler",
            name="ConfigHandler",
            handler="confighandler",
            asset=f"{constants.LAMBDA_ASSET_ROOT}/confighandler/",
            description="Lambda Config Handler: Triggered by Api Gateway",
             lambda_env={
                "DB_PROJECTION_CREDENTIAL_SECRET_NAME": cluster.secret.secret_name,
                "REGION": self.region},
            config=self.config,
            layers=[layer],
            vpc_id=vpc_id
            )
        
        cluster.secret.grant_read(self.config_handler_lambda.function)

        pathVersion=f"{constants.LAMBDA_ASSET_ROOT}/confighandler/version.go"
        version = utils.get_version(pathVersion)
        
        if version:
            Tags.of(self).add(key = "Lambda_Version", value =f"{version}" )
