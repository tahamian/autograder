import pytest
import json
from marker.scripts import get_result

class TestHelloWorld:

    def test_stdout(self):

        with open('marker/tests/test_hello_world/output.json') as f:
            expected_output = json.load(f)

        output = get_result('marker/tests/test_hello_world/input.json')

        assert expected_output['output']['stdout'] == output['output']['stdout']

    def test_function_output(self):

        with open('marker/tests/test_hello_world/output.json') as f:
            expected_output = json.load(f)

        output = get_result('marker/tests/test_hello_world/input.json')

        assert expected_output['functions'] == output['functions']


class TestPythagoreanTheorem:

    def test_stdout(self):
        with open('marker/tests/test_pythagorean_therom/output.json') as f:
            expected_output = json.load(f)

        output = get_result('marker/tests/test_pythagorean_therom/input.json')
        assert expected_output['output']['functions'] == output['output']['functions']

    def test_function_output(self):
        with open('marker/tests/test_pythagorean_therom/output.json') as f:
            expected_output = json.load(f)

        output = get_result('marker/tests/test_pythagorean_therom/input.json')
        assert expected_output['functions'] == output['functions']
