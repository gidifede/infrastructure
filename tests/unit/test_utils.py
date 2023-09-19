from infrastructure.constructs.utils import get_version 


def test_get_version():

    version = get_version("tests/unit/assets/test_version.txt")
    no_version = get_version("tests/unit/assets/test_with_no_version.txt")
    version2 = get_version("tests/unit/assets/test_version2.txt")
    no_path = get_version("tests/unit/assets/no_path.txt")
    assert version == "0.0.1-SNAPSHOT"
    assert no_version == None
    assert version2 == "1.1.1-SNAPSHOT"
    assert no_path == None
