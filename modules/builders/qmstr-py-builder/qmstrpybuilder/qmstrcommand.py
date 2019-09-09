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

        # Ensure metadata is up-to-date
        self.reinitialize_command('build_py', inplace=0)
        self.run_command('build_py')
        bpy_cmd = self.get_finalized_command("build_py")
        build_path = pkg_resources.normalize_path(bpy_cmd.build_lib)

        # Build extensions
        self.reinitialize_command('egg_info', egg_base=build_path)
        self.run_command('egg_info')

        self.reinitialize_command('build_ext', inplace=0)
        self.run_command('build_ext')

        self.reinitialize_command('bdist_dumb', inplace=0, format="gztar", keep_temp=True)
        self.run_command('bdist_dumb')
