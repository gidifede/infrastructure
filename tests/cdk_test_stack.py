import argparse
import os
import sys
import random
from tests.stack import Stack
from tests.integration.conftest import clean_resources

class bcolors:
    HEADER = '\033[95m'
    OKBLUE = '\033[94m'
    OKCYAN = '\033[96m'
    OKGREEN = '\033[92m'
    WARNING = '\033[93m'
    FAIL = '\033[91m'
    ENDC = '\033[0m'
    BOLD = '\033[1m'
    UNDERLINE = '\033[4m'

def print_out(tid, app):
    print("---------------------------------------")
    print(f"RUN")
    print(bcolors.OKGREEN + f"pytest -s tests/integration/{app}/ --create-stack false --tid {tid} --profile logistic-dev" + bcolors.ENDC)
    print(f"to run test on this stack \n\n")
    
    print(f"RUN")
    print(bcolors.OKGREEN + f"python tests/cdk_test_stack.py --app {app} --tid {tid} --mode destroy" + bcolors.ENDC)
    print(f"to destroy this stack")
    print("---------------------------------------")

parser = argparse.ArgumentParser(description="create test stack",
                                 formatter_class=argparse.ArgumentDefaultsHelpFormatter)

parser.add_argument("-t", "--tid", help="test id")
parser.add_argument("-a", "--app", required=True, help="integration app to deploy (the folder name under tests/integration)")
parser.add_argument("-m", "--mode", default="create", help="create or destroy")
parser.add_argument("-p", "--profile", default="logistic-dev", help="aws profile to use")
args = parser.parse_args()
config = vars(args)

tid = config["tid"]
app = config["app"]
mode = config["mode"]
profile = config["profile"]

os.environ.update({"AWS_PROFILE": profile})

if not app:
    print("usage: python cdk_test_stack.py --app <app_name>")

if not tid and mode=="create":
    tid = random.randint(10000, 99999)

s = Stack(app=f"tests/integration/{app}/app.py", tid=tid)

if mode=="create":
    s.create()
    print_out(tid, app)
    
if mode=="destroy":
    clean_resources(tid=tid)
    s.destroy()