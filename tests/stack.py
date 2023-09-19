""" Stack Class for testing """
import os

class Stack():
    """ Class Stack. Implements the CDK deploy/destroy of the test stack """
    def __init__(self, app:str, tid:str):
        self.app=app
        self.tid=tid

    def create(self):
        """ Creates the stack """
        print(f"cdk deploy -c test_id={self.tid} -c config=test --require-approval never -a 'python3 {self.app}'")
        os.system(f"cdk deploy -c test_id={self.tid} -c config=test --require-approval never -a 'python3 {self.app}'")

    def destroy(self):
        """ Destroys the stack"""
        os.system(f"cdk destroy -c test_id={self.tid} -c config=test --force -a 'python3 {self.app}'")
