#!/usr/bin/env python2
import argparse
from pyqmstr.module.module import Analyzer
import logging
import sys

filename_key = "spdxfile"
fileformat_key = "fileformat"


class SpdxAnalyzer(object):

    def __init__(self):
        self.parse_func_map = {
            'rdf': self.__parse_rdf,
            'tag': self.__parse_tagvalue
        }

    def configure(self, config_map):
        print("Configuring spdx analyzer module")
        if not filename_key in config_map:
            logging.error(
                "spdx-analyzer misconfigured. {} missing.".format(filename_key))
            sys.exit(2)
        else:
            self.spdx_file = config_map[filename_key]

        if not fileformat_key in config_map:
            if self.spdx_file.endswith(".rdf"):
                self.format = "rdf"
            elif self.spdx_file.endswith(".tag"):
                self.format = "tag"
            else:
                logging.error("unable to guess file format")
                sys.exit(3)
        else:
            self.format = config_map[fileformat_key]

        self.doc = self._parse_spdx()
        self._processPackageNodeData()

    def analyze(self, node):
        logging.info("Analyze node {}".format(node.path))
        filtered_files = filter(
            lambda f: node.path.endswith(f.name), self.doc.files)
        if not filtered_files:
            logging.warn(
                "File {} not found in SPDX document".format(node.path))
            return
        spdx_doc_file_info = filtered_files[0]
        logging.info("Concluded license {}".format(
            spdx_doc_file_info.conc_lics))

    def post_analyze(self):
        logging.info(self.get_package_node())

    def _processPackageNodeData(self):
        logging.warn("Package node not yet available")
        # self.packageNode.Name = self.doc.package.name

    def _parse_spdx(self):
        if not self.format in self.parse_func_map:
            logging.error("Unsupported format {}".format(self.format))
            sys.exit(4)
        with open(self.spdx_file, "r") as spdxfile:
            spdxdata = spdxfile.read()
            doc, error = self.parse_func_map[self.format](spdxdata)
            if error != None:
                logging.error(error)
                sys.exit(5)
            return doc

    def __parse_tagvalue(self, data):
        from spdx.parsers.tagvalue import Parser
        from spdx.parsers.tagvaluebuilders import Builder
        from spdx.parsers.loggers import StandardLogger
        p = Parser(Builder(), StandardLogger())
        p.build()
        document, error = p.parse(data)
        if error:
            return (None, error)
        else:
            return (document, None)

    def __parse_rdf(self, data):
        from spdx.parsers.rdf import Parser
        from spdx.parsers.rdfbuilders import Builder
        from spdx.parsers.loggers import StandardLogger
        p = Parser(Builder(), StandardLogger())
        document, error = p.parse(data)
        if error:
            return (None, error)
        else:
            return (document, None)


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
