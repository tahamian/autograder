import pytest
import scripts


class TestHelloWorld:

    def test_program(self):
        expected_output = "Hello World"
        output = scripts.get_result(filename="scripts/tests/hello_world.py", stdout=True, functions=None)

        assert expected_output == output['stdout'].decode('utf-8').rstrip()


# class TestFunction:
#     pass

