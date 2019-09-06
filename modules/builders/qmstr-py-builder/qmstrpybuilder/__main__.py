"""qmstr-py-builder

Usage:
  qmstr-py-builder build <dir>
  qmstr-py-builder (-h | --help)
  qmstr-py-builder --version

Options:
  -h --help                     Show this screen.
  --version                     Show version.
  --qmstraddr

"""

import docopt
import os
from os.path import join
import sys
import logging

from qmstr.module.module import QMSTR_Builder


class QMSTRPythonBuilder(QMSTR_Builder):
    def __init__(self, address, args):
        super(QMSTRPythonBuilder, self).__init__(address)
        if args['build']:
            self.mode = 'build'
        self.pythondir = args['dir']

    def start(self):
        self.__index()
        pass

    def _index(self):
        for root, dirs, files in os.walk(self.pythondir):
            full_path_files = [join(root, f) for f in files]
            self.send_files(full_path_files)


def main():
    arguments = docopt.docopt(__doc__, version='0.1')

    qmstr_addr = os.environ.get('QMSTR_MASTER')
    if not qmstr_addr:
        arguments.get("--qmstraddr", None)

    if not qmstr_addr:
        logging.error(
            'No qmstr address given; please provide address of qmstr-master')
        sys.exit(1)

    qmstrpybuilder = QMSTRPythonBuilder(qmstr_addr, arguments)
    qmstrpybuilder.start()


if __name__ == "__main__":
    main()
