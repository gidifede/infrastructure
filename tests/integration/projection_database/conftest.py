import pytest

# region Fixtures
@pytest.fixture(scope="package")
def app_name():
    yield "tests/integration/projection_database/app.py"

# endregion
