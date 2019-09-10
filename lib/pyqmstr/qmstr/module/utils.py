import os
from tempfile import mkdtemp
import hashlib
import json


def _ignore_path(path_list):
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
        dirnames[:] = _ignore_path(dirnames)
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
        print("ERROR: ", e)
        return None


def create_temp_folder(prefix="qmstr-"):
    """
    Generates a tmp folder with a given prefix.
    Default prefix: "qmstr-"
    """
    tmpdir = mkdtemp(prefix=prefix)

    return tmpdir
