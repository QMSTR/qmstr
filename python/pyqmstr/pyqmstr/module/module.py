import grpc
from pyqmstr.service.analyzerservice_pb2 import AnalyzerConfigRequest, NodeRequest
from pyqmstr.service.analyzerservice_pb2_grpc import AnalysisServiceStub
import logging


def sanitize_address(address):
    if address[0] == ":":
        logging.warn("sanitizing localhost address")
        return "localhost{}".format(address)
    return address


class QMSTR_Module(object):
    def __init__(self, address, aid):
        self.id = aid
        self.name = ""
        aserv_address = sanitize_address(address)
        logging.info("Connecting to qmstr-master at %s", aserv_address)
        channel = grpc.insecure_channel(aserv_address)
        self.aserv = AnalysisServiceStub(
            channel)

    def getName(self):
        return self.name


class Analyzer(QMSTR_Module):
    def __init__(self, module, address, aid):
        super(Analyzer, self).__init__(address, aid)
        self.analyzer_module = module

    def run_analyzer(self):
        conf_request = AnalyzerConfigRequest(
            analyzerID=self.id)
        conf_response = self.aserv.GetAnalyzerConfig(conf_request)
        self.analyzer_module.configure(conf_response.configMap)

        node_request = NodeRequest(
            type=conf_response.typeSelector
        )

        node_response = self.aserv.GetNodes(node_request)

        for node in node_response.fileNodes:
            self.analyzer_module.analyze(node)
