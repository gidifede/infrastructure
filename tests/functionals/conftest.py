""" Pytest Conftest """
import pytest
import os
import random
import string
from dotenv import load_dotenv

load_dotenv()  # take environment variables from .env.

# File for define configuration that we need on funzional tests

# Method for retrieve commons configuration for call the api
@pytest.fixture
def config():
    return {
        "BASE_API": os.getenv(key="FUNCTIONAL_TESTS_BASE_API"),
        "API_KEY": os.getenv(key="FUNCTIONAL_TESTS_API_KEY")
    }


# Method for help us to generate random product id
@pytest.fixture
def product_id():
    
    return ''.join(random.choices(string.digits, k=10))


@pytest.fixture
def cluster_id():
    
    return ''.join(random.choices(string.digits, k=10))

# def pytest_addoption(parser):
#     pass