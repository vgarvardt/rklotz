#!/bin/bash
set -e

MY_PATH="$(dirname "$0")"        # relative
DIR="$( cd "$MY_PATH" && pwd )"  # absolutized and normalized

fly -t rklotz set-pipeline -p rklotz -c ${DIR}/pipeline.yml -l ${DIR}/vars.yml
