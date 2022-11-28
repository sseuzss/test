#!/bin/bash

set -e

SCRIPT_DIR="$(dirname "$0")"

docker run --rm -it \
        -v $(realpath "$SCRIPT_DIR/../../../../"):/go/src/ \
        -v $(go env GOMODCACHE):/go/pkg/mod \
        --workdir=/go/src/soldr \
        vxcontrol/vxbuild-cross:1.19.1 /bin/bash -c "build/package/server/build-local.sh"
