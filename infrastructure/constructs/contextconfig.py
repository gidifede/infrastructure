"""
This Module implements the Stack Configuration
"""
from aws_cdk import(
    RemovalPolicy,
    App,
    aws_logs as logs,
    aws_dynamodb as dyndb,
)

removal_policy_map = {
    "DESTROY" : RemovalPolicy.DESTROY,
    "RETAIN" : RemovalPolicy.RETAIN,
    "SNAPSHOT" : RemovalPolicy.SNAPSHOT

}

billing_modes_map={
    "PROVISIONED": dyndb.BillingMode.PROVISIONED,
    "PAY_PER_REQUEST": dyndb.BillingMode.PAY_PER_REQUEST,

}
class ContextConfig():
    """
    The stack configuration
    """
    def __init__(self):
        """
        Empty init
        """
        self.env = None
        self.configuration={}

    def get(self, key:str, default=None):
        """  Gets a value from a key """
        return self.configuration.get(key,default)

    def get_resources_tags(self):
        """ Gets key and value form resources tags """
        key_list = list()
        value_list = list()
        resources_tags = self.configuration.get("resources-tags", [])
        for tag in resources_tags:
            key_list.append(tag.get("key"))
            value_list.append(tag.get("value"))
        return key_list, value_list

    def get_dynamodb_billing_mode(self):
        """ Translates DynamoDB billing mode """
        billing_mode  = self.get("dynamodb").get("billing-mode", "PAY_PER_REQUEST")
        return billing_modes_map[billing_mode]

    def get_dynamodb_table_removal_policy(self):
        """ Translates DynamoDB removal policy """
        removal_policy  = self.get("dynamodb").get("removal-policy", "RETAIN")
        return  removal_policy_map[removal_policy]

    def get_bucket_removal_policy(self):
        """ Translates S3 removal policy """
        removal_policy  = self.get("bucket-removal-policy", "RETAIN")
        return  removal_policy_map[removal_policy]

    def get_removal_policy(self,r_string: str):
        """ Translates removal policy """
        return  removal_policy_map[r_string]

    def get_log_retention(self):
        """ Translates Log retention """
        ret_map={
            "EIGHTEEN_MONTHS":logs.RetentionDays.EIGHTEEN_MONTHS,
            "FIVE_DAYS"      :logs.RetentionDays.FIVE_DAYS,
            "FIVE_MONTHS"    :logs.RetentionDays.FIVE_MONTHS,
            "FIVE_YEARS"     :logs.RetentionDays.FIVE_YEARS,
            "FOUR_MONTHS"    :logs.RetentionDays.FOUR_MONTHS,
            "INFINITE"       :logs.RetentionDays.INFINITE,
            "ONE_DAY"        :logs.RetentionDays.ONE_DAY,
            "ONE_MONTH"      :logs.RetentionDays.ONE_MONTH,
            "ONE_WEEK"       :logs.RetentionDays.ONE_WEEK,
            "ONE_YEAR"       :logs.RetentionDays.ONE_YEAR,
            "SIX_MONTHS"     :logs.RetentionDays.SIX_MONTHS,
            "TEN_YEARS"      :logs.RetentionDays.TEN_YEARS,
            "THIRTEEN_MONTHS":logs.RetentionDays.THIRTEEN_MONTHS,
            "THREE_DAYS"     :logs.RetentionDays.THREE_DAYS,
            "THREE_MONTHS"   :logs.RetentionDays.THREE_MONTHS,
            "TWO_MONTHS"     :logs.RetentionDays.TWO_MONTHS,
            "TWO_WEEKS"      :logs.RetentionDays.TWO_WEEKS,
            "TWO_YEARS"      :logs.RetentionDays.TWO_YEARS,

        }
        retention = self.configuration.get('log-retention-days', 'ONE_DAY')
        return ret_map.get(retention)
    
    def get_command_handler_idempotency_check_source(self):
        check_source  = self.get("command-handler").get("idempotency-check-source", "false")
        return  check_source

    def get_command_handler_state_machine_config_origin(self):
        origin  = self.get("command-handler").get("state-machine-config").get("config-origin", "configuration/stateMachine")
        return  origin
    
    def get_command_handler_state_machine_config_prefix(self):
        prefix  = self.get("command-handler").get("state-machine-config").get("config-bucket-prefix", "commandHandler")
        return  prefix
    
    def get_command_validator_json_schema_config_origin(self):
        origin  = self.get("command-validator").get("json-schema-config").get("config-origin", "configuration/JSON/JSONSCHEMA")
        return  origin
    
    def set_command_validator_json_schema_config_origin(self, originPath):
        self.configuration["command-validator"]["command-api-map-config"]["config-origin"] = originPath
    
    def get_command_validator_json_schema_config_prefix(self):
        prefix  = self.get("command-validator").get("json-schema-config").get("config-bucket-prefix", "product/commandValidator/v1/JSONSCHEMA")
        return  prefix
    
    def get_command_validator_command_api_map_config_origin(self):
        origin  = self.get("command-validator").get("command-api-map-config").get("config-origin", "swagger/commands/v1/")
        return  origin
    
    def get_command_validator_command_api_map_config_prefix(self):
        prefix  = self.get("command-validator").get("command-api-map-config").get("config-bucket-prefix", "product/commandValidator/command-api-map")
        return  prefix

    def get_projection_db_projection_credential_secret_name(self):
        db_projection_credential_secret_name = self.get("projection").get("db_projection_credential_secret_name", "product_projection_application_user")
        return db_projection_credential_secret_name

    def get_projection_region(self):
        region = self.get("projection").get("region", "eu-central-1")
        return region

    
    def get_db_queryhandler_credential_secret_name(self):
        db_queryhandler_credential_secret_name = self.get("query-handler").get("db_queryhandler_credential_secret_name", "product_queryhandler_application_user")
        return db_queryhandler_credential_secret_name

    def get_queryhandler_region(self):
        region = self.get("query-handler").get("region", "eu-central-1")
        return region

    
    def get_lambda_build_image(self):
        lambda_build_image  = self.get("lambda").get("build-image", "golang:1.19")
        return  lambda_build_image
    
    def set_api_key(self, api_key):
        self.configuration["api"]["api-key"]["value"] = api_key

    @classmethod
    def build_config(cls,app: App):
        """ Builds the configuration from the cdk.json file"""
        env = app.node.try_get_context("config")
        if env is None:
            raise Exception("please provide a value for the configuration with -c config=<name>")

        ctxcfg = ContextConfig()
        env_params = app.node.try_get_context(env)
        ctxcfg.env = env
        ctxcfg.configuration = ctxcfg.configuration | env_params

        return ctxcfg
