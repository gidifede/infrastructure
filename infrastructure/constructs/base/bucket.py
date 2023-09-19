"""
This Module implements a Base Bucket with object lock
"""
import json
from collections.abc import Sequence
import datetime
import hashlib
import base64
import isodate
from aws_cdk import(
    Stack,
    RemovalPolicy,
    Duration,
    aws_s3 as s3,
    aws_backup as backup,
    aws_iam as iam,
    aws_kms as kms,
    aws_events as events,
    aws_cloudtrail as trail,
    aws_s3_deployment as s3d
)
from constructs import Construct
from infrastructure.constructs.contextconfig import ContextConfig
from infrastructure.constructs import constants
from infrastructure.constructs.base.baseconstruct import BaseConstruct

class Cors:
    """
    Utility class to implement Cors
    """
    def __init__(self,
                 enabled: bool,
                 allowed_headers: Sequence[str],
                 allowed_origins: Sequence[str],
                 allowed_methods: Sequence[str],
                 exposed_headers: Sequence[str],
                 max_age: int
                 ):
        self.enabled = enabled
        self.allowed_headers = allowed_headers
        self.allowed_origins = allowed_origins
        self.allowed_methods = allowed_methods
        self.exposed_headers=exposed_headers
        self.max_age = max_age

def get_start_window(entry):
    """ Gets start window in hours """
    if entry is None:
        return None
    return Duration.hours(entry)

def get_duration(entry):
    """ Gets duration in days from configuration or None"""
    if entry is None:
        return None
    return Duration.days(entry)

def get_duration_in_days(d_string: str):
    """ Returns a duration in days from a string"""
    isod = f'P{d_string.upper()}'
    duration = isodate.parse_duration(isod)

    if not isinstance(duration, datetime.timedelta):
        return int(duration.years * 365)
    return int(duration.days)

def get_lifecycle_rules():
    """ Reads the lifecycle rules from a file and returns a L1 Construct """
    with open("init/storage.json",encoding="utf-8") as lf_file:
        storage_classes = json.loads(lf_file.read())
        storage_map = { storage_class['name']:storage_class for storage_class in storage_classes}

        lifecycle_rules = []
        for lf_rule in storage_map:
            retention_days = get_duration_in_days(storage_map[lf_rule]['retentionPeriod'])
            hot_days       = get_duration_in_days(storage_map[lf_rule]['hotPeriod'])

            transitions = []
            # no transition if retention and hot are equal
            if retention_days != hot_days:
                transitions=[
                    s3.CfnBucket.TransitionProperty(
                        storage_class="GLACIER",
                        transition_in_days=hot_days
                    )
                ]

            rule = s3.CfnBucket.RuleProperty(
                id=lf_rule,
                status="Enabled",
                expiration_in_days=retention_days,
                tag_filters=[ s3.CfnBucket.TagFilterProperty( key=constants.TAG_LIFECYCLE_RULE, value=lf_rule) ],
                transitions=transitions
            )
            lifecycle_rules.append(rule)

        return lifecycle_rules

class Bucket(BaseConstruct):
    """
    This class implements the Base Bucket by using the L1 construct

    Parameters
    ----------

    bucket_name: str
        Bucket Name, it will be set by the classmethod get_name

    config: ContextConfig
        The Stack configuration

    object_lock_enabled: bool
        True if the object lock must be enabled for this bucket

    object_lock_retention_days: int
        default object lock retention days

    object_lock_mode: str
        default object lock mode
    """
    def __init__(self,
        scope: Construct,
        construct_id: str,
        bucket_name: str,
        config: ContextConfig,
        cors: Cors=None,
        object_lock_enabled: bool=False,
        object_lock_retention_days: int=10,
        object_lock_mode: str="GOVERNANCE",
        backup_enabled: bool=False,
        logging_enabled: bool=False,
        log_prefix: str="safe-storage/",
        **kwargs) -> None:
        super().__init__(scope, construct_id,config=config, **kwargs)


        self.bucket_name = bucket_name
        object_lock_configuration_property = None
        versioning_configuration_property = None
        encryption_configuration_property = None

        # object lock
        if object_lock_enabled:
            object_lock_configuration_property = s3.CfnBucket.ObjectLockConfigurationProperty(
                object_lock_enabled="Enabled",
                rule=s3.CfnBucket.ObjectLockRuleProperty(
                    default_retention=s3.CfnBucket.DefaultRetentionProperty(
                        days=object_lock_retention_days,
                        mode=object_lock_mode,
                    )
                )
            )
            versioning_configuration_property = s3.CfnBucket.VersioningConfigurationProperty(
               status="Enabled"
            )

        # encryption
        encryption_key_arn = config.get('encryption-key-arn',None)
        if encryption_key_arn is not None:
            encryption_configuration_property = s3.CfnBucket.BucketEncryptionProperty(
                server_side_encryption_configuration=[
                    s3.CfnBucket.ServerSideEncryptionRuleProperty(
                         bucket_key_enabled=True,
                         server_side_encryption_by_default=s3.CfnBucket.ServerSideEncryptionByDefaultProperty(
                             sse_algorithm="aws:kms",
                             kms_master_key_id=encryption_key_arn
                         )
                    )
                ]
            )

        # cors
        cors_cfg = None
        if cors is not None and cors.enabled:
            cors_cfg = s3.CfnBucket.CorsConfigurationProperty(
                cors_rules = [s3.CfnBucket.CorsRuleProperty(
                    allowed_headers=cors.allowed_headers,
                    allowed_methods=cors.allowed_methods,
                    allowed_origins=cors.allowed_origins,
                    exposed_headers=cors.exposed_headers,
                    max_age=cors.max_age
                )]
            )

        # getting lifecycle rules
        # rules = get_lifecycle_rules()

        # full bucket name
        full_bucket_name = Bucket.get_name(stack_name=Stack.of(self).stack_name,
                                              region=Stack.of(self).region,
                                              bucket_name=self.bucket_name)
        # enabling logging if requested
        logging_configuration_property = None
        if logging_enabled:
            encryption_key = None
            if encryption_key_arn is not None:
                encryption_key = kms.Key.from_key_arn(self,"Key", encryption_key_arn)
            log_bucket_name = Bucket.get_name(stack_name=Stack.of(self).stack_name,
                                              region=Stack.of(self).region,
                                              bucket_name=f"{self.bucket_name}-log",
                                             )
            logging_bucket = s3.Bucket(self,f"{construct_id}Log",
                                               bucket_name=log_bucket_name,
                                               block_public_access=s3.BlockPublicAccess.BLOCK_ALL,
                                               encryption_key=encryption_key,
                                               )
            log_resource = log_prefix.rstrip("/")

            logging_bucket.add_to_resource_policy(
                iam.PolicyStatement(
                    effect=iam.Effect.ALLOW,
                    principals=[iam.ServicePrincipal("logging.s3.amazonaws.com")],
                    actions=[
                        "s3:PutObject"
                    ],
                    resources=[
                        f"arn:aws:s3:::{log_bucket_name}/{log_resource}*"
                    ]
                )
            )


            logging_configuration_property = s3.CfnBucket.LoggingConfigurationProperty(destination_bucket_name=logging_bucket.bucket_name,
                                                                                       log_file_prefix=log_prefix)



        public_access_block_configuration_property = s3.CfnBucket.PublicAccessBlockConfigurationProperty(
                                                        block_public_acls=True,
                                                        block_public_policy=True,
                                                        ignore_public_acls=True,
                                                        restrict_public_buckets=True
                                                       )

        cfn_bucket = s3.CfnBucket(self, f"Cfn{construct_id}",
                  bucket_name=full_bucket_name,
                  object_lock_enabled=object_lock_enabled,
                  object_lock_configuration=object_lock_configuration_property,
                  versioning_configuration=versioning_configuration_property,
                  bucket_encryption=encryption_configuration_property,
                #   lifecycle_configuration=s3.CfnBucket.LifecycleConfigurationProperty(rules=rules),
                  cors_configuration=cors_cfg,
                  logging_configuration=logging_configuration_property,
                  public_access_block_configuration=public_access_block_configuration_property
        )

        # cfn_bucket.apply_removal_policy(RemovalPolicy.RETAIN)
        cfn_bucket.apply_removal_policy(config.get_bucket_removal_policy())
        
        # Enabling backup
        if backup_enabled:
            self.queue_encrypted_custom_key = True
            encryption_key = None
            if constants.BACKUP_PLAN is None: # create backup plan
                if encryption_key_arn is not None:
                    encryption_key = kms.Key.from_key_arn(self,"Key", encryption_key_arn)
                constants.BACKUP_ROLE = iam.Role(self, "BackupRole",
                                           role_name=f"{Stack.of(self).stack_name}BackupRole",
                                           assumed_by=iam.ServicePrincipal('backup.amazonaws.com'),
                                           managed_policies=[
                                                iam.ManagedPolicy.from_aws_managed_policy_name("service-role/AWSBackupServiceRolePolicyForBackup"),
                                                iam.ManagedPolicy.from_aws_managed_policy_name("service-role/AWSBackupServiceRolePolicyForRestores"),

                                           ],
                                      )
                self.attach_policies(constants.BACKUP_ROLE)

                backup_vault = backup.BackupVault(self, "SafeStorageBackupVault",
                                                   backup_vault_name=f"{Stack.of(self).stack_name}BackupVault",
                                                   encryption_key=encryption_key,
                                                   removal_policy=config.get_bucket_removal_policy())
                # build backup plan rules
                backup_plan_rules = []
                for rule in config.get("backup").get("rules"):
                    name = rule.get("name")
                    delete_after = get_duration(rule.get("delete_after_days",None))
                    move_to_cold_storage = get_duration(rule.get("move_to_cold_storage_after_days",None))
                    schedule     = rule.get("schedule",None)
                    continuous = rule.get("continuous", False)
                    start_window = get_start_window(rule.get("start_window_hours"))

                    if schedule is not None:
                        backup_plan_rules.append( backup.BackupPlanRule(
                            rule_name=  name,
                            delete_after=delete_after,
                            schedule_expression=events.Schedule.expression(schedule),
                            enable_continuous_backup=continuous,
                            move_to_cold_storage_after=move_to_cold_storage,
                            backup_vault=backup_vault,
                            start_window=start_window
                        ))
                    else:
                        backup_plan_rules.append( backup.BackupPlanRule(
                            rule_name=  name,
                            delete_after=delete_after,
                            enable_continuous_backup=continuous,
                            move_to_cold_storage_after=move_to_cold_storage,
                            backup_vault=backup_vault,
                            start_window=start_window

                        ))

                constants.BACKUP_PLAN = backup.BackupPlan(self, "SafeStorageBackupPlan",
                                            backup_plan_name=f"{Stack.of(self).stack_name}BackupPlan",
                                            backup_plan_rules=backup_plan_rules,
                                        )
            constants.BACKUP_PLAN.add_selection(
                f"{bucket_name}Backup",
                role=constants.BACKUP_ROLE,
                resources=[ backup.BackupResource.from_arn(cfn_bucket.attr_arn)]
            )

        self.bucket = s3.Bucket.from_bucket_arn(self,construct_id,bucket_arn=cfn_bucket.attr_arn)




    @classmethod
    def get_name(cls, stack_name:str, region: str, bucket_name:str):
        """ Class method to return the full bucket name"""

        name =  f"{stack_name.lower()}-{bucket_name.lower()}-{region}"
        if len(name) >= 63:
            m_hash = hashlib.sha256()
            m_hash.update(name.encode())
            suffix = (base64.b64encode(m_hash.digest())).decode()[0:12]
            suffix = suffix.replace("+","-")
            suffix = suffix.replace("/","-")
            suffix = suffix.lower()
            name=name[0:50]+suffix

        return name

    def add_asset(self, prefix: str, path: str, key = ""):
        stack_name = Stack.of(self).stack_name
        s3d.BucketDeployment(Stack.of(self), f"{stack_name}-Upload configuration {path}{key}",
            sources=[s3d.Source.asset(path)],
            destination_bucket=self.bucket,
            destination_key_prefix=prefix
        )

    def attach_policies(self,role: iam.Role):
        """ Adds backup policies to role """
        backup_statements=[
            iam.PolicyStatement(
                effect=iam.Effect.ALLOW,
                actions=[
                    "s3:GetInventoryConfiguration",
                    "s3:PutInventoryConfiguration",
                    "s3:ListBucketVersions",
                    "s3:ListBucket",
                    "s3:GetBucketVersioning",
                    "s3:GetBucketNotification",
                    "s3:PutBucketNotification",
                    "s3:GetBucketLocation",
                    "s3:GetBucketTagging"
                ],
                resources=[
                    "arn:aws:s3:::*"
                ]
            ),
            iam.PolicyStatement(
                effect=iam.Effect.ALLOW,
                actions=[
                    "s3:GetObjectAcl",
                    "s3:GetObject",
                    "s3:GetObjectVersionTagging",
                    "s3:GetObjectVersionAcl",
                    "s3:GetObjectTagging",
                    "s3:GetObjectVersion"
                ],
                resources=[
                        "arn:aws:s3:::*/*"
                ]
            ),
            iam.PolicyStatement(
                effect=iam.Effect.ALLOW,
                actions=[
                    "s3:ListAllMyBuckets"
                ],
                resources=[
                    "*"
                ]
            ),
            iam.PolicyStatement(
                effect=iam.Effect.ALLOW,
                actions=[
                    "kms:Decrypt",
                    "kms:DescribeKey"
                ],
                resources=[
                    "*"
                ],
                conditions={
                    "StringLike":{
                        "kms:ViaService":"s3.*.amazonaws.com"
                    }
                }
            ),
            iam.PolicyStatement(
                effect=iam.Effect.ALLOW,
                actions=[
                    "events:DescribeRule",
                    "events:EnableRule",
                    "events:PutRule",
                    "events:DeleteRule",
                    "events:PutTargets",
                    "events:RemoveTargets",
                    "events:ListTargetsByRule",
                    "events:DisableRule"
                ],
                resources=[
                    "arn:aws:events:*:*:rule/AwsBackupManagedRule*"
                ]
            ),
            iam.PolicyStatement(
                effect=iam.Effect.ALLOW,
                actions=[
                    "cloudwatch:GetMetricData",
                    "events:ListRules"
                ],
                resources=[
                    "*"
                ]
            )
        ]
        restore_statements=[
            iam.PolicyStatement(
                effect=iam.Effect.ALLOW,
                actions=[
                    "s3:CreateBucket",
                    "s3:ListBucketVersions",
                    "s3:ListBucket",
                    "s3:GetBucketVersioning",
                    "s3:GetBucketLocation",
                    "s3:PutBucketVersioning"
                ],
                resources=[
                    "arn:aws:s3:::*"
                ]
            ),
            iam.PolicyStatement(
                effect=iam.Effect.ALLOW,
                actions=[
                    "s3:GetObject",
                    "s3:GetObjectVersion",
                    "s3:DeleteObject",
                    "s3:PutObjectVersionAcl",
                    "s3:GetObjectVersionAcl",
                    "s3:GetObjectTagging",
                    "s3:PutObjectTagging",
                    "s3:GetObjectAcl",
                    "s3:PutObjectAcl",
                    "s3:PutObject",
                    "s3:ListMultipartUploadParts"
                ],
                resources=[
            "arn:aws:s3:::*/*"
                ]
            ),
            iam.PolicyStatement(
                effect=iam.Effect.ALLOW,
                actions=[
                    "kms:Decrypt",
                    "kms:DescribeKey",
                    "kms:GenerateDataKey"
                ],
                resources=[
                    '*'
                ],
                conditions={
                    "StringLike":{
                        "kms:ViaService":"s3.*.amazonaws.com"
                    }
                }
            )
        ]

        backup_policy = iam.Policy(self, "BackupPolicy", policy_name=f"{Stack.of(self).stack_name}BackupPolicy", statements=backup_statements)
        restore_policy = iam.Policy(self, "RestorePolicy", policy_name=f"{Stack.of(self).stack_name}RestorePolicy", statements=restore_statements)

        role.attach_inline_policy(backup_policy)
        role.attach_inline_policy(restore_policy)
