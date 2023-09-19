import pytest
import boto3
from tests.integration import util

# region Fixtures
@pytest.fixture(scope="package")
def app_name():
    yield "tests/integration/access/app.py"

# endregion
