#!/bin/bash

set -e

SCRIPT_DIR="$(dirname "$0")"

docker run --rm -it \
        -v $(realpath "$SCRIPT_DIR/../../../../"):/go/src/ \
        -v $(go env GOMODCACHE):/go/pkg/mod \
        --workdir=/go/src/soldr \
        vxcontrol/vxbuild-cross /bin/bash -c "build/package/agent/build-$(go env GOOS)-$(go env GOARCH).sh"
