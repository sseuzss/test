#!/bin/bash
set -e
set -x

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
echo "${DIR}/build-$(go env GOOS)-$(go env GOARCH).sh"
"${DIR}"/build-$(go env GOOS)-$(go env GOARCH).sh
