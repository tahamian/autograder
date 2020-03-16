""" cli tool that runs python scripts can gets output """

__version__ = "0.1"

__license__ = "MIT"
__author__ = "Mostafa Ayesh"
__maintainer__ = "Mostafa Ayesh"

import click
import json
import logging

from .context_manager import Assignment

logger = logging.getLogger(__name__)


@click.command()
@click.option('--config-file', help='get a stdout of python file', required=True, type=click.Path())
@click.option('--log-file', help='set the path of the log file', type=click.Path(), default='marker.log')
def run_script(config_file, log_file, *args, **kwargs):
    logging.basicConfig(filemode='a', format='%(asctime)s %(levelname)s-%(message)s',
                        datefmt='%Y-%m-%d %H:%M:%S', level=logging.DEBUG, filename=log_file)

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

    logger.info(output)

    return output
