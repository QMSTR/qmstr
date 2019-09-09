from setuptools import Command

import pkg_resources

class QMSTRCommand(Command):

    """QMSTR setuptools Command"""

    description = "create build graph for the python module"

    user_options = []

    def initialize_options(self):
        """init options"""
        pass

    def finalize_options(self):
        """finalize options"""
        pass

    def run(self):
        """runner"""

        self.reinitialize_command('bdist_dumb', inplace=0, format="gztar", keep_temp=True, bdist_dir="/tmp/qmstrsomething")
        self.run_command('bdist_dumb')
