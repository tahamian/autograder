""" cli tool that runs python scripts can gets output """

import os
import sys
import subprocess
import click

__version__ = "0.1"

__license__ = "MIT"
__author__ = "Mostafa Ayesh"
__maintainer__ = "Mostafa Ayesh"


class Assignment:

    def __init__(self, script_name, stdout=False, functions=None):
        self.scripts_name = script_name
        self.stdout = stdout
        self.functions = functions


@click.command()
@click.option('--filename', help='name of the file that being run')
@click.option('--stdout', help='get a stdout of python file', default=False)
@click.option('--functions', help='get the function output')
# @click.pass_context
def run_script(filename, stdout=False, functions=None, *args, **kwargs):
    return get_result(filename, stdout, functions)


def get_result(filename, stdout, functions):
    output = {
        'stdout': None,
        'functions': None
    }

    if stdout:
        try:
            output['stdout'] = subprocess.check_output(['python3', filename])

        except Exception as e:
            raise Exception(str(os.getcwd()))

    return output


if __name__ == "__main__":
    run_script()
