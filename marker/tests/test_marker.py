import pytest
import scripts


class TestHelloWorld:

    def test_program(self):
        expected_output = "Hello World"
        # output = scripts.get_result(filename="marker/scripts/tests/hello_world.py", stdout=True, functions=None)
        # assert 'success' in output['stdout']
        # assert expected_output == output['stdout']['success'].decode('utf-8').rstrip()

        pass
class TestFunction:

    def test_pythagorean(self):
        a = '3'
        b = '4'
        c = '5.0'

        # output = scripts.get_result(filename="marker/scripts/tests/pythagorean_theorem.py", stdout=False,
        #                    functions=[{'name': 'pythagorean', 'args': [a, b]}])
        #
        # assert 'pythagorean' in output['functions']
        # assert 'success' in output['functions']['pythagorean']
        # assert output['functions']['pythagorean']['success'].rstrip() == c


