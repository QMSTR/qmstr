from setuptools import setup

setup(
    name='pyqmstr-spdx-analyzer',
    version='0.2',
    description='QMSTR SPDX-Analyzer',
    url='http://qmstr.org',
    license='GPLv3',

    packages=['spdxanalyzer'],
    install_requires=['pyqmstr', 'spdx-tools'],
    entry_points={
        'console_scripts': [
            'pyqmstr-spdx-analyzer = spdxanalyzer.__main__:main',
        ],
    },

)
