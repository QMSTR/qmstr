"""qmstr-py-builder

Usage:
  qmstr-py-builder build <dir>
  qmstr-py-builder (-h | --help)
  qmstr-py-builder --version

Options:
  -h --help                     Show this screen.
  --version                     Show version.

"""

import docopt

class QMSTRPythonBuilder(object):
    def __init__(self, args):
        if args['build']:
            self.mode = build
        self.pythondir = args['dir']

    def start(self):
        pass

def main():
    arguments = docopt.docopt(__doc__, version='0.1')
    qmstrpybuilder = QMSTRPythonBuilder(arguments)
    qmstrpybuilder.start()

if __name__ == "__main__":
    main()