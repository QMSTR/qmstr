#!/bin/bash
set -e

source /common.inc

create_qmstr_user

start_dgraph

start_qmstr
