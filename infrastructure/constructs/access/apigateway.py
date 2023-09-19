'''
API Gateway construct module
'''
from aws_cdk import (RemovalPolicy,
                     aws_iam as iam,
                     aws_apigateway as apigw,
                     aws_logs as logs,
                     aws_ec2 as ec2,
                     aws_certificatemanager as cert,
                     aws_route53 as route53,
                     aws_route53_targets as targets
                     )
from constructs import Construct

from infrastructure.constructs.base.baseconstruct import BaseConstruct
from infrastructure.constructs.contextconfig import ContextConfig


class ApiGatewayConstruct(BaseConstruct):
    '''
    Construct to create the API gateway
    '''
    def __init__(self,
                 scope: Construct,
                 id: str,
                 rest_api_name: str,
                 rest_api_description: str,
                 stage_name: str,
                 config: ContextConfig) -> None:
        super().__init__(scope, id, config=config)

        access_logs_enabled = config.get("api", {}).get("access-logs-enabled", False)

        if access_logs_enabled:
            #Create Log Group
            log_group = logs.LogGroup(self,self.get_name("ApiGatewayAccessLog"),
                                    log_group_name=self.get_name("ApiGatewayAccessLog"),
                                    removal_policy=RemovalPolicy.DESTROY,
                                    retention = logs.RetentionDays.ONE_MONTH)
            #Crea un ruolo utilizzando le policies definite da amazon
            # le quali permettono di scrivere i log su cloudwatch
            role = iam.Role(self, self.get_name("CWRole"),
                        assumed_by=iam.ServicePrincipal("apigateway.amazonaws.com"),
                        managed_policies=[
                            iam.ManagedPolicy.from_aws_managed_policy_name("service-role/AmazonAPIGatewayPushToCloudWatchLogs")
                            ]
                        )
            #Configura il ruolo, creato sopra, come ruolo comune a tutte le api dell'apigateway
            apigw.CfnAccount(self, self.get_name("MyCfnAccount"), cloud_watch_role_arn=role.role_arn)
            access_log_destination = apigw.LogGroupLogDestination(log_group)
            access_log_format = apigw.AccessLogFormat.json_with_standard_fields(caller=True,          #The principal identifier of the caller will be output to the log.
                                                                                http_method=True,     #The http method will be output to the log.
                                                                                ip=True,              #The source IP of request will be output to the log.
                                                                                protocol=True,        #The request protocol will be output to the log.
                                                                                request_time=True,    #The CLF-formatted request time((dd/MMM/yyyy:HH:mm:ss +-hhmm) will be output to the log.
                                                                                resource_path=True,   #The path to your resource will be output to the log.
                                                                                response_length=True, #The response payload length will be output to the log.
                                                                                status=True,          #The method response status will be output to the log.
                                                                                user=True             #The principal identifier of the user will be output to the log.
                                                                                )
        else:
            access_log_destination = None
            access_log_format = None

        endpoint_configuration=apigw.EndpointConfiguration(types=[apigw.EndpointType.REGIONAL])
        api_policy = None
        if config.get("vpc").get("enabled"):
            api_policy = iam.PolicyDocument(
                statements=[
                    iam.PolicyStatement(
                            actions = ['execute-api:Invoke'],
                            principals =  [iam.AnyPrincipal()],
                            resources= ['execute-api:/*/*/*'],
                    ),
                    iam.PolicyStatement(
                            effect= iam.Effect.ALLOW,
                            principals= [iam.AnyPrincipal()],
                            resources= ['execute-api:/*/*/*'],
                    )
                ]
            )

            vpcendpoint = ec2.InterfaceVpcEndpoint.from_interface_vpc_endpoint_attributes(
                scope=self,
                id="VPC",
                port=443,
                vpc_endpoint_id=config.get("vpc").get("vpce_id"))

            endpoint_configuration=apigw.EndpointConfiguration(
                types=[apigw.EndpointType.PRIVATE],
                vpc_endpoints=[vpcendpoint])

        custom_domain_enabled = config.get("api", {}).get("custom-domain", {}).get("enabled", True)

        custom_domain = None
        if custom_domain_enabled:
            domain_name = config.get("api", {}).get("custom-domain", {}).get("domain-name", None)
            full_domain_name = f"{self.stack_name.lower()}.{domain_name}"

            certificate_arn = config.get("api", {}).get("certificate-arn", None)
            certificate = cert.Certificate.from_certificate_arn(self,
                                                                self.get_name("imported-cert"),
                                                                certificate_arn=certificate_arn)

            custom_domain = apigw.DomainNameOptions(
                domain_name=full_domain_name,
                certificate=certificate,
                security_policy=apigw.SecurityPolicy.TLS_1_2,
                endpoint_type=apigw.EndpointType.REGIONAL,
            )

        # Create API Gateway
        deploy_options = apigw.StageOptions(stage_name=stage_name,
                                            data_trace_enabled= True,                      #When enabled, API gateway will log the full API requests and responses. This can be useful to troubleshoot APIs, but can result in logging sensitive data. We recommend that you don’t enable this feature for production APIs.
                                            logging_level = apigw.MethodLoggingLevel.INFO, #INFO,ERROR,OFF
                                            metrics_enabled= True,                         #Each method will generate these metrics: API calls, Latency, Integration latency, 400 errors, and 500 errors. The metrics are charged at the standard CloudWatch rates.
                                            tracing_enabled= True,                         #Specifies whether Amazon X-Ray tracing is enabled
                                            access_log_destination = access_log_destination,
                                            access_log_format = access_log_format
                                        )

        rest_api = apigw.RestApi(self,
                                 id=id,
                                 cloud_watch_role=False,
                                 rest_api_name=rest_api_name,
                                 description=rest_api_description,
                                 endpoint_configuration=endpoint_configuration,
                                 policy=api_policy,
                                 domain_name=custom_domain,
                                 deploy_options=deploy_options
                                )

        usage_plan = rest_api.add_usage_plan(f"{rest_api.rest_api_name}UsagePlan",
                                             name=f"{rest_api.rest_api_name}UsagePlan",
                                             api_stages=[apigw.UsagePlanPerApiStage(
                                                 api=rest_api, stage=rest_api.deployment_stage)]
                                             )

        if custom_domain_enabled:
            domain_name = config.get("api", {}).get("custom-domain", {}).get("route53-hosted-zone", None)
            full_domain_name = f"{self.stack_name.lower()}.{domain_name}"
            hosted_zone = route53.HostedZone.from_lookup(self,
                                                         self.get_name("HostedZone"),
                                                         domain_name=domain_name)
            route53.ARecord(
                self,
                self.get_name("ApiRecord"),
                record_name=full_domain_name,
                zone=hosted_zone,
                target=route53.RecordTarget.from_alias(targets.ApiGateway(rest_api)),
            )

        if config.get("api").get("api-key").get("enabled", False):
            common_api_key = config.get("api").get("api-key").get("id", None)
            imported_key = apigw.ApiKey.from_api_key_id(self,
                                                    self.get_name("imported-key"),
                                                    common_api_key)
            usage_plan.add_api_key(api_key=imported_key)

        self.api = rest_api
