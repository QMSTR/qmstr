from setuptools import setup, find_packages

setup(
    name='pyqmstr',
    version='0.1',
    description='Interface with qmstr-master from python',
    url='http://qmstr.org',
    license='GPLv3',

    packages=find_packages(exclude=['tests']),
    install_requires=['grpcio'],
)
