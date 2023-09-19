from aws_cdk import (
    aws_ec2 as ec2,
    aws_docdb as docdb
)
import aws_cdk as cdk


from infrastructure.constructs.base.baseconstruct import BaseConstruct
from infrastructure.constructs.base.baseconstruct import Construct
from infrastructure.constructs.contextconfig import ContextConfig

class ProjectionDatabase(BaseConstruct):

    def __init__(self, scope: Construct, id: str,
                 config: ContextConfig, **kwargs) -> None:
        super().__init__(scope, id, config, **kwargs)

        vpc_id = config.get("vpc", {}).get("services", {}).get("id", None)
        subnets = []
        vpc = ec2.Vpc.from_lookup(self, "VPC", vpc_id=vpc_id)

        for subnet in config.get("vpc", {}).get("services", {}).get("private_subnet_ids", []):
            subnets.append(
                ec2.Subnet.from_subnet_id(self, f"Subnet-{subnet}", subnet)
            )
            
        sg = []
        for security_group in config.get("vpc").get("services").get("security_groups", []):
            sg.append(ec2.SecurityGroup.from_security_group_id(
                self, id=f"{security_group}-sg", security_group_id=security_group))
            
        

        SECRET_NAME = f"{self.stack_name}-docdb-secret"
        self.cluster = docdb.DatabaseCluster(self, "Database",
            master_user=cdk.aws_docdb.Login(
                username="docdb",
                secret_name=SECRET_NAME
            ),
            instance_type=ec2.InstanceType.of(ec2.InstanceClass.R5, ec2.InstanceSize.LARGE),
            vpc=vpc,
            security_group=sg[0],
            removal_policy=cdk.RemovalPolicy.DESTROY
        )
        # cluster.add_depends_on(subnetGroup)
        self.cluster.add_rotation_single_user()
     

        
