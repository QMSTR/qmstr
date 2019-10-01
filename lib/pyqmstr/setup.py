from setuptools import setup, find_packages
import os


setup(
    name='pyqmstr',
    version=os.environ["QMSTR_VERSION"],
    description='Interface with qmstr-master from python',
    url='http://qmstr.org',
    license='GPLv3',
    packages=find_packages(exclude=['tests']),
    install_requires=["grpcio=={}".format(
        os.environ["GRPCIO_VERSION"]), 'protobuf'],
)
