import pytest
from marker.scripts import run_script


class TestHelloWorld:

    def test_stdout(self):
        # print(dir(marker))
        print("Hello world")

        assert 0 == 0

    def test_function_output(self):
        pass


class TestPythagoreanTheorem:

    def test_stdout(self):
        pass

    def test_function_output(self):
        pass
