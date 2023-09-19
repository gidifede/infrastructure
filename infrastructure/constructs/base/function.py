"""
This module implements a base Function for every lambda in the stack
"""
from collections.abc import Mapping, Sequence
from aws_cdk import(
    Stack,
    Duration,
    aws_lambda as lambda_,
    aws_iam as iam,
    aws_s3 as s3,
    aws_ec2 as ec2,
    aws_efs as efs,
    aws_sns as sns,
    aws_sqs as sqs,
    aws_secretsmanager as secretesmanager,
    BundlingOptions,
    DockerImage,
    ILocalBundling
)
from constructs import Construct
from infrastructure.constructs import constants
from infrastructure.constructs.contextconfig import ContextConfig
import os
import jsii
from typing import Optional

class GoFunction(Construct):
    """
    Function Construct.

    Parameters
    ----------

    name: str
       Lambda's name

    vpc: ec2.Vpc
       VPC

    fs: efs.FileSystem
       EFS shared filesystem

    lambda_env: dict
       additional lambda environment variables

    layers: Sequence[lambda_.LayerVersion]
       additional layers

    description: str
       this lambda description

    handler: str
       this lambda handler

    asset: str
       this lambda asset dir

    handler: str
       the Lambda handler
    """
    def __init__(self,
        scope: Construct,
        construct_id: str,
        config: ContextConfig,
        name: str,
        handler: str,
        vpc_id : Optional[str]=None,
        # vpc_id: str=None,
        fs: efs.FileSystem=None,
        lambda_env: Mapping[str,str]=None,
        layers: Sequence[lambda_.LayerVersion]=None,
        description: str=None,
        asset: str= None,
        **kwargs) -> None:
        super().__init__(scope, construct_id, **kwargs)

        if lambda_env is None:
            lambda_env = {}

        if layers is None:
            layers = []

        environment = config.env
        common_env={
          "LOG_LEVEL": config.get("loglevel", "DEBUG"),
          "ENVIRONMENT": environment.upper()
        }

        stack_name = Stack.of(self).stack_name

        # Get VPC
        vpc = None
        security_groups = []
        if vpc_id is not None:

            vpc = ec2.Vpc.from_vpc_attributes(self, "VPC",
                                            vpc_id=vpc_id,
                                            # vpc_id=config.get("vpc").get("lambda").get("id"),
                                            availability_zones=config.get("vpc").get("lambda").get("availability_zones"),
                                            private_subnet_ids=config.get("vpc").get("lambda").get("private_subnet_ids"),
                                            private_subnet_route_table_ids=config.get("vpc").get("lambda").get("private_subnet_route_table_ids"))
            for security_group in config.get("vpc").get("lambda").get("security_groups",[]):
                security_groups.append(ec2.SecurityGroup.from_security_group_id(self,id=f"{security_group}-sg",security_group_id=security_group))


        print(f"PROJECTION_NAME: {GoFunction.get_name(stack_name=stack_name, name=name)}")


        if asset is not None:
            self.function = lambda_.Function(self,
                                    construct_id,
                                    function_name=GoFunction.get_name(stack_name=stack_name, name=name),
                                    tracing=lambda_.Tracing.ACTIVE,
                                    runtime=lambda_.Runtime.GO_1_X,
                                    security_groups=security_groups,
                                    code=lambda_.Code.from_asset(asset, bundling=BundlingOptions(
                                        # image=DockerImage.from_registry('postesviluppo.azurecr.io/caat-images/caat-ubi-golang:8'),
                                        image=DockerImage.from_registry(config.get_lambda_build_image()),
                                        command=[
                                            'bash', '-c',
                                            # Need to set GOCACHE because, by default, GOCACHE=/.cache/go-build, but user in docker container (unless you specify root user) 
                                            # doesn't have write privileges on /.cache.

                                            # You can set privileges with a chmod on desired folder or set GOCACHE to a folder where user has privileges. 
                                            # Look at https://github.com/aws/aws-cdk/issues/8707.

                                            # When https://docs.aws.amazon.com/cdk/api/v2/docs/aws-lambda-go-alpha-readme.html will go in GA, 
                                            # there will be no need to set GOCACHE, since it's done by the library. Look at https://github.com/aws/aws-cdk/blob/main/packages/%40aws-cdk/aws-lambda-go/lib/Dockerfile 
                                            'GOCACHE=$GOPATH/.cache/go-build GOOS=linux go build -o /asset-output/{0} ./cmd/main.go'.format(handler),
                                        ],
                                        local= LocalBundle(wrk_dir=asset, handler=handler)
                                    )),
                                    
                                    # module_dir=asset,
                                    handler=handler,
                                    environment=common_env | lambda_env ,
                                    #layers=layers,
                                    vpc=vpc,
                                    filesystem=fs,
                                    description=description,
                                    log_retention=config.get_log_retention(),
                                    layers=layers,
                                    timeout=Duration.seconds(config.get('lambda-timeout-seconds',10)))
            
        else:
            raise Exception("asset must be provided")


    @classmethod
    def get_name(cls,stack_name:str, name: str):
        """ gets lambda complete name"""
        return f"{stack_name}-{name}"


    def grant_publish_on(self, topic: sns.Topic):
        """ grants publish on sns topic """
        self.function.role.add_to_policy(
            iam.PolicyStatement(
                actions=['sns:Publish'],
                effect=iam.Effect.ALLOW,
                resources=[topic.topic_arn]
            )
        )

    def grant_invoke_on(self, function: lambda_.Function):
        """ grants invoke on another lambda """
        self.function.role.add_to_policy(
            iam.PolicyStatement(
                effect=iam.Effect.ALLOW,
                actions=[
                    'lambda:InvokeFunction',
                    'lambda:InvokeAsync',
                ],
                resources=[
                    function.function_arn,
                ]
            )
        )

    def grant_read_on_bucket(self, bucket: s3.Bucket, path: str=""):
        """ Grants read on a bucket"""
        self.function.role.add_to_policy(
            iam.PolicyStatement(
                effect=iam.Effect.ALLOW,
                actions=[
                    's3:GetObject',
                    's3:GetObjectVersion',
                    's3:GetObjectAttributes',
                    's3:GetObjectTagging',
                    's3:GetObjectVersionAttributes',
                    's3:GetObjectRetention',
                    's3:ListBucket',
                ],
                resources=[
                    f"{bucket.bucket_arn}",
                    f"{bucket.bucket_arn}{path}/*",
                ]
            )
        )

    def grant_secret_read(self, secret: secretesmanager.Secret):
        """ Grants read on a secret """
        self.function.role.add_to_policy(
            iam.PolicyStatement(
                effect=iam.Effect.ALLOW,
                actions=[
                    "secretsmanager:DescribeSecret",
                    "secretsmanager:GetSecretValue",
                ],
                resources=[
                    secret.secret_arn,
                ]
            )
        )

    def grant_write_on_bucket(self, bucket: s3.Bucket, path: str=""):
        """ Grants write on a bucket"""
        self.function.role.add_to_policy(
            iam.PolicyStatement(
                effect=iam.Effect.ALLOW,
                actions=[
                    's3:PutObject',
                    's3:PutObjectTagging',
                    's3:PutBucketObjectLockConfiguration',
                    's3:PutObjectLegalHold',
                    's3:PutObjectRetention',
                    's3:GetBucketObjectLockConfiguration',
                ],
                resources=[
                    f"{bucket.bucket_arn}",
                    f"{bucket.bucket_arn}{path}/*",
                ]
            )
        )

    def grant_delete_on_bucket(self, bucket: s3.Bucket, path: str=""):
        """ Grants delete on a bucket """
        self.function.role.add_to_policy(
            iam.PolicyStatement(
                effect=iam.Effect.ALLOW,
                actions=[
                    's3:DeleteObject',
                    's3:DeleteObjectVersion',
                    's3:PutLifeCycleConfiguration'
                ],
                resources=[
                    f"{bucket.bucket_arn}",
                    f"{bucket.bucket_arn}{path}/*",
                ]
            )
        )

    def grant_readwrite_on_bucket(self, bucket: s3.Bucket, path: str=""):
        """ Grants read/write on a bucket """
        self.function.role.add_to_policy(
            iam.PolicyStatement(
                effect=iam.Effect.ALLOW,
                actions=[
                    's3:GetObject',
                    's3:GetObjectVersion',
                    's3:GetObjectAttributes',
                    's3:GetObjectTagging',
                    's3:GetObjectVersionAttributes',
                    's3:GetObjectRetention',
                    's3:GetObjectLegalHold',
                    's3:ListBucket',
                    's3:PutObject',
                    's3:PutObjectTagging',
                    's3:PutObjectRetention',
                    's3:PutBucketObjectLockConfiguration',
                    's3:PutObjectLegalHold',
                    's3:GetBucketObjectLockConfiguration',
                ],
                resources=[
                    f"{bucket.bucket_arn}",
                    f"{bucket.bucket_arn}{path}/*",
                ]
            )
        )

    def grant_generate_data_key(self, encryption_key_arn: str):
        """ Grants generate data key on a key """
        self.function.role.add_to_policy(
            iam.PolicyStatement(
                effect=iam.Effect.ALLOW,
                actions=[
                    'kms:GenerateDataKey',
                ],
                resources=[
                    encryption_key_arn
                ]
            )
        )

    def grant_decrypt(self, encryption_key_arn: str):
        """ Grants decrypt on a key """
        self.function.role.add_to_policy(
            iam.PolicyStatement(
                effect=iam.Effect.ALLOW,
                actions=[
                    'kms:Decrypt',
                ],
                resources=[
                    encryption_key_arn
                ]
            )
        )

    def grant_step_function_execute(self, step_function_arn: str):
        """ Grants execute on a step function """
        self.function.role.add_to_policy(
            iam.PolicyStatement(
                effect=iam.Effect.ALLOW,
                actions=[
                    "states:StartExecution"
                ],
                resources=[
                    step_function_arn
                ]
            )
        )

@jsii.implements(ILocalBundling)
class LocalBundle:
    def __init__(self, wrk_dir: str, handler: str)-> None:
        self.wrk_dir = wrk_dir
        self.handler = handler
        super().__init__()

    def try_bundle(self, output_dir, *, image, entrypoint=None, command=None, volumes=None, volumesFrom=None, environment=None, workingDirectory=None, user=None, local=None, outputType=None, securityOpt=None, network=None):
        try:
            os.system('go version')
        except:
            return False
        try:
            os.system('cd {0} && GOOS=linux go build -o {1}/{2} ./cmd/main.go'.format(self.wrk_dir, output_dir, self.handler))
        except:
            return False
        return True
