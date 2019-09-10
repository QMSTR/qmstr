from setuptools import Command
from qmstr.lib.pyqmstr import utils
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

        # FIXME: function to generate tmp folder
        tmp_path = utils.create_temp_folder()

        self.reinitialize_command(
            'bdist_dumb',
            inplace=0,
            format="gztar",
            keep_temp=True,
            bdist_dir=tmp_path
        )
        self.run_command('bdist_dumb')

        # TODO: trigger path walk/hash/etc
