#!/bin/bash
set -e

source /common.inc

start_dgraph

create_qmstr_user

start_qmstr
