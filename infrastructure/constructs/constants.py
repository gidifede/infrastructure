"""
Constants
"""
from aws_cdk import(
    BundlingOptions,
    aws_lambda as lambda_,
)

# global var for backup plan
BACKUP_PLAN = None
BACKUP_ROLE = None

# global var for layer (singleton)
LAMBDA_LAYER = None

# Lambda asset root
# this is the path in the repo and must not be changed
# unless the repo changes
LAMBDA_ASSET_ROOT="services"
MOCK_ASSET_ROOT="services/mock"

# S3 parameters
TAG_LIFECYCLE_RULE='LifeCycleRule'

