import logging
from setuptools import Command
from qmstr.module import utils
from qmstrpybuilder.bdistbuilder import BdistBuilder
import pkg_resources
from tempfile import TemporaryDirectory
import sys
import os


class QMSTRCommand(Command):

    """QMSTR setuptools Command"""


    description = "create build graph for the python module"

    user_options = []

    def initialize_options(self):
        """init options"""
        qmstr_env = 'QMSTR_MASTER'
        try:
            self.master_address = os.environ[qmstr_env]
        except KeyError:
            logging.error("environment variable %s not set", qmstr_env)
            sys.exit(1)

    def finalize_options(self):
        """finalize options"""
        pass

    def run(self):
        """runner"""

        with TemporaryDirectory(prefix="qmstr") as temp_dir:
            bdist_builder = BdistBuilder(self.master_address, os.curdir, temp_dir)
            bdist_builder.index()

            self.reinitialize_command(
                'bdist_dumb',
                inplace=0,
                format="gztar",
                keep_temp=True,
                bdist_dir=temp_dir,
            )
            self.run_command('bdist_dumb')

            bdist_builder.package()
