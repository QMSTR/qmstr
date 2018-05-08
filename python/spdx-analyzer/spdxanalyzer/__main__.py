#!/usr/bin/env python2
import argparse
from pyqmstr.module.module import Analyzer
import logging


class SpdxAnalyzer(object):
    def __init__(self):
        pass

    def configure(self, config_map):
        print("Configure module")

    def analyze(self, node):
        print("Analyze node")


def main():
    logging.basicConfig(level=logging.INFO)
    logging.info("This is the qmstr spdx analyzer")
    parser = argparse.ArgumentParser()
    parser.add_argument("--aserv", help="qmstr-master address")
    parser.add_argument("--aid", help="analyzer id", type=int)
    args = parser.parse_args()
    spdx_analyzer = Analyzer(SpdxAnalyzer(), args.aserv, args.aid)
    spdx_analyzer.run_analyzer()


if __name__ == "__main__":
    main()
