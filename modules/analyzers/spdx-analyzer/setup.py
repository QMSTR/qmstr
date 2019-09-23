from setuptools import setup
import os

setup(
    name='pyqmstr-spdx-analyzer',
    version=os.environ["QMSTR_VERSION"],
    description='QMSTR SPDX-Analyzer',
    url='http://qmstr.org',
    license='GPLv3',

    packages=['spdxanalyzer'],
    install_requires=["pyqmstr=={}".format(
        os.environ["QMSTR_VERSION"]), 'spdx-tools==0.5.4'],
    entry_points={
        'console_scripts': [
            'pyqmstr-spdx-analyzer = spdxanalyzer.__main__:main',
        ],
    },

)
