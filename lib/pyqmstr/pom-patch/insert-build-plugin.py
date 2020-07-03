#/usr/bin/python3
"""
This script inserts a build plugin into an existing POM file.
"""

import defusedxml.ElementTree as ET
import click

import logging
from pathlib import Path

@click.command()
@click.option('-t', '--target',
              help='Path of the pom.xml to patch. Relative to root directory.')
@click.argument('plugin')
def insert_build_plugin(target, plugin):
    '''
    Search, parse, patch and return specified POM file.
    
    Parameters:
        target(str): Target POM file path
        backup(bool): Whether to create a backup of the original POM file.
        plugin(str): Plugin XML file path
    
    Returns:
        True: The patched pom.xml has been written
        False: Otherwise
    '''
    
    if target:
        pom_path = target
    else:
        pom_path = find_pom()
    if not pom_path:
        return False
    pom_xml = parse_xml(pom_path)
    plugin_xml = parse_xml(Path(plugin))
    if not pom_xml or not plugin_xml:
        return False

    if not plugin_section:
        logging.info('Plugin section not found in pom.xml.')
        build_section = pom_xml.find('./build')
        if not build_section:
            logging.info('Build section not found in pom.xml.')
            build_section = ET.SubElement(pom_xml.getroot(), 'build')
        plugin_section = ET.SubElement(build_section)
    plugin_section.append(plugin_xml.getroot())
    return write_pom(pom_xml, pom_path)

def find_pom():
    '''
    Find POM file in current working directory (recursive)
    '''

    target_candidates = list(Path.cwd().rglob('pom.xml'))
    if len(target_candidates) == 0:
        logging.error('No file named pom.xml found in {}.'.format(Path.cwd()))
        return None
    if len(target_candidates) > 1:
        logging.warning('Multiple POM files found. Using {}.'.format(
            target_candidates[0]), extra=dict(pom_files=target_candidates))
    return Path(target_candidates[0])
            
def write_pom(pom_data, out_path):
    '''
    Write POM to specified file after optionally creating a backup.
    '''

    try:
        pom_data.write(out_path)
        logging.info('Wrote pom.xml to {}'.format(str(out_path)))
        return True
    except PermissionError:
        logging.error('Failed to write pom.xml: Permission denied.')
        return False
    except Exception as e:
        logging.error('Failed to write pom.xml: {}'.format(e))
        return False
    
def parse_xml(path):
    '''
    Parse XML file or return None if a) the file doesn't exist or b) is
    unreadable or c) the XML parser encounters an error.
    '''
    
    logging.info('Parsing {}'.format(str(path)))
    try:
        it = ET.parse(path)
        for _, el in it:
            prefix, has_namespace, postfix = el.tag.partition('}')
            if has_namespace:
                el.tag = postfix
        return it.getroot()
    except ET.ParseError as pe:
        logging.error('Failed to parse {}: {}'.format(str(path), pe))
        return None
    except PermissionError:
        logging.error('Failed to read {}: Permission denied.'.format(str(path)))
        return None

if __name__ == '__main__':
    exit(int(insert_build_plugin()))



