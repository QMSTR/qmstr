from qmstr.module.module import QMSTR_Builder
from qmstr.module.utils import get_files_list, generate_iterator, hash_file, new_file_node, new_package_node
import logging
import os


class BdistBuilder(QMSTR_Builder):
    def __init__(self, address, work_dir, temp_dir):
        super(BdistBuilder, self).__init__(address)
        self.work_dir = work_dir
        self.temp_dir = temp_dir

    def configure(self):
        pass

    def index(self):
        logging.debug("indexing the %s", self.work_dir)
        file_list = get_files_list(self.work_dir)
        logging.debug("collected files %s", file_list)
        file_nodes = [new_file_node(f) for f in file_list]
        self.send_files(file_nodes)

    def package(self, name, version):
        logging.debug("package %s", self.temp_dir)
        file_list = get_files_list(self.temp_dir)
        logging.debug("collected files %s", file_list)
        file_nodes = [new_file_node(f, hash=True) for f in file_list]
        BdistBuilder.connect_bytecode(file_nodes)
        pkg_node = new_package_node(name, version, file_nodes)
        self.send_package(pkg_node)

    @staticmethod
    def connect_bytecode(file_nodes):
        bytecode_nodes = [f for f in file_nodes if f.path.endswith(".pyc")]
        target_source = dict()
        for bcfn in bytecode_nodes:
            directory, filename = os.path.split(bcfn.path)
            if directory.endswith("__pycache__"):
                directory = os.path.dirname(directory)

            filenameparts = filename.split('.')
            if len(filenameparts) > 2:
                boundary = -2
            else:
                boundary = -1
            source_filename = os.path.join(
                directory, *filenameparts[:boundary])
            source_filename = source_filename + ".py"
            target_source[source_filename] = bcfn

        source_nodes = [f for f in file_nodes if f.path.endswith("py")]
        for scfn in source_nodes:
            try:
                target = target_source[scfn.path]
                target.derivedFrom.append(scfn)
            except KeyError:
                logging.warn("no bytecode for %s found", scfn.path)
