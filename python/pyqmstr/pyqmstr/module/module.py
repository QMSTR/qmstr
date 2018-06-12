import grpc
from pyqmstr.service.analyzerservice_pb2 import AnalyzerConfigRequest
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
        self.cserv = ControlServiceStub(
            channel
        )

    def get_name(self):
        return self.name

    def set_package_node(self, pkg):
        self.pkg = pkg

    def get_package_node(self):
        return self.pkg

    def configure(self, config):
        raise NotImplementedError()


class QMSTR_Analyzer(QMSTR_Module):
    def __init__(self, address, aid):
        super(QMSTR_Analyzer, self).__init__(address, aid)

    def run_analyzer(self):
        conf_request = AnalyzerConfigRequest(
            analyzerID=self.id)
        conf_response = self.aserv.GetAnalyzerConfig(conf_request)
        self.configure(conf_response.configMap)

        package_request = PackageRequest(
            session=conf_response.session
        )

        package_response = self.cserv.GetPackageNode(package_request)
        self.set_package_node(package_response.packageNode)

        self.analyzer_module.analyze(self.cserv)

        self.post_analyze()

    def analyze(self, cserv):
        raise NotImplementedError()

    def post_analyze(self):
        raise NotImplementedError()
