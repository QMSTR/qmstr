from qmstr.module.module import QMSTR_Builder
from qmstr.module.utils import get_files_list, generate_iterator, hash_file, new_file_node
import logging
import os


class BdistBuilder(QMSTR_Builder):
    def __init__(self, address, work_dir, temp_dir):
        super(BdistBuilder, self).__init__(address)
        self.work_dir = work_dir
        self.temp_dir = temp_dir

    def configure(self):
        # TODO: do we need it?
        pass

    def index(self):
        logging.warn("indexing the %s", self.work_dir)
        file_list = get_files_list(self.work_dir)
        logging.warn("collected files %s", file_list)
        file_nodes = [new_file_node(f) for f in file_list]
        self.send_files(file_nodes)

    def package(self):
        logging.warn("running the py builder")
        file_list = get_files_list(self.temp_dir)
        logging.warn("collected files %s", file_list)
        file_nodes = [new_file_node(f, hash=True) for f in file_list]
        #create package node and connect file_nodes

