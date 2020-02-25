""" cli tool that runs python scripts can gets output """

import os
import sys
import subprocess
import click
import ast

__version__ = "0.1"

__license__ = "MIT"
__author__ = "Mostafa Ayesh"
__maintainer__ = "Mostafa Ayesh"


class Assignment:

    def __init__(self, script_name, stdout=False, functions=None):
        self.scripts_name = script_name
        self.stdout = stdout
        self.functions = functions


class PythonLiteralOption(click.Option):

    def type_cast_value(self, ctx, value):
        try:
            return ast.literal_eval(value)
        except:
            raise click.BadParameter(value)


@click.command()
@click.option('--filename', help='name of the file that being run', type=click.Path())
@click.option('--stdout', help='get a stdout of python file', default=False)
@click.option('--functions', help='get the function output', cls=PythonLiteralOption, default=[])
def run_script(filename, stdout, functions, *args, **kwargs):
    return get_result(filename, stdout, functions)


def get_result(filename, stdout, functions):
    output = {
        'stdout': {},
        'functions': {}
    }

    if stdout:
        try:
            output['stdout'] = {'success': subprocess.check_output(['python3', filename])}

        except Exception as e:
            output['stdout'] = {'error': str(e)}

    if bool(functions):

        for func in functions:
            name = func['name']
            args = ''
            if func['args']:
                args = ','.join(func['args'])
            try:
                import_path = (os.path.splitext(filename)[0]).replace('/', '.')

                command = 'from {} import *; print({}({}))'.format(import_path, name, args)
                func_output = subprocess.check_output(['python3', '-c', command])

                output['functions'][func['name']] = dict(success=func_output.decode('utf-8'))

            except Exception as e:
                output['functions'][func['name']] = {'error': str(e)}

    return output


if __name__ == "__main__":
    run_script()
