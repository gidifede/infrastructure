from aws_cdk import (
    aws_rds as rds,
    aws_ec2 as ec2,
    aws_rds as rds,
    aws_lambda as _lambda,
    aws_iam as iam,
    aws_secretsmanager as secretsmanager,
    aws_cloudformation as cfn,
    custom_resources as cr,
    Duration, SecretValue
)


from aws_cdk.custom_resources import (
    AwsCustomResource,
    AwsCustomResourcePolicy,
    PhysicalResourceId,
    AwsSdkCall
)

import pathlib
import os
import json
import time
from infrastructure.constructs.base.baseconstruct import BaseConstruct
from infrastructure.constructs.base.baseconstruct import Construct
from infrastructure.constructs.contextconfig import ContextConfig
from typing import Sequence
# Define a custom resource to execute a MySQL script against the RDS instance


class LambdaInitializer(BaseConstruct):

    def __init__(self, scope: Construct, id: str, vpc: ec2.Vpc, subnets: Sequence[ec2.Subnet], 
                 database: rds.DatabaseCluster,
                 cluster_endpoint: rds.Endpoint, cluster_db_name: str,
                 secret: rds.DatabaseSecret, config: ContextConfig, **kwargs) -> None:
        super().__init__(scope, id, config, **kwargs)

        local_path = str(pathlib.Path(__file__).parent.resolve())
        lambda_asset = os.path.join(
            local_path, "init-lambda", "init-lambda.zip")
        sql_script_asset = os.path.join(
            local_path, "scripts", "product-projection-init.sql")

        security_groups = []
        if vpc:
            for security_group in config.get("vpc").get("lambda").get("security_groups",[]):
                security_groups.append(ec2.SecurityGroup.from_security_group_id(self,id=f"{security_group}-sg",security_group_id=security_group))

        function = _lambda.Function(self, "lambda_function",
                                    # allow_public_subnet=True,
                                    runtime=_lambda.Runtime.PYTHON_3_8,
                                    vpc=vpc,
                                    vpc_subnets=ec2.SubnetSelection(
                                        subnets=subnets,
                                    ),
                                    security_groups=security_groups,
                                    handler="init-lambda.handler",
                                    timeout=Duration.seconds(600),
                                    code=_lambda.Code.from_asset(lambda_asset))

        invoke_lambda_statement = iam.PolicyStatement(
            actions=["lambda:InvokeFunction"],
            resources=[function.function_arn],
            effect=iam.Effect.ALLOW
        )

        lambda_policy = iam.Policy(
            self, "LambdaPolicy",
            statements=[invoke_lambda_statement]
        )

        lambda_role = iam.Role(
            self, "LambdaRole",
            assumed_by=iam.ServicePrincipal("lambda.amazonaws.com"),
        )

        lambda_policy.attach_to_role(lambda_role)

        with open(sql_script_asset, "r") as f:
            init_sql = f.read()

        f.close()
        init_sql_clean = init_sql.replace('\n', '').replace('\t', '')

        projection_secret = self.create_database_secret(
            username="projection_handler", db_name=cluster_db_name, host=cluster_endpoint.hostname, port=cluster_endpoint.port)
        query_secret = self.create_database_secret(
            username="query_handler", db_name=cluster_db_name, host=cluster_endpoint.hostname, port=cluster_endpoint.port)

        secret.grant_read(function)
        projection_secret.grant_read(function)
        query_secret.grant_read(function)

        # Custom resource that will invoke the lambda on creation
        self.rds_custom_resource = cr.AwsCustomResource(self, "RDS_Initializer_CustomResource",
                                                   on_create=self.lambda_invokation(admin_secret=secret.secret_name, query_secret=query_secret.secret_name, project_secret=projection_secret.secret_name,
                                                                                    init_script=init_sql,
                                                                                    resourceId=function.function_name),
                                                   policy=cr.AwsCustomResourcePolicy.from_statements(
                                                       statements=[invoke_lambda_statement]),
                                                   role=lambda_role,

                                                   )

        self.function = function
        self.projection_secret = projection_secret
        self.query_secret = query_secret

    def lambda_invokation(self, admin_secret, query_secret, project_secret, init_script, resourceId):

        payload = {
            "admin_secret": admin_secret,
            "query_secret": query_secret,
            "projection_secret": project_secret,
            "sql-script": init_script
        },

        parameter = {
            "FunctionName": resourceId,
            "InvocationType": "RequestResponse",
            "Payload": json.dumps(payload)
        }

        return AwsSdkCall(
            action='invoke',
            service='Lambda',
            output_paths=[""],
            parameters=parameter,
            physical_resource_id=PhysicalResourceId.of(resourceId)
        )

    def create_database_secret(self, username: str, db_name: str, host: str, port: int):

        secret_json = {
            "username": username,
            "dbname": db_name,
            "host": host,
            "port": port
        }

        secret_value = secretsmanager.Secret(self,
                                             self.get_name(f"{username}-secret"),
                                             secret_name=self.get_name(f"{username}-secret"),
                                             generate_secret_string=secretsmanager.SecretStringGenerator(
                                                 generate_string_key="password", 
                                                 exclude_characters=" ^%+~`$&*()|[]{}:;,-<>?!'/\\\",=", 
                                                 exclude_punctuation=False, 
                                                 secret_string_template=(json.dumps(secret_json))
                                             ))

        return secret_value
