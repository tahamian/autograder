import json
import logging

import click

from .grader import get_result


@click.command()
@click.option('--config-file', help='Config file for running', required=True, type=click.Path())
@click.option('--output-file', help='Output file path', required=True, type=click.Path())
@click.option('--log-file', help='Log file path', type=click.Path(), default='marker.log')
def run_script(config_file, output_file, log_file):
    logging.basicConfig(
        filemode='a',
        format='%(asctime)s %(levelname)s-%(message)s',
        datefmt='%Y-%m-%d %H:%M:%S',
        level=logging.DEBUG,
        filename=log_file,
    )

    output = get_result(config_file)
    with open(output_file, 'w', encoding='utf-8') as f:
        json.dump(output, f, ensure_ascii=False, indent=4)
