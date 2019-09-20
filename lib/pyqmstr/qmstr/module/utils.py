import os
from tempfile import mkdtemp
import hashlib
import json
import logging
from qmstr.service.datamodel_pb2 import FileNode, InfoNode, PackageNode


def _filter_out_hidden_files(path_list):
    """
    Returns the path list ignoring all paths and files starting with a dot.
    """
    # FIXME: is it needed?
    return [path for path in path_list if not path.startswith(".")]


def get_files_list(path):
    """
    Returns a list of all files from a given directory, including its
    subfolders.
    """
    all_files = []

    for dirpath, dirnames, filenames in os.walk(path, topdown=True):
        dirnames[:] = _filter_out_hidden_files(dirnames)
        for name in filenames:
            filename = (os.path.join(dirpath, name))
            all_files.append(filename)

    return all_files


def hash_file(file_path):
    """
    Simple function to generate SHA1 of a given file.
    """
    BLOCKSIZE = 4096
    hasher = hashlib.sha1()

    try:
        with open(file_path, 'rb') as fp:
            chunk = fp.read(BLOCKSIZE)
            while len(chunk) > 0:
                hasher.update(chunk)
                chunk = fp.read(BLOCKSIZE)
            chunk = fp.read()
            hasher.update(chunk)
        file_hash = hasher.hexdigest()
        return file_hash
    except FileNotFoundError as e:
        # FIXME: should we just ignore files not found?
        logging.error("ERROR: ", e)
        return None

def generate_iterator(collection):
    for i in collection:
        yield i

def new_file_node(path, hash=False):
    """
    Returns a filenode with calculated checksum if hash parameter is True
    """

    if hash:
        chksum = hash_file(path)
    else:
        chksum = None

    file_node = FileNode(
        path=path,
        fileType=FileNode.UNDEF,
        hash=chksum,
        name=os.path.basename(path)
    )

    return file_node

def new_package_node(name, version, file_nodes):
    return PackageNode(
        name=name,
        version=version,
        targets=file_nodes
    )