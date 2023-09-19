from aws_cdk import (
    aws_route53 as route53,
)
from constructs import Construct

from infrastructure.constructs.base.baseconstruct import BaseConstruct
from infrastructure.constructs.contextconfig import ContextConfig


class Route53Construct(BaseConstruct):

    def __init__(self, scope: Construct, id: str, api_custom_domain: str,
                 hosted_zone_name: str, domain_name: str,config: ContextConfig) -> None:
        super().__init__(scope, id, config=config)
        
        # Retrieve a Route53 Hosted Zone in the
        #hosted_zone = route53.HostedZone.from_lookup(self, "HostedZone", domain_name=hosted_zone_name)

        # Create a CNAME record set in the Hosted Zone to point to the API Gateway custom domain name
        # record_name = "test.logistic-backbone.com"
        # cname_record = route53.CnameRecord(self, "CnameRecord",
        #     record_name=record_name,
        #     domain_name=api_custom_domain,
        #     zone=hosted_zone)
        
        route53.CfnHealthCheck(self, "MyCfnHealthCheck",
            health_check_config=route53.CfnHealthCheck.HealthCheckConfigProperty(   
                type="HTTPS",
                fully_qualified_domain_name=domain_name,
                resource_path="/health",
                port=443
            ),

        )