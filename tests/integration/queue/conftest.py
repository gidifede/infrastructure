import pytest

# region Fixtures
@pytest.fixture(scope="package")
def app_name():
    yield "tests/integration/queue/app.py"

# endregion
