import pytest

# region Fixtures
@pytest.fixture(scope="package")
def app_name():
    yield "tests/integration/eventstore_cdc/app.py"

# endregion
