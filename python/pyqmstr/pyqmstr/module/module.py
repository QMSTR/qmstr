import grpc
from pyqmstr.service.analyzerservice_pb2 import AnalyzerConfigRequest, NodeRequest
from pyqmstr.service.analyzerservice_pb2_grpc import AnalysisServiceStub


class QMSTR_Module(object):
    def __init__(self, address, aid):
        self.id = aid
        self.name = ""
        channel = grpc.insecure_channel(address)
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

