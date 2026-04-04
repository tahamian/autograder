import json
from pathlib import Path

from marker import get_result

TESTS_DIR = Path(__file__).parent


class TestHelloWorld:

    def test_stdout(self):
        with open(TESTS_DIR / "test_hello_world" / "output.json") as f:
            expected_output = json.load(f)

        output = get_result(str(TESTS_DIR / "test_hello_world" / "input.json"))
        assert expected_output["stdout"] == output["stdout"]

    def test_function_output(self):
        with open(TESTS_DIR / "test_hello_world" / "output.json") as f:
            expected_output = json.load(f)

        output = get_result(str(TESTS_DIR / "test_hello_world" / "input.json"))
        assert expected_output["functions"] == output["functions"]


class TestPythagoreanTheorem:

    def test_stdout(self):
        with open(TESTS_DIR / "test_pythagorean_theorem" / "output.json") as f:
            expected_output = json.load(f)

        output = get_result(str(TESTS_DIR / "test_pythagorean_theorem" / "input.json"))
        assert expected_output["functions"] == output["functions"]

    def test_function_output(self):
        with open(TESTS_DIR / "test_pythagorean_theorem" / "output.json") as f:
            expected_output = json.load(f)

        output = get_result(str(TESTS_DIR / "test_pythagorean_theorem" / "input.json"))
        assert expected_output["functions"] == output["functions"]
