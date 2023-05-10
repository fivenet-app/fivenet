#!/bin/bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
cd "$SCRIPT_DIR" || { echo "Failed to enter script directory."; exit 1; }

set -xe

if [ -z ${1+x} ]; then
    echo "Requires new version as first argument!"
    exit 2
fi

export VERSION="$1"

# Go to root of repo
cd ../../ || { echo "Failed to cd to root of repo."; exit 1; }

echo "v${VERSION}" > VERSION

sed \
    --in-place \
    --regexp-extended \
    --expression 's~"version": "[\d.]+",~"version": "'"${VERSION}"'",~' \
        ./package.json \
        ./gen/js/package.json

sed \
    --in-place \
    --regexp-extended \
    --expression 's~appVersion: "v[\d.]+"~appVersion: "v'"${VERSION}"'"~' \
        ./charts/fivenet/Chart.yaml

make yarn-upgrade-gen-js
