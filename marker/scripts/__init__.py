""" cli tool that runs python scripts can gets output """

import click
import json

from .context_manager import Assignment

__version__ = "0.1"

__license__ = "MIT"
__author__ = "Mostafa Ayesh"
__maintainer__ = "Mostafa Ayesh"


@click.command()
@click.option('--config-file', help='get a stdout of python file', required=True, type=click.Path())
def run_script(config_file, *args, **kwargs):
    return get_result(config_file)


def get_result(config_file):

    with open(config_file) as f:
        input = json.load(f)

    filename = input['filename']
    stdout = input['stdout']
    function = input['functions']

    output = {
        'stdout': {},
        'functions': {}
    }


    assignment = Assignment(filename, stdout, function)

    output['stdout'] = assignment.get_stdout()
    output['functions'] = assignment.get_functions()

    print(output)

    return output

