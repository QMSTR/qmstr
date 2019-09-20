import grpc
from qmstr.service.analyzerservice_pb2 import AnalyzerConfigRequest
from qmstr.service.analyzerservice_pb2_grpc import AnalysisServiceStub
from qmstr.service.controlservice_pb2_grpc import ControlServiceStub
from qmstr.service.buildservice_pb2_grpc import BuildServiceStub
from qmstr.service.datamodel_pb2 import FileNode, PackageNode
from qmstr.module.utils import generate_iterator
import logging


def sanitize_address(address):
    if address[0] == ":":
        logging.warn("sanitizing localhost address")
        return "localhost{}".format(address)
    return address


class QMSTR_Module(object):
    def __init__(self, address):
        self.name = ""
        aserv_address = sanitize_address(address)
        logging.info("Connecting to qmstr-master at %s", aserv_address)
        self.channel = grpc.insecure_channel(aserv_address)
        self.cserv = ControlServiceStub(
            self.channel
        )

    def get_name(self):
        return self.name

    def configure(self, config):
        raise NotImplementedError()


class QMSTR_Analyzer(QMSTR_Module):
    def __init__(self, address, aid):
        super(QMSTR_Analyzer, self).__init__(address)
        self.id = aid
        self.aserv = AnalysisServiceStub(
            self.channel)

    def run_analyzer(self):
        conf_request = AnalyzerConfigRequest(
            analyzerID=self.id)
        conf_response = self.aserv.GetAnalyzerConfig(conf_request)
        self.configure(conf_response.configMap)
        self.token = conf_response.token
        self.analyze()

        self.post_analyze()

    def analyze(self):
        raise NotImplementedError()

    def post_analyze(self):
        raise NotImplementedError()


class QMSTR_Builder(QMSTR_Module):
    def __init__(self, address):
        super(QMSTR_Builder, self).__init__(address)
        self.buildserv = BuildServiceStub(
            self.channel)

    def send_files(self, file_nodes):
        """
        Get file nodes and send it to the server
        """
        logging.warn("sending file nodes %s", file_nodes)
        fn_iterator = generate_iterator(file_nodes)
        response = self.buildserv.Build(fn_iterator)

    def send_package(self, package_node):
        logging.warn("sending packagenode %s", package_node)
        response = self.buildserv.CreatePackage(package_node)