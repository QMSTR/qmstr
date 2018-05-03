#!/usr/bin/env python2
import argparse
from pyqmstr.module.module import Analyzer


class SpdxAnalyzer(object):
    def __init__(self):
        pass

    def configure(self, config_map):
        print("Configure module")

    def analyze(self, node):
        print("Analyze node")


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--aserv", help="increase output verbosity")
    parser.add_argument("--aid", help="increase output verbosity", type=int)
    args = parser.parse_args()
    spdx_analyzer = Analyzer(SpdxAnalyzer(), args.aserv, args.aid)
    spdx_analyzer.run_analyzer()


if __name__ == "__main__":
    main()
