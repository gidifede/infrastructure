from aws_cdk import aws_sns as sns
import aws_cdk as cdk
import json
from aws_cdk.assertions import Template
from infrastructure.constructs import constants
from infrastructure.constructs.base import bucket
from aws_cdk import Names


def test_bucket(config):
    stack = cdk.Stack()
    
    s3 = bucket.Bucket(stack, "TestBucket", config=config, bucket_name="TestBucket")
    s3.add_asset(path="tests/unit/constructs/base", prefix="test/config")
    
    template = Template.from_stack(stack)
    
    # Assert that the correct number of resources are being created
    template.resource_count_is("AWS::S3::Bucket", 1)
    template.resource_count_is("AWS::IAM::Policy", 1)
    
    # Assert that policy is being created with the right parameters 
    template.has_resource_properties("AWS::IAM::Policy", props= {
        'PolicyDocument':{
               'Statement':[
                  {
                     'Action':[
                        's3:GetObject*',
                        's3:GetBucket*',
                        's3:List*'
                     ],
                     'Effect':'Allow',
                     'Resource':[
                        {
                           'Fn::Join':[
                              '',
                              [
                                 'arn:',
                                 {
                                    'Ref':'AWS::Partition'
                                 },
                                 ':s3:::',
                                 {
                                    'Fn::Sub':'cdk-hnb659fds-assets-${AWS::AccountId}-${AWS::Region}'
                                 }
                              ]
                           ]
                        },
                        {
                           'Fn::Join':[
                              '',
                              [
                                 'arn:',
                                 {
                                    'Ref':'AWS::Partition'
                                 },
                                 ':s3:::',
                                 {
                                    'Fn::Sub':'cdk-hnb659fds-assets-${AWS::AccountId}-${AWS::Region}'
                                 },
                                 '/*'
                              ]
                           ]
                        }
                     ]
                  },
                  {
                     'Action':[
                        's3:GetObject*',
                        's3:GetBucket*',
                        's3:List*',
                        's3:DeleteObject*',
                        's3:PutObject',
                        's3:PutObjectLegalHold',
                        's3:PutObjectRetention',
                        's3:PutObjectTagging',
                        's3:PutObjectVersionTagging',
                        's3:Abort*'
                     ],
                     'Effect':'Allow',
                     'Resource':[
                        {
                           'Fn::GetAtt':[
                              'TestBucketCfnTestBucketD97F5813',
                              'Arn'
                           ]
                        },
                        {
                           'Fn::Join':[
                              '',
                              [
                                 {
                                    'Fn::GetAtt':[
                                       'TestBucketCfnTestBucketD97F5813',
                                       'Arn'
                                    ]
                                 },
                                 '/*'
                              ]
                           ]
                        }
                     ]
                  }
               ],
               'Version':'2012-10-17'
            },
            'PolicyName':'CustomCDKBucketDeployment8693BB64968944B69AAFB0CC9EB8756CServiceRoleDefaultPolicy88902FDF',
            'Roles':[
               {
                  'Ref':'CustomCDKBucketDeployment8693BB64968944B69AAFB0CC9EB8756CServiceRole89A01265'
               }
            ]
    })