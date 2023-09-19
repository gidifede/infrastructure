from aws_cdk import (
    Stack, Duration,
    aws_lambda as _lambda,
    aws_ec2 as ec2,
)
import pathlib
from constructs import Construct
import aws_cdk as cdk
import os
from infrastructure.constructs.contextconfig import ContextConfig
from infrastructure.constructs.projection_database.rds import ViewDatabase
from infrastructure.constructs.projection_database.custom_initializer import LambdaInitializer
from infrastructure.constructs.projection_database.projection_database import ProjectionDatabase

class ProjectionDatabaseConstruct_TestStack(Stack):

    def __init__(self, scope: Construct, construct_id: str, config: ContextConfig, **kwargs) -> None:
        super().__init__(scope, construct_id, tags=config.get("tags", {}), **kwargs)
        
        db = ProjectionDatabase(self,
                        "Projection_DB",
                        config=config,
                        )

        vpc_id = config.get("vpc", {}).get("services", {}).get("id", None)
        subnets = []
        vpc = ec2.Vpc.from_lookup(self, "VPC", vpc_id=vpc_id)

        for subnet in config.get("vpc", {}).get("services", {}).get("private_subnet_ids", []):
            subnets.append(
                ec2.Subnet.from_subnet_id(self, f"Subnet-{subnet}", subnet)
            )

        local_path = str(pathlib.Path(__file__).parent.resolve())
        lambda_asset = os.path.join(
            local_path, "lambda", "lambda.zip")
        function = _lambda.Function(self, "lambda_function",
                                    # allow_public_subnet=True,
                                    runtime=_lambda.Runtime.PYTHON_3_8,
                                    vpc=vpc,
                                    vpc_subnets=ec2.SubnetSelection(
                                        subnets=subnets,
                                    ),
                                    handler="handler.lambda_handler",
                                    timeout=Duration.seconds(600),
                                    code=_lambda.Code.from_asset(lambda_asset))

        
        function.add_environment(key="secret_name_projection", value=db.db_secret_projection.secret_name)
        function.add_environment(key="secret_name_queryhandler", value=db.db_secret_queryhandler.secret_name)
        db.db_secret_projection.grant_read(function)
        db.db_secret_queryhandler.grant_read(function)

        db.allow_connection_on_aurora(function, "Connection from test lambda")


# the app
app = cdk.App()
cfg = ContextConfig.build_config(app)

test_id = app.node.try_get_context("test_id")
if test_id is None:
    raise Exception("no test_id provided")

cfg.configuration['tags']['test_id'] = test_id
s_name = f"TestProjectionDatabase-{test_id}"

ProjectionDatabaseConstruct_TestStack(app, s_name, env=cdk.Environment(account=os.getenv(
    'CDK_DEFAULT_ACCOUNT'), region=os.getenv('CDK_DEFAULT_REGION')), config=cfg)

app.synth()
