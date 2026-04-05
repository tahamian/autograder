import os
import tempfile

import pytest

from marker.sandbox import (
    Assignment,
    BlackListedImport,
    Function,
    exec_function,
    load_module,
    types,
)


class TestExecFunction:

    def test_simple_return(self):
        result, buf, status = exec_function(lambda x: x * 2, [5])
        assert result == 10
        assert status == 0

    def test_exception_returns_error(self):

        def bad():
            raise ValueError("boom")

        result, buf, status = exec_function(bad, [])
        assert status == 1
        assert isinstance(result, ValueError)

    def test_captures_stdout(self):

        def noisy(x):
            print(f"debug: {x}")
            return x

        result, buf, status = exec_function(noisy, [42])
        assert result == 42
        assert "debug: 42" in buf
        assert status == 0

    def test_no_args(self):
        result, buf, status = exec_function(lambda: "hello", [])
        assert result == "hello"
        assert status == 0

    def test_multiple_args(self):
        result, buf, status = exec_function(lambda a, b, c: a + b + c, [1, 2, 3])
        assert result == 6
        assert status == 0


class TestTypeCoercion:

    def test_str(self):
        assert types['str'](42) == '42'

    def test_int(self):
        assert types['int']('7') == 7

    def test_float(self):
        assert types['float']('3.14') == 3.14

    def test_complex(self):
        assert types['complex']('1+2j') == (1 + 2j)

    def test_bool(self):
        assert types['bool'](1) is True
        assert types['bool'](0) is False

    def test_invalid_type_raises(self):
        with pytest.raises(KeyError):
            types['invalid']('x')


class TestLoadModule:

    def test_captures_stdout(self):
        with tempfile.NamedTemporaryFile(mode='w', suffix='.py', delete=False) as f:
            f.write('print("hello from module")')
            f.flush()
            output = load_module(f.name)
        os.unlink(f.name)
        assert output == "hello from module"

    def test_file_not_found(self):
        output = load_module("/nonexistent/path.py")
        assert output == "file not found"

    def test_empty_file(self):
        with tempfile.NamedTemporaryFile(mode='w', suffix='.py', delete=False) as f:
            f.write('')
            f.flush()
            output = load_module(f.name)
        os.unlink(f.name)
        assert output == ""


class TestBlacklist:

    def _write_script(self, code):
        f = tempfile.NamedTemporaryFile(mode='w', suffix='.py', delete=False)
        f.write(code)
        f.flush()
        f.close()
        return f.name

    def test_safe_import_allowed(self):
        path = self._write_script('import math\nprint(math.sqrt(4))')
        try:
            a = Assignment(path, True, None)
            output = a.get_stdout()
            assert output == "2.0"
        finally:
            os.unlink(path)

    def test_os_import_blocked(self):
        path = self._write_script('import os')
        try:
            with pytest.raises(BlackListedImport):
                Assignment(path, False, None)
        finally:
            os.unlink(path)

    def test_subprocess_import_blocked(self):
        path = self._write_script('import subprocess')
        try:
            with pytest.raises(BlackListedImport):
                Assignment(path, False, None)
        finally:
            os.unlink(path)

    def test_sys_import_blocked(self):
        path = self._write_script('import sys')
        try:
            with pytest.raises(BlackListedImport):
                Assignment(path, False, None)
        finally:
            os.unlink(path)


class TestFunction:

    def test_basic_function(self):
        fn = Function(
            {
                'function_name': 'add',
                'function_args': [{
                    'type': 'int',
                    'value': '3'
                }, {
                    'type': 'int',
                    'value': '4'
                }],
                'testcase_name': 'test_add',
            }
        )
        assert fn.function_name == 'add'
        assert fn.function_args == [3, 4]
        assert fn.testcase_name == 'test_add'

    def test_missing_testcase_name(self):
        fn = Function({
            'function_name': 'foo',
            'function_args': [],
        })
        assert fn.testcase_name == ''

    def test_float_args(self):
        fn = Function({
            'function_name': 'calc',
            'function_args': [{
                'type': 'float',
                'value': '3.14'
            }],
        })
        assert fn.function_args == [3.14]

    def test_evaluate_function(self):
        with tempfile.NamedTemporaryFile(mode='w', suffix='.py', delete=False) as f:
            f.write('def add(a, b):\n    return a + b\n')
            f.flush()
            fn = Function(
                {
                    'function_name': 'add',
                    'function_args': [{
                        'type': 'int',
                        'value': '2'
                    }, {
                        'type': 'int',
                        'value': '3'
                    }],
                    'testcase_name': 'tc',
                }
            )
            result = fn.evaluate_function(f.name)
        os.unlink(f.name)
        assert result['result'] == '5'
        assert result['status'] == 0
        assert result['testcase_name'] == 'tc'
