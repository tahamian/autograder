import os
import subprocess
import click
# from importlib import __import__

class Assignment:

    def __init__(self, script_name, stdout=False, functions=None):
        self.scripts_name = script_name
        self.stdout = stdout
        self.functions = functions




@click.command()
@click.option('--filename', help='name of the file that being run')
@click.option('--stdout', help='get a stdout of python file')
@click.command('--functions', help='get the function output')
def run_script(filename, stdout, function_names):

    # file = __import__(filename)
    #
    # if stdout:
    #     file()
    pass

if __name__ == "__main__":
    run_script()