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

# package.json
sed \
    --in-place \
    --regexp-extended \
    --expression 's~"version": "[0-9\.]+"~"version": "'"${VERSION}"'"~' \
        ./package.json

# Helm Chart
sed \
    --in-place \
    --regexp-extended \
    --expression 's~appVersion: "v[0-9\.]+"~appVersion: "v'"${VERSION}"'"~' \
        ./charts/fivenet/Chart.yaml

HELM_CHART_VERSION=$(grep \
    --only-matching \
    --perl-regexp \
    '^version: ([0-9\.]+)' \
        ./charts/fivenet/Chart.yaml)

version_rest=$(echo "${HELM_CHART_VERSION}" | cut -d':' -f 2 | cut -d'.' -f 1-2)
version_rest="${version_rest// /}"
version_patch=$(echo "${HELM_CHART_VERSION}" | rev | cut -d'.' -f 1 | rev)
version_patch_new=$(( version_patch + 1 ))

sed \
    --in-place \
    --regexp-extended \
    --expression 's~^version: [0-9\.]+~version: '"${version_rest}.${version_patch_new}"'~' \
        ./charts/fivenet/Chart.yaml

git add --all

git commit \
    --signoff \
    --gpg-sign \
    --message "version: bump to v${VERSION}"

git push
echo "Pushing the version bump commit and sleeping 60 seconds before tagging"
sleep 60

git tag "v${VERSION}"
git push --tags
