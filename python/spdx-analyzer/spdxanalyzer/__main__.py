#!/usr/bin/env python2
import argparse
from qmstr.service.datamodel_pb2 import FileNode, InfoNode
from qmstr.service.controlservice_pb2 import PackageRequest
from qmstr.service.analyzerservice_pb2 import InfoNodeMessage
from qmstr.module.module import QMSTR_Analyzer
from spdx.document import License
import logging
import sys

filename_key = "spdxfile"
fileformat_key = "fileformat"


class SpdxAnalyzer(QMSTR_Analyzer):

    @staticmethod
    def is_primitive(attr):
        primitive = (int, str, bool, float)
        return isinstance(attr, primitive)

    @staticmethod
    def __stringify(t):
        if SpdxAnalyzer.is_primitive(t):
            return t
        if isinstance(t, list):
            return ",".join(map(lambda x: SpdxAnalyzer.__stringify(x), t))
        if isinstance(t, tuple):
            return ":".join(str(i) for i in t)
        if isinstance(t, dict):
            return ",".join(map(lambda x: SpdxAnalyzer.__stringify(x), t.iteritems()))
        members = SpdxAnalyzer.__membersof(t)
        return ",".join(map(lambda y: SpdxAnalyzer.__stringify(y), zip(members, map(lambda x: SpdxAnalyzer.__stringify(t.__getattribute__(x)), members))))

    @staticmethod
    def __membersof(t):
        return [attr for attr in dir(t) if not callable(getattr(t, attr)) and not attr.startswith("_")]

    def __init__(self, address, aid):
        super(SpdxAnalyzer, self).__init__(address, aid)
        self.parse_func_map = {
            'rdf': self.__parse_rdf,
            'tag': self.__parse_tagvalue
        }

    def configure(self, config_map):
        logging.info("Configuring spdx analyzer module")
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

    def analyze(self):
        self.doc = self._parse_spdx()
        self._process_filenodes()
        self._process_packagenode()

    def _process_filenodes(self):

        query_node = FileNode(
            fileType=FileNode.SOURCE
        )
        stream_resp = self.cserv.GetFileNode(query_node)

        for node in stream_resp:
            logging.info("Analyze node {}".format(node.path))
            filtered_files = list(filter(
                lambda f: node.path.endswith(f.name), self.doc.files))
            if not filtered_files:
                logging.warn(
                    "File {} not found in SPDX document".format(node.path))
                continue
            spdx_doc_file_info = filtered_files[0]
            if not isinstance(spdx_doc_file_info.conc_lics, License):
                continue

            data_nodes = []
            logging.info("Concluded license {}".format(
                spdx_doc_file_info.conc_lics))
            data_nodes.append(InfoNode.DataNode(
                type="spdxIdentifier",
                data=spdx_doc_file_info.conc_lics.identifier
            ))
            data_nodes.append(InfoNode.DataNode(
                type="name",
                data=spdx_doc_file_info.conc_lics.full_name
            ))

            info_node = InfoNode(
                type="license",
                dataNodes=data_nodes
            )

            info_nodes = []
            info_nodes.append(InfoNodeMessage(
                uid=node.uid,
                token=self.token,
                infonode=info_node))

            info_iterator = _generate_iterator(info_nodes)

            self.aserv.SendInfoNodes(info_iterator)

    def post_analyze(self):
        pass

    def _process_packagenode(self):
        logging.info("Processing package node")

        data_nodes = []
        for member in SpdxAnalyzer.__membersof(self.doc.package):
            value = self.doc.package.__getattribute__(member)
            if member == "license_declared":
                data_nodes.append(InfoNode.DataNode(
                    type=member,
                    data=value.full_name
                ))
                continue
            data_nodes.append(InfoNode.DataNode(
                type=member,
                data=SpdxAnalyzer.__stringify(value)
            ))

        info_node = InfoNode(
            type="metadata",
            dataNodes=data_nodes
        )

        package_request = PackageRequest(
        )

        package_node = self.cserv.GetPackageNode(package_request)

        info_nodes = []
        info_nodes.append(InfoNodeMessage(
            uid=package_node.uid,
            token=self.token,
            infonode=info_node))

        info_iterator = _generate_iterator(info_nodes)

        self.aserv.SendInfoNodes(info_iterator)

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


def _generate_iterator(collection):
    for i in collection:
        yield i


def main():
    logging.basicConfig(level=logging.INFO)
    logging.info("This is the qmstr spdx analyzer")
    parser = argparse.ArgumentParser()
    parser.add_argument("--aserv", help="qmstr-master address")
    parser.add_argument("--aid", help="analyzer id", type=int)
    args = parser.parse_args()
    spdx_analyzer = SpdxAnalyzer(args.aserv, args.aid)
    spdx_analyzer.run_analyzer()


if __name__ == "__main__":
    main()
