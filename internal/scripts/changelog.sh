#!/usr/bin/env bash

NODE_VERSION=$(node -p -e "require('./package.json').version")

if [ "$1" = "stdout" ]; then
    # Second argument can be used to specify the "new" version
    pnpx git-cliff --output - --unreleased --tag "${2:-$NODE_VERSION}"
else
    pnpx git-cliff --unreleased --tag "${NODE_VERSION}" --prepend ./CHANGELOG.md
fi
