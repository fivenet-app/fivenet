#!/usr/bin/env bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
cd "$SCRIPT_DIR" || { echo "Failed to enter script directory."; exit 1; }

set -xe

if [ -z ${1+x} ]; then
    echo "Requires new version as first argument!"
    exit 2
fi

export VERSION="${1//v}"

# Go to root of repo
cd ../../ || { echo "Failed to cd to root of repo."; exit 1; }

echo "v${VERSION}" > VERSION

# package.json
sed \
    --in-place \
    --regexp-extended \
    --expression 's~"version": "[0-9\.]+"~"version": "'"${VERSION}"'"~' \
        ./package.json

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
