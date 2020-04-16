""" cli tool that runs python scripts can gets output """

__version__ = "0.1"

__license__ = "MIT"
__author__ = "Mostafa Ayesh"
__maintainer__ = "Mostafa Ayesh"

import click
import json
import logging
import os

from .context_manager import Assignment

logger = logging.getLogger(__name__)


@click.command()
@click.option('--config-file', help='config File for running', required=True, type=click.Path())
@click.option('--output-file', help='set output file', required=True, type=click.Path())
@click.option('--log-file', help='set the path of the log file', type=click.Path(), default='marker.log')
def run_script(config_file, output_file, log_file, *args, **kwargs):
    logging.basicConfig(filemode='a', format='%(asctime)s %(levelname)s-%(message)s',
                        datefmt='%Y-%m-%d %H:%M:%S', level=logging.DEBUG, filename=log_file)

    output =  get_result(config_file)
    with open(output_file, 'w', encoding='utf-8') as f:
        json.dump(output, f, ensure_ascii=False, indent=4)


def get_result(config_file):
    output = {
        'stdout': {},
        'functions': {}
    }

    try:
        with open(config_file) as f:
            input = json.load(f)

    except FileNotFoundError:
        output['error'] = 'file not found'
        return output

    filename = input['filename']
    stdout = input['stdout']
    function = input['functions']

    assignment = Assignment(filename, stdout, function)
    try:

        output['stdout'] = assignment.get_stdout()
        output['functions'] = assignment.get_functions()

    except Exception as e:
        output['error'] = 'Failed to evaluate {}'.format(e)

    logger.info(output)

    return output
