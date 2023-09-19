"""
This module implements a base Construct
"""
import base64
import hashlib
from aws_cdk import Stack, Tags
from constructs import Construct
from infrastructure.constructs.contextconfig import ContextConfig


class BaseConstruct(Construct):
    '''
    Base construct
    '''
    def __init__(self, scope: "Construct", id: str, config: ContextConfig) -> None:
        super().__init__(scope, id)
        self.stack_name = Stack.of(self).stack_name
        self.region = Stack.of(self).region
        self.config = config

        key_list, value_list = config.get_resources_tags()
        for i in range(len(key_list)):
            self.tag = Tags.of(self).add(key=key_list[i], value=value_list[i])

    def get_name(self, contruct_name: str):
        """ Class method to return the full contruct name"""

        name = f"{self.stack_name.lower()}-{contruct_name.lower()}"
        if len(name) >= 63:
            m_hash = hashlib.sha256()
            m_hash.update(name.encode())
            suffix = (base64.b64encode(m_hash.digest())).decode()[0:12]
            suffix = suffix.replace("+", "-")
            suffix = suffix.replace("/", "-")
            suffix = suffix.lower()
            name = name[0:50]+suffix

        return name
