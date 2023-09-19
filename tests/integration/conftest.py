""" Pytest Conftest """
import os
from datetime import datetime
import pytest
from tests.integration import util
from tests.stack import Stack
import boto3
import itertools


def pytest_addoption(parser):
    """ Add option for pytest """
    parser.addoption("--create-stack", action="store", default="true")
    parser.addoption("--tid", action="store", default="")
    parser.addoption("--profile", action="store", default="default")


@pytest.fixture(scope="session", autouse=True)
def tests_setup_and_teardown(request):
    # Will be executed before the first test
    old_environ = dict(os.environ)
    os.environ.update({"AWS_PROFILE": request.config.getoption("--profile")})

    yield
    # Will be executed after the last test
    os.environ.clear()
    os.environ.update(old_environ)


@pytest.fixture(scope="session")
def must_create_stack(request):
    """ Fixture to return if stack must be created """
    return request.config.getoption("--create-stack")


@pytest.fixture(scope="session")
def tid_option(pytestconfig):
    """ Fixture to add tid option """
    return pytestconfig.getoption("tid")


@pytest.fixture(scope="package")
def tid(request, must_create_stack):
    """ Fixture to get tid option or to create tid"""
    if must_create_stack != "true":
        yield request.config.getoption("--tid")
    else:
        yield datetime.now().strftime("%Y%m%dT%H%M%S")


@pytest.fixture(scope="package")
def stack(tid, app_name):
    """ Stack fixture """
    yield Stack(app=app_name, tid=tid)


@pytest.fixture(scope="package")
def create_stack(must_create_stack, stack):
    """ Fixture to return if stack must be created """
    if must_create_stack == "true":
        stack.create()

    yield stack

    if must_create_stack == "true":
        clean_resources(stack.tid)
        stack.destroy()


@pytest.fixture(scope="module")
def event_topic(tid):
    """ Fixture to return a topic """
    topics = util.get_resources_by_tag(tid, ["sns:topic"])
    yield list(filter(lambda x: "EventStoreCDC" in x["ResourceARN"], topics))[0]


@pytest.fixture(scope="module")
def sqs_queue(tid):
    """ Fixture to return a queue"""
    queues = util.get_resources_by_tag(tid, ["sqs:queue"])
    # remove DLQ Queue from list
    queues = list(itertools.filterfalse(lambda x: "TestQueueDLQ" in x["ResourceARN"], queues))
    yield list(filter(lambda x: "TestQueue" in x["ResourceARN"], queues))[0]


@pytest.fixture(scope="module")
def sqs_dlq_queue(tid):
    """ Fixture to return a queue"""
    queues = util.get_resources_by_tag(tid, ["sqs:queue"])
    yield list(filter(lambda x: "TestQueueDLQ" in x["ResourceARN"], queues))[0]



@pytest.fixture(scope="module")
def sqs_dlq_command_handler(tid):
    """ Fixture to return a queue"""
    queues = util.get_resources_by_tag(tid, ["sqs:queue"])
    yield list(filter(lambda x: "productQueueDLQ" in x["ResourceARN"], queues))[0]

@pytest.fixture(scope="module")
def sqs_command_handler(tid):
    """ Fixture to return a queue"""
    queues = util.get_resources_by_tag(tid, ["sqs:queue"])
    yield list(filter(lambda x: "productQueue" in x["ResourceARN"] and "DLQ" not in x["ResourceARN"], queues))[0]

@pytest.fixture(scope="module")
def sqs_dlq_projection(tid):
    """ Fixture to return a queue"""
    queues = util.get_resources_by_tag(tid, ["sqs:queue"])
    yield list(filter(lambda x: "productProjectionQueueDLQ" in x["ResourceARN"], queues))[0]

@pytest.fixture(scope="module")
def sqs_projection(tid):
    """ Fixture to return a queue"""
    queues = util.get_resources_by_tag(tid, ["sqs:queue"])
    yield list(filter(lambda x: "productProjectionQueue" in x["ResourceARN"] and "DLQ" not in x["ResourceARN"], queues))[0]


@pytest.fixture(scope="module")
def sns_projection(tid):
    """ Fixture to return a queue"""
    queues = util.get_resources_by_tag(tid, ["sns:topic"])
    yield list(filter(lambda x: "productEventTopic" in x["ResourceARN"] and "DLQ" not in x["ResourceARN"], queues))[0]





@pytest.fixture(scope="module")
def command_handler_lambda(tid):
    """ Fixture to return the lambda for command handler"""
    lambdas = util.get_lambdas(tid, {"product_CommandHandler": "product_CommandHandler"})
    yield lambdas['product_CommandHandler']


@pytest.fixture(scope="module")
def eventstore_table(tid):
    """ Fixture to return the observer table """
    tables = util.get_dynamodb_tables(tid)
    yield tables['product_events']


@pytest.fixture(scope="module")
def cdc_lambda(tid):
    """ Fixture to return the observer table """
    lambdas = util.get_lambdas(tid, {"Test_CDC": "Test_CDC"})
    yield lambdas['Test_CDC']

@pytest.fixture(scope="package")
def validator_lambda(tid):
    """ Fixture to return the observer table """
    lambdas = util.get_lambdas(tid, {"product_Validator" : "product_Validator"})
    yield  lambdas['product_Validator']

@pytest.fixture(scope="module")
def test_query_lambda(tid):
    """ Fixture to return the lambda test to query projection database"""
    lambdas = util.get_lambdas(tid, {"lambdafunction": "lambdafunction"})
    yield lambdas['lambdafunction']

@pytest.fixture(scope="module")
def api_endpoint(tid):
    """ Fixture to return API endpoint """
    api_gw_resources = util.get_resources_by_tag(tid, ["apigateway"])
    url=util.get_apigateway_endpoint(api_gw_resources)
    yield url

@pytest.fixture
def bucket(tid):
    the_buckets = util.get_buckets(tid, {"TestBucket": "TestBucket"})
    yield the_buckets


@pytest.fixture
def confBucket(tid):
    the_buckets = util.get_buckets(tid, {"ConfigurationBucket": "ConfigurationBucket"})
    yield the_buckets

@pytest.fixture(scope="package")
def query_handler_secret(tid):
    """ Fixture to return a secret """
    secrets = util.get_secret(tid, {"query_handler" : "query_handler"})
    yield  secrets['query_handler']

@pytest.fixture(scope="package")
def projection_secret(tid):
    """ Fixture to return a secret """
    secrets = util.get_secret(tid, {"projection" : "projection_handler"})
    yield  secrets['projection']


def empty_buckets(*buckets):
    s3 = boto3.resource('s3')
    for bucket in buckets:
        if not bucket:
            continue
        try:
            print(f"deleting {bucket}")
            b = s3.Bucket(bucket)
            b.objects.all().delete()
        except:
            print(f"error deleting {bucket}")


def clean_resources(tid):
    try:
        conf_bucket = util.get_buckets(tid, {"ConfigurationConfigurationBucket": "ConfigurationConfigurationBucket"})
    except AssertionError:
        conf_bucket = {}
    try:
        test_bucket = util.get_buckets(tid, {"TestBucket": "TestBucket"})
    except AssertionError:
        test_bucket = {}

    empty_buckets(
        conf_bucket.get("ConfigurationConfigurationBucket"),
        test_bucket.get("TestBucket")
    )
