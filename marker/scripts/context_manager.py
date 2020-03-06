from dis import get_instructions
import logging

from contextlib import contextmanager, redirect_stdout
from importlib.util import module_from_spec, spec_from_file_location
from io import StringIO
from unittest.mock import patch

blacklist = ['os', 'sys', 'subprocess', 'threading', 'cmd', 'shlex', 'trace', 'tracemalloc', 'gc', 'pipes', 'syslog',
             'termios', 'posix', 'multiprocessing', 'concurrent', '_dummy_thread', 'dummy_threading']

types = {
    'str': lambda x: str(x),
    'int': lambda x: int(x),
    'float': lambda x: float(x),
    'complex': lambda x: complex(x),
    'bool': lambda x: bool(x),
    'list': lambda x: list(x),
    'dict': lambda x: dict(x)
}

logger = logging.getLogger(__name__)

@contextmanager
def load_module(module_path, redirect=True):
    with open(module_path) as f:
        data = ''.join(f.readlines())
        program = get_instructions(data)
        imports = [__ for __ in program if 'IMPORT' in __.opname]

    imported_libs = list(map(lambda x: x.argval, imports))

    if any([imported_lib in blacklist for imported_lib in imported_libs]):
        print('using black listed package')

    spec = spec_from_file_location("module.sub_mod", module_path)
    sub_mod = module_from_spec(spec)

    with StringIO() as buf, redirect_stdout(buf):
        with patch('builtins.input', side_effect='0'):
            spec.loader.exec_module(sub_mod)

        yield (sub_mod, buf.getvalue())
        del sub_mod


def exec_function(fn, *args):
    status = 1
    with StringIO() as buf, redirect_stdout(buf):
        try:
            retval = fn(*args)
            status = 0
        except Exception as e:
            retval = e
        return retval, buf.getvalue(), status


class Assignment:

    def __init__(self, filename, stdout, functions):
        self.filename = filename
        self.stdout = stdout
        self.functions = list(map(lambda x: Function(x), functions))

    def get_stdout(self):
        if self.stdout:
            return "Hello"

        return None

    def get_functions(self):
        return list(map(lambda x: x.evaluate_function(self.filename), self.functions))


class Function:

    def __init__(self, function):
        self.function_name = function['function_name']
        self.function_args = list(map(lambda x: types[x['type']](x['value']), function['function_args']))

    def evaluate_function(self, module):
        with load_module(module) as (sub_mod, buf):
            f = getattr(sub_mod, self.function_name)
            result, buffer, status = exec_function(f, self.function_args)
        return dict(result=result, buffer=buffer, status=status, function_name=self.function_name)


