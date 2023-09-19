from infrastructure.constructs.base.baseconstruct import BaseConstruct
from infrastructure.constructs.projection_database.custom_initializer import LambdaInitializer
from infrastructure.constructs.base.baseconstruct import Construct
from infrastructure.constructs.contextconfig import ContextConfig
from infrastructure.constructs import constants
import secrets
from aws_cdk import (
    aws_ec2 as ec2,
    aws_rds as rds,
    aws_iam as iam,
    aws_ec2 as ec2,
    Fn, App, RemovalPolicy, Stack, Duration
)

from typing import Sequence


class ViewDatabase(BaseConstruct):

    def __init__(self, scope: Construct, id: str, db_name: str,
                 config: ContextConfig, vpc: ec2.Vpc, subnets: Sequence[ec2.Subnet], **kwargs) -> None:
        super().__init__(scope, id, config, **kwargs)

        security_groups = []
        for security_group in config.get("vpc").get("services").get("security_groups", []):
            security_groups.append(ec2.SecurityGroup.from_security_group_id(
                self, id=f"{security_group}-sg", security_group_id=security_group))

        # Admin credentials
        credentials = rds.Credentials.from_generated_secret(
            "admin", secret_name=self.get_name("admin-cred"))

        # Create an IAM policy for accessing Secrets Manager
        secretsmanager_policy = iam.PolicyStatement(
            effect=iam.Effect.ALLOW,
            actions=[
                "secretsmanager:GetSecretValue",
                "secretsmanager:DescribeSecret",
                "secretsmanager:ListSecrets"
            ],
            resources=["*"]
        )

        # Create an IAM role for accessing Secrets Manager
        secretsmanager_role = iam.Role(self, "SecretsManagerRole",
                                       assumed_by=iam.ServicePrincipal(
                                           "lambda.amazonaws.com")
                                       )

        # Attach the IAM policy to the role
        secretsmanager_role.add_to_policy(secretsmanager_policy)

        # Database Cluster
        database_cluster = rds.DatabaseCluster(self, id="Projection_DB_Cluster",
                                               engine=rds.DatabaseClusterEngine.aurora_mysql(
                                                   version=rds.AuroraMysqlEngineVersion.VER_3_02_0),
                                               instances=1,
                                               instance_props=rds.InstanceProps(
                                                   vpc=vpc,
                                                   vpc_subnets=ec2.SubnetSelection(
                                                       #    subnet_type=ec2.SubnetType.PUBLIC
                                                       subnets=subnets
                                                   ),
                                                   instance_type=ec2.InstanceType(
                                                       "serverless"),
                                                   security_groups=security_groups,
                                                   publicly_accessible=False
                                               ),
                                               storage_encrypted=True,
                                               default_database_name=db_name,
                                               credentials=credentials,
                                               removal_policy=RemovalPolicy.DESTROY
                                               )

        # Edit the generated CloudFormation construct directly: see https://github.com/aws/aws-cdk/issues/20197
        min_capacity = config.get("rds", {}).get("min-capacity", 1)
        max_capacity = config.get("rds", {}).get("min-capacity", 16)
        database_cluster.node.find_child("Resource").add_property_override(
            "ServerlessV2ScalingConfiguration",
            {
                "MinCapacity": min_capacity,
                "MaxCapacity": max_capacity,
            })

        self.database = database_cluster
        self.db_name = db_name
        self.vpc = vpc
