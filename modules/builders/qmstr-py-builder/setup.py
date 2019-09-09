from setuptools import setup
import os

setup(
    name='qmstr-py-builder',
    version=os.environ["QMSTR_VERSION"],
    description='QMSTR Python Builder',
    url='http://qmstr.org',
    license='GPLv3',

    packages=['qmstrpybuilder'],
    install_requires=["pyqmstr=={}".format(
        os.environ["QMSTR_VERSION"]), 'docopt'],
    entry_points={
        'distutils.commands': [
            "qmstr = qmstrpybuilder.qmstrcommand:QMSTRCommand"]}
)
