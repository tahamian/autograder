import json
import logging

from .sandbox import Assignment

logger = logging.getLogger(__name__)


def get_result(config_file):
    output = {'stdout': {}, 'functions': {}}

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
