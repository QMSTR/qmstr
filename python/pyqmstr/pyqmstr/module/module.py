import grpc
from pyqmstr.service.analyzerservice_pb2 import AnalyzerConfigRequest, NodeRequest, AnalysisMessage
from pyqmstr.service.analyzerservice_pb2_grpc import AnalysisServiceStub
from pyqmstr.service.controlservice_pb2 import PackageRequest
from pyqmstr.service.controlservice_pb2_grpc import ControlServiceStub
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
        self.cserv = ControlServiceStub(channel)

    def getName(self):
        return self.name

    def setPackageNode(self, pkg):
        self.pkg = pkg

    def getPackageNode(self):
        return self.pkg

    def analyze(self, node):
        raise NotImplementedError()


class Analyzer(QMSTR_Module):
    def __init__(self, module, address, aid):
        super(Analyzer, self).__init__(address, aid)
        self.analyzer_module = module

    def run_analyzer(self):
        conf_request = AnalyzerConfigRequest(
            analyzerID=self.id)
        conf_response = self.aserv.GetAnalyzerConfig(conf_request)
        self.analyzer_module.configure(conf_response.configMap)

        package_request = PackageRequest(
            session=conf_response.session
        )

        package_response = self.cserv.GetPackageNode(package_request)
        self.analyzer_module.setPackageNode(package_response.packageNode)

        node_request = NodeRequest(
            type=conf_response.typeSelector
        )

        node_response = self.aserv.GetNodes(node_request)

        for node in node_response.fileNodes:
            self.analyzer_module.analyze(node)

        pkg_node = self.getPackageNode()

        ana_msg = AnalysisMessage(
            token=conf_response.token,
            packageNode=pkg_node,
            resultMap=None
        )
        anaresp = self.aserv.SendNodes(ana_msg)
